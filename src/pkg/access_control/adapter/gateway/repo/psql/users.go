package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) ListUsers() ([]entity.User, error) {
	const query = `
		SELECT 
			id,
			sir_name,
			first_name,
			last_name,
			user_type,
			gender,
			date_of_birth,
			created_at,
			updated_at
		FROM auth.users
	`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users %v", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var (
			id          sql.NullString
			sirName     sql.NullString
			firstName   sql.NullString
			lastName    sql.NullString
			userType    sql.NullString
			gender      sql.NullString
			dateOfBirth sql.NullTime
			createdAt   sql.NullTime
			updatedAt   sql.NullTime
		)

		if err := rows.Scan(
			&id,
			&sirName,
			&firstName,
			&lastName,
			&userType,
			&gender,
			&dateOfBirth,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user %v", err)
		}

		parsedID, err := uuid.Parse(id.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse UUID for user ID %v", err)
		}

		user := entity.User{
			Id:          parsedID,
			SirName:     sirName.String,
			FirstName:   firstName.String,
			LastName:    lastName.String,
			UserType:    userType.String,
			Gender:      entity.Gender(gender.String),
			DateOfBirth: parseNullTime(dateOfBirth),
			CreatedAt:   parseNullTime(createdAt),
			UpdatedAt:   parseNullTime(updatedAt),
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating over users %v", err)
	}

	return users, nil
}

func parseNullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}
