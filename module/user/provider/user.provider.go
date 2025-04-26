package provider

import (
	"context"
	"errors"
	"modular_monolith/module/user/models"
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
