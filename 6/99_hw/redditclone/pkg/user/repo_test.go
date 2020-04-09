package user

import (
	"fmt"
	"redditclone/pkg/random"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAuthorize(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	if err != nil {
		t.Errorf("Cant create mock %v", err)
		return
	}
	mockUserRepo := NewRepo(db)

	// User not registered
	username := "user"
	password := "pass"
	row := sqlmock.NewRows([]string{"id", "username", "password"})
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnError(fmt.Errorf("not found"))
	mockUserRepo.Authorize(username, password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}

	// User wrong password
	username = "user"
	password = "pass"
	otherpassword := "notpass"
	row = sqlmock.NewRows([]string{"id", "username", "password"})
	row = row.AddRow("1", username, otherpassword)
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(row)
	mockUserRepo.Authorize(username, password)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}

	// OK
	username = "user"
	password = "pass"
	row = sqlmock.NewRows([]string{"id", "username", "password"})
	row = row.AddRow("1", username, password)
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(row)
	mockUserRepo.Authorize(username, password)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}
}

func TestNewUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	if err != nil {
		t.Errorf("Cant create mock %v", err)
		return
	}
	mockUserRepo := NewRepo(db)
	// mocking random
	mockUserRepo.rn = random.NewGenerator(false)

	// User already registered
	username := "user"
	password := "pass"
	row := sqlmock.NewRows([]string{"id", "username", "password"})
	row = row.AddRow("1", username, password)
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnRows(row)
	mockUserRepo.NewUser(username, password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}

	// OK
	username = "user"
	password = "pass"
	row = sqlmock.NewRows([]string{"id", "username", "password"})
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE username = ?").
		WithArgs(username).
		WillReturnError(fmt.Errorf("not found"))
	mock.
		ExpectExec("INSERT INTO users").
		WithArgs("1", username, password).
		WillReturnError(nil)
	mockUserRepo.NewUser(username, password)
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}
}

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close()
	if err != nil {
		t.Errorf("Cant create mock %v", err)
		return
	}
	mockUserRepo := NewRepo(db)

	// User not found
	username := "user"
	password := "pass"
	row := sqlmock.NewRows([]string{"id", "username", "password"})
	row = row.AddRow("1", username, password)
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE id = ?").
		WithArgs("1").
		WillReturnRows(row)
	mockUserRepo.GetUserById("1")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}

	// OK
	username = "user"
	password = "pass"
	mock.
		ExpectQuery("SELECT id, username, password FROM users WHERE id = ?").
		WithArgs("1").
		WillReturnError(fmt.Errorf("user not found"))
	mockUserRepo.GetUserById("1")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
		return
	}
}

// go test -v -coverprofile=repo.out && go tool cover -html=repo.out -o repo.html && rm repo.out
