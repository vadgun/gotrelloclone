package services

import (
	"time"

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

func (s *BoardService) CreateBoard(name, ownerID string) (*models.Board, error) {
	board := &models.Board{
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
	}

	id, err := s.repo.CreateBoard(board)
	if err != nil {
		return nil, err
	}

	board.ID = id.(primitive.ObjectID)

	return board, nil
}

func (s *BoardService) GetBoardsByUser(userID string) ([]models.Board, error) {
	return s.repo.GetBoardsByUser(userID)
}

func (s *BoardService) GetBoardByID(boardID string) (*models.Board, error) {
	return s.repo.GetBoardByID(boardID)
}
