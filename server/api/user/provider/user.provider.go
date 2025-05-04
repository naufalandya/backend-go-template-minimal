package provider

import (
	"context"
	"errors"
	"modular_monolith/server/api/user/models"
	"modular_monolith/server/config/db"
)

func FetchAllUsers() ([]models.User, error) {
	query := `SELECT id, full_name, email FROM users`
	rows, err := db.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func FetchUserByID(id string) (*models.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = $1 LIMIT 1`
	row := db.DB.QueryRow(context.Background(), query, id)

	var user models.User
	err := row.Scan(&user.Id, &user.Name, &user.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func IsEmailOrUsernameExist(email, username string) (bool, error) {
	query := `SELECT COUNT(1) FROM users WHERE email = $1 OR username = $2`
	var count int
	err := db.DB.QueryRow(context.Background(), query, email, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CreateUser(hashedPassword string, input models.RegisterRequest) error {

	query := `
		INSERT INTO users (full_name, email, username, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err := db.DB.Exec(context.Background(), query, input.FullName, input.Email, input.Username, hashedPassword)
	return err
}
