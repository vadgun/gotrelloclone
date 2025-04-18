package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/vadgun/gotrelloclone/board-service/infra/config"
	"github.com/vadgun/gotrelloclone/board-service/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardRepository struct {
	collection *mongo.Collection
}

func NewBoardRepository() *BoardRepository {
	return &BoardRepository{
		collection: config.DB.Collection("boards"),
	}
}

func (r *BoardRepository) CreateBoard(board *models.Board) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := r.collection.InsertOne(ctx, board)
	return id.InsertedID, err
}

func (r *BoardRepository) GetBoardsByUser(userID string) ([]models.Board, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var boards []models.Board
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": userID})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &boards); err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *BoardRepository) GetBoardByID(boardID string) (*models.Board, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	boardObjectID, errs := primitive.ObjectIDFromHex(boardID)
	if errs != nil {
		return nil, errs
	}

	var board models.Board
	err := r.collection.FindOne(ctx, bson.M{"_id": boardObjectID}).Decode(&board)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("tablero no encontrado")
		}
		return nil, err
	}

	return &board, nil
}

func (r *BoardRepository) DeleteBoardByID(boardID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	boardObjectID, errs := primitive.ObjectIDFromHex(boardID)
	if errs != nil {
		return errs
	}

	mongoResult, err := r.collection.DeleteOne(ctx, bson.M{"_id": boardObjectID})
	if err != nil {
		return err
	}

	if mongoResult.DeletedCount == 0 {
		return errors.New("tablero no encontrado")
	}

	return nil
}

func (r *BoardRepository) UpdateBoardByID(boardID string, newBoardName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	boardObjectID, errs := primitive.ObjectIDFromHex(boardID)
	if errs != nil {
		return errs
	}

	filter := bson.M{"_id": boardObjectID}
	update := bson.M{"$set": bson.M{"name": newBoardName}}
	mongoResult, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if mongoResult.MatchedCount == 0 {
		return errors.New("tablero no encontrado")
	}

	return nil
}

func (r *BoardRepository) GetAllBoards() ([]models.Board, error) {
	var boards []models.Board

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var board models.Board
		if err := cursor.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	return boards, nil
}
