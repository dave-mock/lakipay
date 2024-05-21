package repo

import (
	"auth/src/pkg/auth/core/entity"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
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
	`, repo.schema), phoneAuth.Id, phoneAuth.Token, phoneAuth.Phone.Id, phoneAuth.Code, phoneAuth.Method, phoneAuth.Status, phoneAuth.CreatedAt)

	return err
}

func (repo PsqlRepo) FindPhoneAuth(token string) (entity.PhoneAuth, error) {

	var phoneAuth entity.PhoneAuth

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

	err := repo.db.QueryRow(fmt.Sprintf(`
	SELECT sessions.id, sessions.token, users.id, users.sir_name, users.first_name, users.last_name
	FROM %s.sessions
	INNER JOIN %s.users ON %s.users.id = sessions.user_id
	WHERE sessions.id = $1::UUID
	`, repo.schema, repo.schema, repo.schema), id).Scan(&session.Id, &session.Token,
		&session.User.Id, &sirName, &session.User.FirstName, &lastName,
	)

	if sirName.Valid {
		session.User.SirName = sirName.String
	}

	if lastName.Valid {
		session.User.LastName = lastName.String
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
