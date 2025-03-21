package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vadgun/gotrelloclone/user-service/config"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository maneja las operaciones con MongoDB.
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository crea una nueva instancia del repositorio.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: config.DB.Collection("users"),
	}
}

// CreateUser guarda un usuario en la base de datos.
func (r *UserRepository) CreateUser(user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	mongoResult, err := r.collection.InsertOne(ctx, user)
	id := mongoResult.InsertedID.(primitive.ObjectID).Hex()
	logrus.WithFields(logrus.Fields{
		"user_id": id,
	}).Info("Guardando usuario en la base de datos")
	return id, err
}

// GetUserByEmail busca un usuario por su email.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID busca un usuario por su ID en la base de datos.
func (r *UserRepository) GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userObjID, errs := primitive.ObjectIDFromHex(userID)
	if errs != nil {
		return nil, errors.New("id no valida")
	}

	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}
