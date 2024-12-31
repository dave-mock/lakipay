package repo

import (
	"auth/src/pkg/auth/core/entity"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Pre Session
func (repo PsqlRepo) StorePreSession(preSession entity.PreSession) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.pre_sessions (id, token, created_at)
	VALUES ($1::UUID, $2, $3);
	`, repo.schema), preSession.Id, sql.NullString{Valid: preSession.Token != "", String: preSession.Token}, preSession.CreatedAt)

	return err
}

func (repo PsqlRepo) UpdatePreSession(id string, token string) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.pre_sessions
	SET hash = $2
	WHERE id = $1
	RETURNING id, hash;
	`, repo.schema), id, token)

	return err
}

// Device

func (repo PsqlRepo) StoreDevice(device entity.Device) error {
	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.devices (id, ip, name, agent, created_at)
	VALUES ($1::UUID,$2,$3, $4, $5);
	`, repo.schema), device.Id, device.IP.String(), device.Name, device.Agent, device.CreatedAt)

	return err
}

func (repo PsqlRepo) StoreDeviceAuth(deviceAuth entity.DeviceAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.device_auths (id, token, device, status, created_at)
	VALUES($1::UUID,$2, $3, $4, $5);
	`, repo.schema), deviceAuth.Id, deviceAuth.Token, deviceAuth.Device.Id, deviceAuth.Status, deviceAuth.CreatedAt)

	return err
}

func (repo PsqlRepo) UpdateDeviceAuthStatus(deviceAuthId uuid.UUID, status bool) error {

	_, err := repo.db.Query(fmt.Sprintf(`
	UPDATE %s.device_auths
	SET status = $2
	WHERE id = $1::UUID;
	`, repo.schema), deviceAuthId, status)

	return err
}

func (repo PsqlRepo) FindDeviceAuth(token string) (entity.DeviceAuth, error) {
	var deviceAuth entity.DeviceAuth

	var ip sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT device_auths.id, device_auths.token, devices.id, devices.ip, devices.name, devices.agent ,device_auths.status
	FROM %s.device_auths
	INNER JOIN %s.devices ON %s.devices.id = device_auths.device
	WHERE device_auths.token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&deviceAuth.Id, &deviceAuth.Token,
		&deviceAuth.Device.Id, &ip, &deviceAuth.Device.Name, &deviceAuth.Device.Agent,
		&deviceAuth.Status,
	)

	return deviceAuth, err
}

// Phone

func (repo PsqlRepo) StorePhone(phone entity.Phone) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phones (id, prefix, number, created_at)
	VALUES ($1::UUID,$2,$3,$4);
	`, repo.schema), phone.Id, phone.Prefix, phone.Number, phone.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPhone(prefix string, number string) (*entity.Phone, error) {

	var phone entity.Phone
	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT id, prefix, number 
	FROM %s.phones
	WHERE prefix = $1 AND number = $2;
	`, repo.schema), prefix, number).Scan(
		&phone.Id, &phone.Prefix, &phone.Number,
	)

	repo.log.Println("Find phone")

	if err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			{
				return nil, nil
			}
		}
		return nil, err
	}

	return &phone, err
}

func (repo PsqlRepo) StorePhoneAuth(phoneAuth entity.PhoneAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.phone_auths (id, token, phone_id, code, method, status, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4, $5, $6, $7);
	`, repo.schema), phoneAuth.Id, phoneAuth.Token, phoneAuth.Phone.Id, phoneAuth.Code, phoneAuth.Method, phoneAuth.Status, time.Now())

	return err
}

func (repo PsqlRepo) FindPhoneAuth(token string) (entity.PhoneAuth, error) {

	var phoneAuth entity.PhoneAuth

	fmt.Println("||||||||||||||||||| ", token)

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT phone_auths.id, phone_auths.token, 
		phones.id, phones.prefix, phones.number,
		phone_auths.code, phone_auths.method, phone_auths.status
	FROM %s.phone_auths
	INNER JOIN %s.phones ON %s.phones.id = phone_auths.phone_id
	WHERE token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&phoneAuth.Id, &phoneAuth.Token,
		&phoneAuth.Phone.Id, &phoneAuth.Phone.Prefix, &phoneAuth.Phone.Number,
		&phoneAuth.Code, &phoneAuth.Method, &phoneAuth.Status,
	)

	return phoneAuth, err
}

func (repo PsqlRepo) FindPhoneAuthWithoutPhone(token string) (entity.PhoneAuth, error) {

	var phoneAuth entity.PhoneAuth

	fmt.Println("||||||||||||||||||| ", token)

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT phone_auths.id, phone_auths.token, 
		phone_auths.code, phone_auths.method, phone_auths.status
	FROM %s.phone_auths
	WHERE token = $1
	order by created_at DESC
	limit 1
	`, repo.schema), token).Scan(
		&phoneAuth.Id, &phoneAuth.Token,
		&phoneAuth.Code, &phoneAuth.Method, &phoneAuth.Status,
	)

	return phoneAuth, err
}

