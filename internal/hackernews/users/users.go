package users

import (
	"context"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"

	database "github.com/captain-corgi/go-graphql-example/pkg/db/mysql"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

func (user *User) Create(ctx context.Context) {
	statement, err := database.Db.Prepare("INSERT INTO Users(Username,Password) VALUES(?,?)")
	print(statement)
	if err != nil {
		log.Fatal(err)
	}
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		// WARNING: Hashed password should not be failed!!!
		log.Printf("create hash failed: %s", err)
		hashedPassword = user.Password
	}
	_, err = statement.ExecContext(ctx, user.Username, hashedPassword)
	if err != nil {
		log.Fatal(err)
	}
}

// GetUserIdByUsername check if a user exists in database by given username
func GetUserIdByUsername(ctx context.Context, username string) (int, error) {
	statement, err := database.Db.Prepare("select ID from Users WHERE Username = ?")
	if err != nil {
		log.Fatal(err)
	}
	row := statement.QueryRowContext(ctx, username)

	var Id int
	err = row.Scan(&Id)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Print(err)
		}
		return 0, err
	}

	return Id, nil
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
