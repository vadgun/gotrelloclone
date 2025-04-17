package services

import (
	"encoding/json"
	"time"

	"github.com/vadgun/gotrelloclone/board-service/kafka"
	"github.com/vadgun/gotrelloclone/board-service/models"
	"github.com/vadgun/gotrelloclone/board-service/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardService struct {
	repo *repositories.BoardRepository
}

func NewBoardService(repo *repositories.BoardRepository) *BoardService {
	return &BoardService{repo}
}

func (s *BoardService) CreateBoard(name, ownerID, ownerName string) (*models.Board, error) {
	board := &models.Board{
		Name:      name,
		OwnerID:   ownerID,
		OwnerName: ownerName,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateBoard(board)
	if err != nil {
		return nil, err
	}

	board.ID = id.(primitive.ObjectID)

	var kafkaBoard struct {
		ID string `json:"id" bson:"_id"`
	}

	kafkaBoard.ID = board.ID.Hex()

	jsonID, _ := json.Marshal(kafkaBoard)
	go kafka.ProduceMessage("", string(jsonID), "board-events", "new-board")

	return board, nil
}

func (s *BoardService) GetBoardsByUser(userID string) ([]models.Board, error) {
	return s.repo.GetBoardsByUser(userID)
}

func (s *BoardService) GetBoardByID(boardID string) (*models.Board, error) {
	return s.repo.GetBoardByID(boardID)
}

func (s *BoardService) DeleteBoardByID(boardID string) error {

	var kafkaBoard struct {
		ID string `json:"id" bson:"_id"`
	}
	kafkaBoard.ID = boardID

	jsonID, _ := json.Marshal(kafkaBoard)
	go kafka.ProduceMessage("", string(jsonID), "board-events", "drop-board")
	return s.repo.DeleteBoardByID(boardID)
}

func (s *BoardService) UpdateBoardByID(boardID, newBoardName string) error {
	return s.repo.UpdateBoardByID(boardID, newBoardName)
}

func (s *BoardService) GetAllBoards() ([]models.Board, error) {
	return s.repo.GetAllBoards()
}
