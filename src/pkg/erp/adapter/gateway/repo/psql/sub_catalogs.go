package psql

import (
	"auth/src/pkg/erp/core/entity"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateSubCatalog(subCatalog *entity.SubCatalog) error {
	// Check if the provided catalog_id exists in the catalogs table
	var exists bool
	err := repo.db.QueryRow(`
        SELECT EXISTS (SELECT 1 FROM erp.catalogs WHERE id = $1)
    `, subCatalog.CatalogID).Scan(&exists)

	if err != nil {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Failed to check if catalog exists")
		return err
	}

	if !exists {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Catalog ID does not exist")
		return fmt.Errorf("INVALID_CATALOG_ID: The provided catalog_id does not exist in the catalogs table")
	}

	// Check if a sub-catalog already exists for the given catalog_id and merchant_id
	var subCatalogExists bool
	err = repo.db.QueryRow(`
        SELECT EXISTS (SELECT 1 FROM erp.sub_catalogs WHERE catalog_id = $1 AND merchant_id = $2 AND sub_catalog_id = $3)
    `, subCatalog.CatalogID, subCatalog.MerchantID, subCatalog.SubCatalogID).Scan(&subCatalogExists)

	if err != nil {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Failed to check if sub-catalog already exists")
		return err
	}

	if subCatalogExists {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Sub-catalog already exists for the given catalog and merchant")
		return fmt.Errorf("SUB_CATALOG_ALREADY_EXISTS: A sub-catalog with the given catalog_id and merchant_id already exists")
	}

	// Start the transaction
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Failed to start transaction")
		return err
	}

	if subCatalog.SubCatalogID == uuid.Nil {
		subCatalog.SubCatalogID = uuid.New()
	}

	_, err = repo.db.Exec(`
        INSERT INTO erp.sub_catalogs (id, catalog_id, sub_catalog_id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		subCatalog.ID, subCatalog.CatalogID, subCatalog.SubCatalogID,
		subCatalog.MerchantID, subCatalog.Name, subCatalog.Description, subCatalog.Status,
		subCatalog.CreatedAt, subCatalog.UpdatedAt,
		subCatalog.CreatedBy, subCatalog.UpdatedBy,
	)

	if err != nil {
		tx.Rollback()
		repo.log.Println("CREATE SUB-CATALOG ERROR: Failed to insert sub-catalog into the database")
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		repo.log.Println("CREATE SUB-CATALOG ERROR: Failed to commit transaction")
		tx.Rollback()
		return err
	}

	repo.log.Println("CREATE SUB-CATALOG SUCCESS: Sub-catalog successfully created")
	return nil
}

func (repo PsqlRepo) ListSubCatalogs(userId uuid.UUID) ([]entity.SubCatalog, error) {
	var subCatalogs []entity.SubCatalog
	rows, err := repo.db.Query(`
		SELECT id, catalog_id, sub_catalog_id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
		FROM erp.sub_catalogs
		WHERE created_by = $1
	`, userId)
	if err != nil {
		repo.log.Println("LIST SUB-CATALOGS ERROR: Failed to query sub-catalogs")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var subCatalog entity.SubCatalog
		if err := rows.Scan(&subCatalog.ID,
			&subCatalog.CatalogID,
			&subCatalog.SubCatalogID,
			&subCatalog.MerchantID,
			&subCatalog.Name,
			&subCatalog.Description,
			&subCatalog.Status,
			&subCatalog.CreatedAt,
			&subCatalog.UpdatedAt,
			&subCatalog.CreatedBy,
			&subCatalog.UpdatedBy); err != nil {
			repo.log.Println("LIST SUB-CATALOGS ERROR: Failed to scan sub-catalog row")
			return nil, err
		}
		subCatalogs = append(subCatalogs, subCatalog)
	}

	if err := rows.Err(); err != nil {
		repo.log.Println("LIST SUB-CATALOGS ERROR: Row iteration error")
		return nil, err
	}

	repo.log.Println("LIST SUB-CATALOGS SUCCESS: Sub-catalogs retrieved successfully")
	return subCatalogs, nil
}

// method for retrieving a sub-catalog by ID
func (repo PsqlRepo) GetSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) (*entity.SubCatalog, error) {
	var subCatalog entity.SubCatalog

	err := repo.db.QueryRow(`
		SELECT id, catalog_id, sub_catalog_id, merchant_id, name, description, status, created_at, updated_at, created_by, updated_by
		FROM erp.sub_catalogs
		WHERE id = $1 AND created_by = $2
	`, subCatalogID, userId).Scan(&subCatalog.ID,
		&subCatalog.CatalogID,
		&subCatalog.SubCatalogID,
		&subCatalog.MerchantID,
		&subCatalog.Name,
		&subCatalog.Description,
		&subCatalog.Status,
		&subCatalog.CreatedAt,
		&subCatalog.UpdatedAt,
		&subCatalog.CreatedBy,
		&subCatalog.UpdatedBy)

	if err != nil {
		repo.log.Println("GET SUB-CATALOG ERROR: Sub-catalog not found")
		return nil, err
	}

	repo.log.Println("GET SUB-CATALOG SUCCESS: Sub-catalog retrieved successfully")
	return &subCatalog, nil
}

func (repo PsqlRepo) UpdateSubCatalog(subCatalog *entity.SubCatalog) error {
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		repo.log.Println("UPDATE SUB-CATALOG ERROR: Failed to start transaction")
		return err
	}

	_, err = repo.db.Exec(`
        UPDATE erp.sub_catalogs 
        SET name = $1, description = $2, status = $3, updated_at = $4, updated_by = $5
        WHERE id = $6`,
		subCatalog.Name, subCatalog.Description, subCatalog.Status,
		time.Now(), subCatalog.UpdatedBy, subCatalog.ID,
	)

	if err != nil {
		tx.Rollback()
		repo.log.Println("UPDATE SUB-CATALOG ERROR: Failed to update sub-catalog")
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.log.Println("UPDATE SUB-CATALOG ERROR: Failed to commit transaction")
		return err
	}

	repo.log.Println("UPDATE SUB-CATALOG SUCCESS: Sub-catalog successfully updated")
	return nil
}

func (repo PsqlRepo) DeleteSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) error {
	// Check if the sub-catalog exists
	var exists bool
	err := repo.db.QueryRow(`
        SELECT EXISTS (SELECT 1 FROM erp.sub_catalogs WHERE id = $1 AND created_by = $2)
    `, subCatalogID, userId).Scan(&exists)

	if err != nil {
		repo.log.Println("DELETE SUB-CATALOG ERROR: Failed to check if sub-catalog exists")
		return err
	}

	if !exists {
		repo.log.Println("DELETE SUB-CATALOG ERROR: Sub-catalog does not exist or you do not have permission")
		return fmt.Errorf("SUB_CATALOG_NOT_FOUND: The sub-catalog does not exist or you do not have permission to delete it")
	}

	// Proceed to delete the sub-catalog
	_, err = repo.db.Exec(`
		DELETE FROM erp.sub_catalogs
		WHERE id = $1 AND created_by = $2
	`, subCatalogID, userId)

	if err != nil {
		repo.log.Println("DELETE SUB-CATALOG ERROR: Failed to delete sub-catalog")
		return err
	}

	repo.log.Println("DELETE SUB-CATALOG SUCCESS: Sub-catalog deleted successfully")
	return nil
}
func (repo PsqlRepo) ArchiveSubCatalog(subCatalogID uuid.UUID, userId uuid.UUID) error {
	// Check if the sub-catalog exists
	var exists bool
	err := repo.db.QueryRow(`
        SELECT EXISTS (SELECT 1 FROM erp.sub_catalogs WHERE id = $1 AND created_by = $2)
    `, subCatalogID, userId).Scan(&exists)

	if err != nil {
		repo.log.Println("ARCHIVE SUB-CATALOG ERROR: Failed to check if sub-catalog exists")
		return err
	}

	if !exists {
		repo.log.Println("ARCHIVE SUB-CATALOG ERROR: Sub-catalog does not exist or you do not have permission")
		return fmt.Errorf("SUB_CATALOG_NOT_FOUND: The sub-catalog does not exist or you do not have permission to archive it")
	}

	// Proceed to archive the sub-catalog
	_, err = repo.db.Exec(`
		UPDATE erp.sub_catalogs
		SET status = 'archived', updated_at = $1, updated_by = $2
		WHERE id = $3 AND created_by = $2
	`, time.Now(), userId, subCatalogID)

	if err != nil {
		repo.log.Println("ARCHIVE SUB-CATALOG ERROR: Failed to archive sub-catalog")
		return err
	}

	repo.log.Println("ARCHIVE SUB-CATALOG SUCCESS: Sub-catalog archived successfully")
	return nil
}
