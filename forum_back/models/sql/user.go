package sqls

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"awesome_forum/forum_back/models"
)

type UserModel struct {
	DB *sql.DB
}

func (um *UserModel) Insert(name, email, password string) error {
	_, er := um.GetByEmail(email)

	if er == nil {
		return models.ErrDuplicateEmail //user exists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(username, email, password) VALUES (?, ?, ?)`

	_, err = um.DB.Exec(stmt, name, email, string(hashedPassword))

	return err
}

func (um *UserModel) Authenticate(email, password string) (*models.User, error) {
	user, err := um.GetByEmail(email)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, models.ErrInvalidCredentials
		} else {
			return nil, err
		}
	}
	user.Password = ""
	return user, nil
}

func (um *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}

func (um *UserModel) GetByEmail(email string) (*models.User, error) {
	stmt := `SELECT rowid, username, email, password FROM users WHERE email = ?`
	row := um.DB.QueryRow(stmt, email)
	u := &models.User{}
	err := row.Scan(&u.UserId, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}
