package repo

import (
	"auth/src/pkg/access_control/core/entity"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateGroup(title string) (*entity.Group, error) {
	const query = `
		INSERT INTO auth.groups (title)
		VALUES ($1)
		RETURNING id, title, created_at, updated_at
	`
	var group entity.Group
	err := repo.db.QueryRow(query, title).Scan(&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create group %v", err)
	}
	return &group, nil
}

func (repo PsqlRepo) UpdateGroup(groupID uuid.UUID, title string) (*entity.Group, error) {
	const query = `
		UPDATE auth.groups
		SET title = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING id, title, created_at, updated_at
	`
	var group entity.Group
	err := repo.db.QueryRow(query, title, groupID).Scan(&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update group %v", err)
	}
	return &group, nil
}

func (repo PsqlRepo) DeleteGroup(groupID uuid.UUID) error {
	const query = `
		DELETE FROM auth.groups
		WHERE id = $1
	`
	_, err := repo.db.Exec(query, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete group %v", err)
	}
	return nil
}

func (repo PsqlRepo) ListGroups() ([]entity.Group, error) {
	const query = `
		SELECT g.id, g.title, g.created_at, g.updated_at,
		       COUNT(ug.user_id) AS member_count
		FROM auth.groups g
		LEFT JOIN auth.user_groups ug ON g.id = ug.group_id
		GROUP BY g.id
	`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list groups %v", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var group entity.Group
		var memberCount int
		if err := rows.Scan(&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt, &memberCount); err != nil {
			return nil, fmt.Errorf("failed to scan group %v", err)
		}

		group.Members = memberCount
		groups = append(groups, group)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("no groups found")
	}

	return groups, nil
}

func (repo PsqlRepo) AddUserToGroup(userID, groupID uuid.UUID) error {
	const userExistsQuery = `
		SELECT 1 FROM auth.users WHERE id = $1
	`
	var userExists bool
	err := repo.db.QueryRow(userExistsQuery, userID).Scan(&userExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user with ID %v does not exist", userID)
	} else if err != nil {
		return fmt.Errorf("failed to check if user exists %v", err)
	}

	const groupExistsQuery = `
		SELECT 1 FROM auth.groups WHERE id = $1
	`
	var groupExists bool
	err = repo.db.QueryRow(groupExistsQuery, groupID).Scan(&groupExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("group with ID %v does not exist", groupID)
	} else if err != nil {
		return fmt.Errorf("failed to check if group exists: %v", err)
	}

	const userInGroupQuery = `
		SELECT 1 FROM auth.user_groups WHERE user_id = $1 AND group_id = $2
	`
	var userInGroup bool
	err = repo.db.QueryRow(userInGroupQuery, userID, groupID).Scan(&userInGroup)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check if user is already in group %v", err)
	}
	if userInGroup {
		return fmt.Errorf("user with ID %v is already in group with ID %v", userID, groupID)
	}

	// Add user to group
	const insertQuery = `
		INSERT INTO auth.user_groups (user_id, group_id)
		VALUES ($1, $2)
	`
	_, err = repo.db.Exec(insertQuery, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to add user to group: %v", err)
	}

	return nil
}

func (repo PsqlRepo) RemoveUserFromGroup(userID, groupID uuid.UUID) error {
	const checkAndDeleteQuery = `
		DELETE FROM auth.user_groups
		WHERE user_id = $1 AND group_id = $2
		RETURNING user_id, group_id
	`

	var returnedUserID, returnedGroupID uuid.UUID
	err := repo.db.QueryRow(checkAndDeleteQuery, userID, groupID).Scan(&returnedUserID, &returnedGroupID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user is not a member of the specified group")
		}
		return fmt.Errorf("an unexpected error occurred while removing the user from the group")
	}

	return nil
}

func (repo PsqlRepo) ListGroupUsers(groupID uuid.UUID) ([]entity.User, error) {
	const groupExistsQuery = `
		SELECT 1 FROM auth.groups WHERE id = $1
	`
	var groupExists bool
	err := repo.db.QueryRow(groupExistsQuery, groupID).Scan(&groupExists)
	if err != nil || !groupExists {
		if err == nil {
			err = fmt.Errorf("group with ID %v does not exist", groupID)
		}
		return nil, err
	}

	const query = `
		SELECT u.id, u.sir_name, u.first_name, u.last_name, u.gender, u.date_of_birth, u.user_type, u.created_at, u.updated_at
		FROM auth.user_groups ug
		JOIN auth.users u ON ug.user_id = u.id
		WHERE ug.group_id = $1
	`
	rows, err := repo.db.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list group users %v", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		var gender sql.NullString
		var dateOfBirth sql.NullTime

		if err := rows.Scan(
			&user.Id,
			&user.SirName,
			&user.FirstName,
			&user.LastName,
			&gender,
			&dateOfBirth,
			&user.UserType,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}

		user.Gender = entity.Gender(gender.String)
		if !gender.Valid {
			user.Gender = ""
		}

		if dateOfBirth.Valid {
			user.DateOfBirth = dateOfBirth.Time
		} else {
			user.DateOfBirth = time.Time{}
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no users found in group with ID %v", groupID)
	}

	return users, nil
}

func (repo PsqlRepo) ListUserGroups(userID uuid.UUID) ([]entity.Group, error) {

	const userExistsQuery = `
		SELECT 1 FROM auth.users WHERE id = $1
	`
	var userExists bool
	err := repo.db.QueryRow(userExistsQuery, userID).Scan(&userExists)
	if err != nil || !userExists {
		if err == nil {
			err = fmt.Errorf("user with ID %v does not exist", userID)
		}
		return nil, err
	}

	const query = `
		SELECT g.id, g.title, g.created_at, g.updated_at
		FROM auth.user_groups ug
		JOIN auth.groups g ON ug.group_id = g.id
		WHERE ug.user_id = $1
	`
	rows, err := repo.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user groups %v", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var group entity.Group
		if err := rows.Scan(&group.ID, &group.Title, &group.CreatedAt, &group.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group %v", err)
		}
		groups = append(groups, group)
	}

	if len(groups) == 0 {
		return nil, fmt.Errorf("user with ID %v is not a member of any group", userID)
	}

	return groups, nil
}
