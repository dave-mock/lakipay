package psql

import (
	"auth/src/pkg/erp/core/entity"
	"context"
	"time"

	"github.com/google/uuid"
)

// method for storing a catalog in PostgreSQL
func (repo PsqlRepo) CreateCatalog(userId uuid.UUID, name string, description string, status string) (*entity.Catalog, error) {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		repo.log.Println("CREATE CATALOG ERROR: Failed to start transaction")
		return nil, err
	}

	// Generate the catalog ID
	catalogId := uuid.New()
	err = repo.db.QueryRow(`INSERT INTO erp.catalogs (id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`, catalogId, userId, name, description, status, time.Now(), time.Now(), userId, userId).Scan(&catalogId)

	if err != nil {
		tx.Rollback()
		repo.log.Println("CREATE CATALOG ERROR: Failed to insert catalog into the database")
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		repo.log.Println("CREATE CATALOG ERROR: Failed to commit transaction")
		tx.Rollback()
		return nil, err
	}

	// Return the created catalog
	catalog := &entity.Catalog{
		Id:          catalogId,
		MerchantId:  userId,
		Name:        name,
		Description: description,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   userId,
		UpdatedBy:   userId,
	}

	repo.log.Println("CREATE CATALOG SUCCESS: Catalog successfully created")
	return catalog, nil
}

// method for retrieving catalogs of a specific merchant
func (repo PsqlRepo) ListMerchantCatalogs(merchantID string, userId uuid.UUID) ([]entity.Catalog, error) {
	var catalogs []entity.Catalog

	rows, err := repo.db.Query(`
		SELECT id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
		FROM erp.catalogs
		WHERE merchant_id = $1 AND created_by = $2
	`, merchantID, userId)
	if err != nil {
		repo.log.Println("LIST MERCHANT CATALOGS ERROR: Failed to query catalogs")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var catalog entity.Catalog
		if err := rows.Scan(&catalog.Id, &catalog.MerchantId, &catalog.Name, &catalog.Description, &catalog.Status, &catalog.CreatedAt, &catalog.UpdatedAt, &catalog.CreatedBy, &catalog.UpdatedBy); err != nil {
			repo.log.Println("LIST MERCHANT CATALOGS ERROR: Failed to scan catalog row")
			return nil, err
		}
		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		repo.log.Println("LIST MERCHANT CATALOGS ERROR: Row iteration error")
		return nil, err
	}

	repo.log.Println("LIST MERCHANT CATALOGS SUCCESS: Catalogs retrieved successfully")
	return catalogs, nil
}

// method for listing all catalogs of a user
func (repo PsqlRepo) ListCatalogs(userId uuid.UUID) ([]entity.Catalog, error) {
	var catalogs []entity.Catalog
	rows, err := repo.db.Query(`
		SELECT id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
		FROM erp.catalogs
		WHERE created_by = $1
	`, userId)
	if err != nil {
		repo.log.Println("LIST CATALOGS ERROR: Failed to query catalogs")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var catalog entity.Catalog
		if err := rows.Scan(&catalog.Id, &catalog.MerchantId, &catalog.Name, &catalog.Description, &catalog.Status, &catalog.CreatedAt, &catalog.UpdatedAt, &catalog.CreatedBy, &catalog.UpdatedBy); err != nil {
			repo.log.Println("LIST CATALOGS ERROR: Failed to scan catalog row")
			return nil, err
		}
		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		repo.log.Println("LIST CATALOGS ERROR: Row iteration error")
		return nil, err
	}

	repo.log.Println("LIST CATALOGS SUCCESS: Catalogs retrieved successfully")
	return catalogs, nil
}

// method for retrieving a catalog by ID
func (repo PsqlRepo) GetCatalog(catalogID string, userId uuid.UUID) (*entity.Catalog, error) {
	var catalog entity.Catalog

	err := repo.db.QueryRow(`
		SELECT id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
		FROM erp.catalogs
		WHERE id = $1 AND created_by = $2
	`, catalogID, userId).Scan(&catalog.Id, &catalog.MerchantId, &catalog.Name, &catalog.Description, &catalog.Status, &catalog.CreatedAt, &catalog.UpdatedAt, &catalog.CreatedBy, &catalog.UpdatedBy)

	if err != nil {
		repo.log.Println("GET CATALOG ERROR: Catalog not found")
		return nil, err
	}

	repo.log.Println("GET CATALOG SUCCESS: Catalog retrieved successfully")
	return &catalog, nil
}

// method for archiving a catalog
func (repo PsqlRepo) ArchiveCatalog(catalogID string, userId uuid.UUID) error {
	// Update the status of the catalog to 'archived'
	_, err := repo.db.Exec(`
		UPDATE erp.catalogs
		SET status = 'archived', updated_at = $1, updated_by = $2
		WHERE id = $3 AND created_by = $2
	`, time.Now(), userId, catalogID)
	if err != nil {
		repo.log.Println("ARCHIVE CATALOG ERROR: Failed to archive catalog")
		return err
	}

	repo.log.Println("ARCHIVE CATALOG SUCCESS: Catalog archived successfully")
	return nil
}

// method for updating catalog details
func (repo PsqlRepo) UpdateCatalog(userId uuid.UUID, catalogID, name, description, status string) (*entity.Catalog, error) {
	repo.log.Println("Updating Catalog")

	// Start a new transaction
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		repo.log.Println("UPDATE CATALOG ERROR: Failed to start transaction")
		return nil, err
	}

	// Update the catalog information and immediately retrieve the updated data
	var updatedCatalog entity.Catalog
	err = repo.db.QueryRow(`
		UPDATE erp.catalogs
		SET name = $1, description = $2, status = $3, updated_at = $4, updated_by = $5::UUID
		WHERE id = $6::UUID AND created_by = $5::UUID
		RETURNING id, name, description, status, created_at, updated_at, created_by, updated_by`,
		name, description, status, time.Now(), userId, catalogID).
		Scan(&updatedCatalog.Id, &updatedCatalog.Name, &updatedCatalog.Description, &updatedCatalog.Status, &updatedCatalog.CreatedAt, &updatedCatalog.UpdatedAt, &updatedCatalog.CreatedBy, &updatedCatalog.UpdatedBy)

	if err != nil {
		repo.log.Println("UPDATE CATALOG ERROR: Failed to update catalog or retrieve updated data")
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		repo.log.Println("UPDATE CATALOG ERROR: Failed to commit transaction")
		tx.Rollback()
		return nil, err
	}

	repo.log.Println("UPDATE CATALOG SUCCESS: Catalog updated successfully")
	return &updatedCatalog, nil
}

// method for deleting a catalog
func (repo PsqlRepo) DeleteCatalog(catalogID string, userId uuid.UUID) error {
	_, err := repo.db.Exec(`
		DELETE FROM erp.catalogs
		WHERE id = $1 AND created_by = $2
	`, catalogID, userId)
	if err != nil {
		repo.log.Println("DELETE CATALOG ERROR: Failed to delete catalog")
		return err
	}

	repo.log.Println("DELETE CATALOG SUCCESS: Catalog deleted successfully")
	return nil
}
