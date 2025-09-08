package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/MananLed/evalProjectRailway/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestAddUser_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	user := model.User{ID: uuid.New(), Firstname: "John", Middlename: "M", Lastname: "Doe", Mobile: "123", Email: "john@example.com", Password: "hashed", Role: "passenger"}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(user.ID, user.Firstname, user.Middlename, user.Lastname, user.Mobile, user.Email, user.Password, user.Role).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.AddUser(user); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAddUser_Failure(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	user := model.User{ID: uuid.New()}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO users`)).
		WithArgs(user.ID, "", "", "", "", "", "", "").
		WillReturnError(errors.New("insert failed"))

	if err := repo.AddUser(user); err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestIsEmailUnique_True(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM users WHERE email = $1`)).
		WithArgs("unique@example.com").
		WillReturnRows(rows)

	if !repo.IsEmailUnique("unique@example.com") {
		t.Errorf("expected true")
	}
}

func TestIsEmailUnique_False(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM users WHERE email = $1`)).
		WithArgs("exists@example.com").
		WillReturnRows(rows)

	if repo.IsEmailUnique("exists@example.com") {
		t.Errorf("expected false")
	}
}

func TestIsEmailUnique_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(*) FROM users WHERE email = $1`)).
		WithArgs("bad@example.com").
		WillReturnError(errors.New("db error"))

	if repo.IsEmailUnique("bad@example.com") {
		t.Errorf("expected false due to error")
	}
}

func TestGetUserByEmailAndPassword_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	email, password := "john@example.com", "secret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	rows := sqlmock.NewRows([]string{"id", "firstname", "middlename", "lastname", "mobile", "email", "password", "role"}).
		AddRow(uuid.New(), "John", "M", "Doe", "123", email, string(hashed), "passenger")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnRows(rows)

	if _, err := repo.GetUserByEmailAndPassword(email, password); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetUserByEmailAndPassword_InvalidPassword(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	email := "john@example.com"
	hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)

	rows := sqlmock.NewRows([]string{"id", "firstname", "middlename", "lastname", "mobile", "email", "password", "role"}).
		AddRow(uuid.New(), "John", "M", "Doe", "123", email, string(hashed), "passenger")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnRows(rows)

	if _, err := repo.GetUserByEmailAndPassword(email, "wrong"); err == nil {
		t.Errorf("expected invalid password error")
	}
}

func TestGetUserByEmailAndPassword_NoUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	email := "nouser@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE email = $1`)).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	if _, err := repo.GetUserByEmailAndPassword(email, "pw"); err == nil {
		t.Errorf("expected no user error")
	}
}

func TestGetUserByEmailAndPassword_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE email = $1`)).
		WithArgs("err@example.com").
		WillReturnError(errors.New("db error"))

	if _, err := repo.GetUserByEmailAndPassword("err@example.com", "pw"); err == nil {
		t.Errorf("expected error")
	}
}

func TestGetUserByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "firstname", "middlename", "lastname", "mobile", "email", "password", "role"}).
		AddRow(id, "John", "M", "Doe", "123", "john@example.com", "hashed", "passenger")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	if _, err := repo.GetUserByID(id); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetUserByID_NoUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	if _, err := repo.GetUserByID(id); err == nil {
		t.Errorf("expected no user error")
	}
}

func TestGetUserByID_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, firstname, middlename, lastname, mobile, email, password, role FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(errors.New("db error"))

	if _, err := repo.GetUserByID(id); err == nil {
		t.Errorf("expected error")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	user := model.User{ID: uuid.New(), Firstname: "John", Middlename: "M", Lastname: "Doe", Mobile: "123", Email: "john@example.com", Role: "passenger"}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET firstname = $1, middlename = $2, lastname = $3, mobile = $4, email = $5, role = $6 WHERE id = $7`)).
		WithArgs(user.Firstname, user.Middlename, user.Lastname, user.Mobile, user.Email, user.Role, user.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.UpdateUser(user); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateUser_Failure(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	user := model.User{ID: uuid.New()}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users`)).
		WithArgs("", "", "", "", "", "", user.ID).
		WillReturnError(errors.New("update failed"))

	if err := repo.UpdateUser(user); err == nil {
		t.Errorf("expected error")
	}
}

func TestChangePassword_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET password = $1 WHERE id = $2`)).
		WithArgs("newpw", id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.ChangePassword(id, "newpw"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestChangePassword_Failure(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET password = $1 WHERE id = $2`)).
		WithArgs("newpw", id).
		WillReturnError(errors.New("update failed"))

	if err := repo.ChangePassword(id, "newpw"); err == nil {
		t.Errorf("expected error")
	}
}

func TestDeleteUserByID_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := repo.DeleteUserByID(id); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDeleteUserByID_Failure(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	id := uuid.New()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(id).
		WillReturnError(errors.New("delete failed"))

	if err := repo.DeleteUserByID(id); err == nil {
		t.Errorf("expected error")
	}
}
