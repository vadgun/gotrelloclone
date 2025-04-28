package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vadgun/gotrelloclone/user-service/infra/config"
	"github.com/vadgun/gotrelloclone/user-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// UserRepository maneja las operaciones con MongoDB.
type UserRepository struct {
	Collection  *mongo.Collection
	RedisClient *redis.Client
	Logger      *zap.Logger
}

// NewUserRepository crea una nueva instancia del repositorio.
func NewUserRepository(collection *mongo.Collection, redisClient *redis.Client, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		Collection:  collection,
		RedisClient: redisClient,
		Logger:      logger,
	}
}

// CreateUser guarda un usuario en la base de datos.
func (r *UserRepository) CreateUser(user *models.User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now()
	mongoResult, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	oid, ok := mongoResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("no se pudo obtener el ID del usuario")
	}

	id := oid.Hex()

	return id, nil
}

// GetUserByEmail busca un usuario por su email.
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 1. Buscar en Redis
	val, err := r.RedisClient.Get(ctx, "user_email:"+email).Result()
	if err == nil {
		var cachedUser models.User
		if json.Unmarshal([]byte(val), &cachedUser) == nil {
			r.Logger.Info("✅ User by email from cache")
			return &cachedUser, nil
		}
	}

	// 2. Si no está, ir a Mongo
	var user models.User
	err = r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	// 3. Guardar en cache
	bytes, _ := json.Marshal(user)
	r.RedisClient.Set(ctx, "user_email:"+email, bytes, 10*time.Minute)

	r.Logger.Info("📥 User by email cached")

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
	// 1. Revisar cache
	val, err := config.RedisClient.Get(ctx, userID).Result()
	if err == nil {
		var cachedUser models.User
		if err := json.Unmarshal([]byte(val), &cachedUser); err == nil {
			r.Logger.Info("✅ UserByID from cache")
			return &cachedUser, nil
		}
	}

	// 2. Ir a Mongo si no está en cache
	var user models.User
	err = r.Collection.FindOne(ctx, bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// 3. Guardar en cache
	bytes, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, userID, bytes, 10*time.Minute) // TTL configurable
	r.Logger.Info("📥 UserByID cached")

	return &user, nil
}

// GetAllUsers busca todos los usuarios en la base de datos
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
