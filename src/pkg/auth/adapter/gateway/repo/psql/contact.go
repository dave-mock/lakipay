package repo

import (
	"auth/src/pkg/auth/core/entity"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Save Contact
func (repo PsqlRepo) StoreContact(userId uuid.UUID, v entity.Contact) error {

	_, err := repo.db.Exec(fmt.Sprintf(`
	INSERT INTO %s.contacts (id, sir_name, first_name, last_name, contact_of)
	VALUES ($1::UUID, $2, $3, $4, $5::UUID)
	`, repo.schema),
		v.Id,
		sql.NullString{Valid: v.SirName != "", String: v.SirName},
		sql.NullString{Valid: v.FirstName != "", String: v.FirstName},
		sql.NullString{Valid: v.LastName != "", String: v.LastName},
		userId,
	)

	if err != nil {
		return err
	}

	return nil
}

// // Find contacts of a specified user
// func (repo PsqlRepo) FindUserContacts(userId uuid.UUID) ([]entity.Contact, error) {
// 	var contacts []entity.Contact = make([]entity.Contact, 0)

// 	rows, err := repo.db.Query(fmt.Sprintf(`
// 	SELECT id, sir_name, first_name, last_name,
// 	users.sir_name, users.first_name, users.last_name, users.gender, users.date_of_birth, users.created_at
// 	FROM %s.contacts
// 	INNER JOIN %s.users ON %s.users.id = id
// 	WHERE contact_of = $1;
// 	`), userId)

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var contact entity.Contact
// 		err = rows.Scan()
// 	}

// 	return contacts, nil
// }
