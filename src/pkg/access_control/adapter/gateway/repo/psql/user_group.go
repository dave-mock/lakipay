package repo

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) GrantPermissionToUser(userID uuid.UUID, permissionID uuid.UUID) error {
	const checkUserQuery = `
		SELECT 1 FROM auth.users WHERE id = $1
	`
	var userExists bool
	err := repo.db.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if user exists %v", err)
	}

	const checkPermissionQuery = `
		SELECT 1 FROM auth.permissions WHERE id = $1
	`
	var permissionExists bool
	err = repo.db.QueryRow(checkPermissionQuery, permissionID).Scan(&permissionExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("permission with ID %s does not exist", permissionID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if permission exists %v", err)
	}

	const checkPermissionGrantedQuery = `
		SELECT 1 FROM auth.user_permissions
		WHERE user_id = $1 AND permission_id = $2
	`
	var exists bool
	err = repo.db.QueryRow(checkPermissionGrantedQuery, userID, permissionID).Scan(&exists)
	if err == sql.ErrNoRows {
	} else if err != nil {
		return fmt.Errorf("failed to check permission existence %v", err)
	} else {
		return fmt.Errorf("permission with ID %s is already granted to user with ID %s", permissionID, userID)
	}

	const grantPermissionQuery = `
		INSERT INTO auth.user_permissions (id, user_id, permission_id, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = repo.db.Exec(grantPermissionQuery, userID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to grant permission to user %v", err)
	}

	return nil
}

func (repo PsqlRepo) RevokePermissionFromUser(userID uuid.UUID, permissionID uuid.UUID) error {
	if _, err := uuid.Parse(userID.String()); err != nil {
		return fmt.Errorf("invalid user ID format")
	}
	if _, err := uuid.Parse(permissionID.String()); err != nil {
		return fmt.Errorf("invalid permission ID format")
	}

	const checkUserQuery = `
        SELECT 1 FROM auth.users WHERE id = $1
    `
	var userExists bool
	err := repo.db.QueryRow(checkUserQuery, userID).Scan(&userExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user with ID %v does not exist", userID)
		}
		return fmt.Errorf("failed to check user existence: %v", err)
	}

	const checkPermissionQuery = `
        SELECT 1 FROM auth.permissions WHERE id = $1
    `
	var permissionExists bool
	err = repo.db.QueryRow(checkPermissionQuery, permissionID).Scan(&permissionExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("permission with ID %v does not exist", permissionID)
		}
		return fmt.Errorf("failed to check permission existence: %v", err)
	}

	const checkUserPermissionQuery = `
        SELECT 1 FROM auth.user_permissions
        WHERE user_id = $1 AND permission_id = $2
    `
	var userPermissionExists bool
	err = repo.db.QueryRow(checkUserPermissionQuery, userID, permissionID).Scan(&userPermissionExists)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("permission does not exist for user with ID %v", userID)
		}
		return fmt.Errorf("failed to check user-permission existence %v", err)
	}

	const revokeQuery = `
        DELETE FROM auth.user_permissions
        WHERE user_id = $1 AND permission_id = $2
    `
	_, err = repo.db.Exec(revokeQuery, userID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission from user %v", err)
	}

	return nil
}

func (repo PsqlRepo) GrantPermissionToGroup(groupID uuid.UUID, permissionID uuid.UUID) error {
	const checkGroupQuery = `
		SELECT 1 FROM auth.groups WHERE id = $1
	`
	var groupExists bool
	err := repo.db.QueryRow(checkGroupQuery, groupID).Scan(&groupExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("group with ID %s does not exist", groupID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if group exists %v", err)
	}

	const checkPermissionQuery = `
		SELECT 1 FROM auth.permissions WHERE id = $1
	`
	var permissionExists bool
	err = repo.db.QueryRow(checkPermissionQuery, permissionID).Scan(&permissionExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("permission with ID %s does not exist", permissionID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if permission exists %v", err)
	}

	const checkPermissionGrantedQuery = `
		SELECT 1 FROM auth.group_permissions
		WHERE group_id = $1 AND permission_id = $2
	`
	var exists bool
	err = repo.db.QueryRow(checkPermissionGrantedQuery, groupID, permissionID).Scan(&exists)
	if err == sql.ErrNoRows {
	} else if err != nil {
		return fmt.Errorf("failed to check permission existence %v", err)
	} else {
		return fmt.Errorf("permission with ID %s is already granted to group with ID %s", permissionID, groupID)
	}

	const grantPermissionQuery = `
		INSERT INTO auth.group_permissions (group_id, permission_id, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	_, err = repo.db.Exec(grantPermissionQuery, groupID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to grant permission to group %v", err)
	}
	return nil
}

func (repo PsqlRepo) RevokePermissionFromGroup(groupID uuid.UUID, permissionID uuid.UUID) error {
	const checkGroupQuery = `
		SELECT 1 FROM auth.groups WHERE id = $1
	`
	var groupExists bool
	err := repo.db.QueryRow(checkGroupQuery, groupID).Scan(&groupExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("group with ID %s does not exist", groupID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if group exists %v", err)
	}

	const checkPermissionQuery = `
		SELECT 1 FROM auth.permissions WHERE id = $1
	`
	var permissionExists bool
	err = repo.db.QueryRow(checkPermissionQuery, permissionID).Scan(&permissionExists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("permission with ID %s does not exist", permissionID)
	}
	if err != nil {
		return fmt.Errorf("failed to check if permission exists %v", err)
	}

	const checkPermissionAssignmentQuery = `
		SELECT 1 FROM auth.group_permissions
		WHERE group_id = $1 AND permission_id = $2
	`
	var exists bool
	err = repo.db.QueryRow(checkPermissionAssignmentQuery, groupID, permissionID).Scan(&exists)
	if err == sql.ErrNoRows {
		return fmt.Errorf("permission with ID %s is not assigned to group with ID %s", permissionID, groupID)
	}
	if err != nil {
		return fmt.Errorf("failed to check permission assignment %v", err)
	}

	const query = `
		DELETE FROM auth.group_permissions
		WHERE group_id = $1 AND permission_id = $2
	`
	_, err = repo.db.Exec(query, groupID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission from group %v", err)
	}
	return nil
}
