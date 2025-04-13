package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/vadgun/gotrelloclone/task-service/config"
	"github.com/vadgun/gotrelloclone/task-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository struct {
	collection      *mongo.Collection
	userCollection  *mongo.Collection
	boardCollection *mongo.Collection
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		collection:      config.DB.Collection("tasks"),
		userCollection:  config.DB.Collection("users"),
		boardCollection: config.DB.Collection("boards"),
	}
}

// 1️⃣ Crear tarea
func (r *TaskRepository) CreateTask(ctx context.Context, task *models.Task, userID string) (any, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.UserID = userID
	task.Status = models.TODO
	id, err := r.collection.InsertOne(ctx, task)
	return id.InsertedID, err
}

// 2️⃣ Obtener todas las tareas de un board
func (r *TaskRepository) GetTasksByBoardID(ctx context.Context, boardID string, page, limit int64) ([]models.Task, int64, error) {
	var tasks []models.Task

	skip := (page - 1) * limit
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := r.collection.Find(ctx, bson.M{"board_id": boardID}, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	count, err := r.collection.CountDocuments(ctx, bson.M{"board_id": boardID})
	if err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

// 3️⃣ Obtener una tarea específica
func (r *TaskRepository) GetTaskByID(ctx context.Context, taskID string) (*models.Task, error) {
	var task models.Task
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// 4️⃣ Actualizar tarea
func (r *TaskRepository) UpdateTask(ctx context.Context, taskID string, updatedData bson.M) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	updatedData["updated_at"] = time.Now()
	mongoResult, err := r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": updatedData})
	if mongoResult.MatchedCount == 0 {
		err = errors.New("id no encontrada")
	}
	return err
}

// 5️⃣ Eliminar tarea
func (r *TaskRepository) DeleteTask(ctx context.Context, taskID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	mongoResult, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if mongoResult.DeletedCount == 0 {
		err = errors.New("id no encontrada")
	}
	return err
}

// 6️⃣ Mover tarea
func (r *TaskRepository) UpdateTaskBoard(ctx context.Context, taskID, newBoardID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"board_id": newBoardID, "updated_at": time.Now()}}

	mongoResult, err := r.collection.UpdateOne(ctx, filter, update)
	if mongoResult.MatchedCount == 0 {
		err = errors.New("id no encontrada")
	}
	return err
}

// 7️⃣ Asignar una tarea
func (r *TaskRepository) UpdateTaskAssignee(ctx context.Context, taskID, userID string) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"assignee_id": userID, "updated_at": time.Now()}}

	mongoResult, err := r.collection.UpdateOne(ctx, filter, update)
	if mongoResult.MatchedCount == 0 {
		err = errors.New("id no encontrada")
	}
	return err
}

// 8️⃣ Cmbiar el estado de una tarea
func (r *TaskRepository) UpdateTaskStatus(ctx context.Context, taskID string, status models.TaskStatus) error {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}}

	mongoResult, err := r.collection.UpdateOne(ctx, filter, update)
	if mongoResult.MatchedCount == 0 {
		err = errors.New("id no encontrada")
	}
	return err
}

// 9️⃣ Guarda el usuario al recibir un evento de Kafka
func (r *TaskRepository) SaveUser(user *models.User) error {
	_, err := r.userCollection.InsertOne(context.Background(), user)
	return err
}

// 🔟 Obtiene el usuario por ID en la base de datos de mongo-task
func (r *TaskRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.userCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

// 1️⃣1️⃣ Guarda el tablero ar recibir un evento de Kafka
func (r *TaskRepository) SaveBoard(board *models.Board) error {
	_, err := r.boardCollection.InsertOne(context.Background(), board)
	return err
}

// 1️⃣2️⃣ Elimina un tablero al recibir un evento de Kafka
func (r *TaskRepository) DeleteBoard(board *models.Board) error {
	mongoResult, err := r.boardCollection.DeleteOne(context.Background(), board)
	if err != nil {
		return err
	}

	if mongoResult.DeletedCount == 0 {
		return errors.New("tablero no encontrado")
	}

	return nil
}

// 1️⃣3️⃣ Obtiene el tablero por ID en la base de datos de mongo-task
func (r *TaskRepository) GetBoardByID(id string) (*models.Board, error) {
	var board models.Board
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &board, err
	}

	err = r.boardCollection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&board)
	if err != nil {
		return &board, err
	}
	return &board, nil
}
