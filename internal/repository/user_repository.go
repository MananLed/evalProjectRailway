package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/MananLed/evalProjectRailway/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

func(r *UserRepository) IsEmailUnique(email string) bool {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int 
	
	r.mu.Lock()
	err := r.db.QueryRow(query, email).Scan(&count)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false
	}

	return count == 0
}

func (r *UserRepository) GetUserByEmailAndPassword(email string, password string) (*model.User, error){
	query := `SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE email = $1`

	var user model.User

	r.mu.Lock()
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Firstname, &user.Middlename, &user.Lastname, &user.Mobile, &user.Email, &user.Password, &user.Role)
	r.mu.Unlock()

	if err != nil{
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("no user found")
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil{
		return nil, fmt.Errorf("invalid password")
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uuid.UUID) (*model.User, error) {
	query := `SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE id = $1`

	var user model.User 

	r.mu.Lock()
	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Firstname, &user.Middlename, &user.Lastname, &user.Mobile, &user.Email, &user.Password, &user.Role)
	r.mu.Unlock()

	if err != nil{
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found with id: %v", id)
		}
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err 
	}
	return &user, nil 
}

func (r *UserRepository) UpdateUser(user model.User) error{
	query := `UPDATE users 
				SET firstname = $1, middlename = $2, lastname = $3, mobile = $4, email = $5, role = $6
				WHERE id = $7`

	r.mu.Lock()
	_, err := r.db.Exec(query, user.Firstname, user.Middlename, user.Lastname, user.Mobile, user.Email, user.Role, user.ID)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	return nil 
}

func (r *UserRepository) ChangePassword(id uuid.UUID, newPassword string) error{
	query := `UPDATE users SET password = $1 WHERE id = $2`

	r.mu.Lock()
	_, err := r.db.Exec(query, newPassword, id)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	return nil
}

func (r *UserRepository) DeleteUserByID(id uuid.UUID) error{
	query := `DELETE FROM users WHERE id = $1`

	r.mu.Lock()
	_, err := r.db.Exec(query, id)
	r.mu.Unlock()

	if err != nil{
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err 
	}

	return nil 
}