func (repo PsqlRepo) UpdatePhoneAuthStatus(phoneAuthId uuid.UUID, status bool) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	UPDATE %s.phone_auths
	SET status = $2
	WHERE id = $1::UUID
	`, repo.schema), phoneAuthId, status)
	return err
}

func (repo PsqlRepo) StoreSession(session entity.Session) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.sessions (id, token, user_id, device_id, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4::UUID, $5);
	`, repo.schema), session.Id, session.Token, session.User.Id, session.Device.Id, session.CreatedAt)

	return err
}

func (repo PsqlRepo) FindSessionById(id uuid.UUID) (*entity.Session, error) {
	var session entity.Session

	var sirName sql.NullString
	var lastName sql.NullString
	var userType sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT sessions.id, sessions.token, users.id, users.sir_name, users.first_name, users.last_name, users.user_type
	FROM %s.sessions
	INNER JOIN %s.users ON %s.users.id = sessions.user_id
	WHERE sessions.id = $1::UUID
	`, repo.schema, repo.schema, repo.schema), id).Scan(
		&session.Id, &session.Token,
		&session.User.Id, &sirName, &session.User.FirstName, &lastName, &userType,
	)

	if sirName.Valid {
		session.User.SirName = sirName.String
	}

	if lastName.Valid {
		session.User.LastName = lastName.String
	}

	if userType.Valid {
		session.User.UserType = userType.String
	} else {
		session.User.UserType = "UNKNOWN"
	}

	return &session, err
}

func (repo PsqlRepo) StorePasswordAuth(passwordAuth entity.PasswordAuth) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.password_auths (id, token, password_id, status, created_at)
	VALUES ($1::UUID, $2, $3::UUID, $4, $5);
	`, repo.schema), passwordAuth.Id, passwordAuth.Token, passwordAuth.Password.Id, passwordAuth.Status, passwordAuth.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPasswordAuth(token string) (*entity.PasswordAuth, error) {

	var passwordAuth entity.PasswordAuth
	var hint sql.NullString

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT password_auths.id, password_auths.token, 
		password_identities.id, password_identities.password, password_identities.hint,
		password_auths.status, password_auths.created_at, password_auths.updated_at
	FROM %s.password_auths
	INNER JOIN %s.password_identities ON %s.password_identities.id = password_auths.password_id
	WHERE token = $1
	`, repo.schema, repo.schema, repo.schema), token).Scan(
		&passwordAuth.Id, &passwordAuth.Token,
		&passwordAuth.Password.Id, &passwordAuth.Password.Password, &hint,
		&passwordAuth.Status, &passwordAuth.CreatedAt, &passwordAuth.UpdatedAt,
	)

	if hint.Valid {
		passwordAuth.Password.Hint = hint.String
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
	}

	return &passwordAuth, err
}

func (repo PsqlRepo) CheckPermission(userID uuid.UUID, requiredPermission entity.Permission) (bool, error) {
	log.Printf("Checking permission for userID |||||||||   ||||||||||| %s with payload: %+v\n", userID, requiredPermission)

	query := `
        WITH create_operation AS (
            SELECT id
            FROM auth.operations
            WHERE name = $2
        ),
        user_resource AS (
            SELECT id
            FROM auth.resources
            WHERE name = $3
        )
        SELECT p.resource, p.resource_id, o.name AS operation, p.effect
        FROM auth.permissions p
        JOIN auth.user_permissions up ON up.permission_id = p.id
        JOIN user_resource ur ON ur.id = p.resource
        LEFT JOIN auth.operations o ON o.id = ANY(p.operations)
        WHERE up.user_id = $1
        AND EXISTS (
            SELECT 1
            FROM create_operation co
            WHERE co.id = ANY(p.operations)
        )
        AND p.effect = $4
        AND p.resource_id = $5;
    `

	log.Printf("Query to fetch user permissions |||||||   ||||||| %s\n", query)

	rowsUserPermissions, err := repo.db.Query(query, userID, requiredPermission.Operation, requiredPermission.Resource, requiredPermission.Effect, requiredPermission.ResourceIdentifier)
	if err != nil {
		log.Printf("Failed to fetch user permissions ||||||| ||||||||| %v\n", err)
		return false, fmt.Errorf("failed to fetch user permissions ||||| |||||||| %v", err)
	}
	defer rowsUserPermissions.Close()

	var permissionFound bool
	for rowsUserPermissions.Next() {
		var permission entity.Permission
		if err := rowsUserPermissions.Scan(&permission.Resource, &permission.ResourceIdentifier, &permission.Operation, &permission.Effect); err != nil {
			log.Printf("Failed to scan permission ||||||| || %v\n", err)
			return false, fmt.Errorf("failed to scan permission ||||||||      |||||||| %v", err)
		}

		log.Printf("Found permission |||||||| %+v\n", permission)

		if permission.Effect == requiredPermission.Effect {
			log.Printf("User %s is %s to perform the operation.\n", userID, requiredPermission.Effect)
			permissionFound = true
			break
		}
	}

	if permissionFound {
		return true, nil
	}

	queryGroups := `
        SELECT g.id
        FROM auth.user_groups ug
        JOIN auth.groups g ON ug.group_id = g.id
        WHERE ug.user_id = $1
    `
	log.Printf("Query to fetch user groups ||||||||||||||||  |||||||||||||| %s\n", queryGroups)
	rowsGroups, err := repo.db.Query(queryGroups, userID)
	if err != nil {
		log.Printf("Failed to fetch user groups ||||||||||||||||  |||||||||||||| %v\n", err)
		return false, fmt.Errorf("failed to fetch user groups ||||||||||||||||  |||||||||||||| %v", err)
	}
	defer rowsGroups.Close()

	var groups []uuid.UUID
	for rowsGroups.Next() {
		var groupID uuid.UUID
		if err := rowsGroups.Scan(&groupID); err != nil {
			log.Printf("Failed to scan group ID ||||||||||||||||  |||||||||||||| %v\n", err)
			return false, fmt.Errorf("failed to scan group ID ||||||||||||||||  |||||||||||||| %v", err)
		}
		groups = append(groups, groupID)
	}

	if len(groups) == 0 {
		log.Printf("User %s does not belong to any group\n", userID)
		return false, fmt.Errorf("user does not belong to any group")
	}

	groupIDs := pq.Array(groups)
	queryGroupPermissions := `
        SELECT p.resource, p.resource_id, o.name AS operation, p.effect
        FROM auth.group_permissions gp
        JOIN auth.permissions p ON gp.permission_id = p.id
        LEFT JOIN auth.operations o ON o.id = ANY(p.operations)
        WHERE gp.group_id = ANY($1)
        AND p.resource = $2
        AND o.name = $3
        AND (p.resource_id = $4 OR p.resource_id = '*')
        AND p.effect = $5
    `
	log.Printf("Query to fetch group permissions|||||||||||||||| %s\n", queryGroupPermissions)
	rowsGroupPermissions, err := repo.db.Query(queryGroupPermissions, groupIDs, requiredPermission.Resource, requiredPermission.Operation, requiredPermission.ResourceIdentifier, requiredPermission.Effect)
	if err != nil {
		log.Printf("Failed to fetch group permissions |||||||||||||| %v\n", err)
		return false, fmt.Errorf("failed to fetch group permissions |||||||||||||| %v", err)
	}
	defer rowsGroupPermissions.Close()

	for rowsGroupPermissions.Next() {
		var permission entity.Permission
		if err := rowsGroupPermissions.Scan(&permission.Resource, &permission.ResourceIdentifier, &permission.Operation, &permission.Effect); err != nil {
			log.Printf("Failed to scan permission |||||||||||||| %v\n", err)
			return false, fmt.Errorf("failed to scan permission |||||||||||||| %v", err)
		}
		log.Printf("Found group permission |||||||||||||| %+v\n", permission)

		if permission.Effect == requiredPermission.Effect {
			log.Printf("User %s is %s to perform the operation via group permission.\n", userID, requiredPermission.Effect)
			return true, nil
		}
	}

	log.Printf("User %s does not have the required permission.\n", userID)
	return false, fmt.Errorf("user does not have the required permission")
}

func (repo PsqlRepo) FindUserPermissions(userID uuid.UUID, requiredPermission entity.Permission) ([]entity.Permission, error) {
	var permissions []entity.Permission
	query := `
		SELECT resource, resource_identifier, operation, effect
		FROM auth.user_permissions
		WHERE user_id = $1
		AND resource = $2
		AND operation = $3
		AND resource_identifier = $4
	`
	rows, err := repo.db.Query(query, userID, requiredPermission.Resource, requiredPermission.Operation, requiredPermission.ResourceIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user permissions: %v", err)
	}
	defer rows.Close()

	// Collect all the permissions for the user
	for rows.Next() {
		var permission entity.Permission
		if err := rows.Scan(&permission.Resource, &permission.ResourceIdentifier, &permission.Operation, &permission.Effect); err != nil {
			return nil, fmt.Errorf("failed to scan permission: %v", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}
