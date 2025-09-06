package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	IsEmailUnique(email string) bool
	AddUser(user model.User) error
	GetUserByEmailAndPassword(email string, password string) (*model.User, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateUser(user model.User) error
	ChangePassword(id uuid.UUID, newHashedPassword string) error
	DeleteUserByID(id uuid.UUID) error
}

type UserRepository struct {
	mu sync.Mutex
	db *sql.DB
}


func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(newUser model.User) error {

	query := `
		INSERT INTO users (id, firstname, middlename, lastname, mobile, email, password, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	r.mu.Lock()
	_, err := r.db.Exec(query, newUser.ID, newUser.Firstname, newUser.Middlename, newUser.Lastname, newUser.Mobile, newUser.Email, newUser.Password, newUser.Role)
	r.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return nil
}
