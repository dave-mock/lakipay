package psql

import (
	"auth/src/pkg/erp/core/entity"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

var ErrProductDoesNotExist = errors.New("product does not exist")

func (repo PsqlRepo) ListProducts(userID uuid.UUID, userType string) ([]entity.Product, error) {
	var products []entity.Product
	var query string

	var rows *sql.Rows
	var err error

	if userType == "admin" {
		query = `
            SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
            FROM erp.products
        `
		rows, err = repo.db.Query(query)
		if err != nil {
			repo.log.Println("ERROR LISTING PRODUCTS:", err)
			return nil, err
		}
		defer rows.Close()

	} else {

		query = `
            SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
            FROM erp.products
            WHERE merchant_id = $1 
        `
		rows, err = repo.db.Query(query, userID)
		if err != nil {
			repo.log.Println("ERROR LISTING PRODUCTS:", err)
			return nil, err
		}
		defer rows.Close() // Ensure rows are closed after use
	}

	// Iterate through the rows and scan data into the products slice
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.MerchantID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
			&product.SKU,
			&product.Weight,
			&product.Dimensions,
			&product.ImageURL,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CreatedBy,
			&product.UpdatedBy,
		)
		if err != nil {
			repo.log.Println("ERROR SCANNING PRODUCT:", err)
			return nil, err
		}
		products = append(products, product)
	}

	// Check for any row iteration errors
	if err := rows.Err(); err != nil {
		repo.log.Println("ERROR IN ROWS:", err)
		return nil, err
	}

	return products, nil
}

func (repo PsqlRepo) GetProduct(productID, userID uuid.UUID) (*entity.Product, error) {
	var product entity.Product

	err := repo.db.QueryRow(`
		SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
		FROM erp.products
		WHERE id = $1 AND merchant_id = $2
	`, productID, userID).Scan(
		&product.Id,
		&product.MerchantID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Currency,
		&product.SKU,
		&product.Weight,
		&product.Dimensions,
		&product.ImageURL,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
		&product.UpdatedBy)

	if err != nil {
		repo.log.Println("GET PRODUCT ERROR:", err)
		return nil, err
	}

	return &product, nil
}

func (repo PsqlRepo) CreateProduct(name, description string, price float64, currency, sku string, weight float64, dimensions, imageURL string, createdBy uuid.UUID) error {
	const (
		ErrFailedToStoreProduct = "FAILED_TO_STORE_PRODUCT"
		ErrDuplicateSKU         = "DUPLICATE_SKU"
	)

	var existingProductID uuid.UUID
	err := repo.db.QueryRow(`
		SELECT id FROM erp.products WHERE sku = $1
	`, sku).Scan(&existingProductID)

	if err == nil {
		repo.log.Printf("Duplicate SKU detected for SKU: %s", sku)
		return fmt.Errorf("%s: A product with the SKU '%s' already exists.", ErrDuplicateSKU, sku)
	} else if err != sql.ErrNoRows {
		repo.log.Println("Error checking for duplicate SKU:", err)
		return fmt.Errorf("%s: %v", ErrFailedToStoreProduct, err)
	}

	// Proceed to insert the product if no duplicate
	productID := uuid.New()
	createdAt := time.Now()
	updatedAt := createdAt

	_, err = repo.db.Exec(`
		INSERT INTO erp.products (
			id, name, description, price, currency, sku, weight, dimensions, image_url,
			created_at, updated_at, created_by, updated_by, merchant_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $12
		)
	`, productID, name, description, price, currency, sku, weight, dimensions, imageURL, createdAt, updatedAt, createdBy, createdBy)

	if err != nil {
		repo.log.Println("CREATE PRODUCT ERROR:", err)
		return fmt.Errorf("%s: %v", ErrFailedToStoreProduct, err)
	}

	repo.log.Println("PRODUCT CREATED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) UpdateProduct(productID uuid.UUID, name string, description string, price float64, currency string,
	SKU string,
	weight float64,
	dimensions string,
	imageURL string,
	status string,
	userId uuid.UUID) error {

	updatedAt := time.Now()

	repo.log.Println("Updating Product with values:")
	repo.log.Printf("ID: %s", productID)
	repo.log.Printf("Name: %s", name)
	repo.log.Printf("Description: %s", description)
	repo.log.Printf("Price: %.2f", price)
	repo.log.Printf("Currency: %s", currency)
	repo.log.Printf("SKU: %s", SKU)
	repo.log.Printf("Weight: %.2f", weight)
	repo.log.Printf("Dimensions: %s", dimensions)
	repo.log.Printf("Image URL: %s", imageURL)
	repo.log.Printf("Status: %s", status)
	repo.log.Printf("Updated At: %s", updatedAt)
	repo.log.Printf("Merchant ID: %s", userId)

	_, err := repo.db.Exec(`
	UPDATE erp.products
	SET 
		name = $1, 
		description = $2, 
		price = $3, 
		currency = $4, 
		sku = $5, 
		status = $6, 
		weight = $7, 
		dimensions = $8, 
		image_url = $9, 
		updated_at = $10, 
		updated_by = $11
	WHERE id = $12 AND merchant_id = $13
	`, name, description, price, currency, SKU, status, weight, dimensions, imageURL, updatedAt, userId, productID, userId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			repo.log.Println("SKU already exists, cannot update product")
			return fmt.Errorf("duplicate SKU: product with this SKU already exists")
		}
		repo.log.Println("UPDATE PRODUCT ERROR:", err)
		return fmt.Errorf("failed to update product: %v", err)
	}

	repo.log.Println("PRODUCT UPDATED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) DeactivateProduct(productID uuid.UUID) error {
	_, err := repo.db.Exec(`
		UPDATE erp.products
		SET status = 'inactive', updated_at = $1
		WHERE id = $2
	`, time.Now(), productID)
	if err != nil {
		repo.log.Println("DEACTIVATE PRODUCT ERROR:", err)
		return err
	}
	return nil
}

func (repo PsqlRepo) DeleteProduct(productID uuid.UUID, userID uuid.UUID) error {
	const ErrFailedToDeleteProduct = "FAILED_TO_DELETE_PRODUCT"

	_, err := repo.db.Exec(`
		DELETE FROM erp.products
		WHERE id = $1
	`, productID)

	if err != nil {
		repo.log.Println("ERROR DELETING PRODUCT:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeleteProduct, err)
	}

	repo.log.Println("PRODUCT DELETED SUCCESSFULLY", "productID", productID, "userID", userID)
	return nil
}

func (repo PsqlRepo) ListAllProducts() ([]entity.Product, error) {
	var products []entity.Product
	rows, err := repo.db.Query(`
		SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
		FROM erp.products
	`)
	if err != nil {
		repo.log.Println("ERROR LISTING PRODUCTS:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.MerchantID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
			&product.SKU,
			&product.Weight,
			&product.Dimensions,
			&product.ImageURL,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CreatedBy,
			&product.UpdatedBy,
		)
		if err != nil {
			repo.log.Println("ERROR SCANNING PRODUCT:", err)
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		repo.log.Println("ERROR IN ROWS:", err)
		return nil, err
	}

	return products, nil
}

func (repo PsqlRepo) GetProductBySKU(sku string) (*entity.Product, error) {
	var product entity.Product

	err := repo.db.QueryRow(`
		SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
		FROM erp.products
		WHERE sku = $1
	`, sku).Scan(
		&product.Id,
		&product.MerchantID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Currency,
		&product.SKU,
		&product.Weight,
		&product.Dimensions,
		&product.ImageURL,
		&product.Status,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CreatedBy,
		&product.UpdatedBy)
	if err != nil {
		repo.log.Println("GET PRODUCT BY SKU ERROR:", err)
		return nil, err
	}

	return &product, nil
}

func (repo PsqlRepo) ListProductsByStatus(merchantID uuid.UUID, status string) ([]entity.Product, error) {
	var products []entity.Product

	rows, err := repo.db.Query(`
		SELECT id, merchant_id, name, description, price, currency, sku, weight, dimensions, image_url, status, created_at, updated_at, created_by, updated_by
		FROM erp.products
		WHERE merchant_id = $1 AND status = $2
	`, merchantID, status)
	if err != nil {
		repo.log.Println("LIST PRODUCTS BY STATUS ERROR:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		if err := rows.Scan(
			&product.Id,
			&product.MerchantID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
			&product.SKU,
			&product.Weight,
			&product.Dimensions,
			&product.ImageURL,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CreatedBy,
			&product.UpdatedBy); err != nil {
			repo.log.Println("LIST PRODUCTS BY STATUS SCAN ERROR:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (repo PsqlRepo) CountProductsByMerchant(merchantID uuid.UUID) (int, error) {
	var count int

	err := repo.db.QueryRow(`
		SELECT COUNT(*) FROM erp.products
		WHERE merchant_id = $1
	`, merchantID).Scan(&count)
	if err != nil {
		repo.log.Println("COUNT PRODUCTS ERROR:", err)
		return 0, err
	}

	return count, nil
}

func (repo PsqlRepo) AddProductToCatalog(productId, catalogId uuid.UUID, displayOrder int, userId uuid.UUID) ([]entity.ProductCatalog, error) {
	var exists bool
	err := repo.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM erp.products WHERE id = $1);
	`, productId).Scan(&exists)
	if err != nil {
		repo.log.Println("ERROR CHECKING IF PRODUCT EXISTS:", err)
		return nil, fmt.Errorf("failed to check if product exists: %w", err)
	}
	if !exists {
		repo.log.Println("ERROR: Product with the given ID does not exist")
		return nil, ErrProductDoesNotExist // Return the specific error
	}

	var catalogs []entity.ProductCatalog
	rows, err := repo.db.Query(`
		INSERT INTO erp.product_catalog (product_id, catalog_id, display_order, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, product_id, catalog_id, display_order, created_by, updated_by, created_at, updated_at
	`, productId, catalogId, displayOrder, userId)

	if err != nil {
		repo.log.Println("ERROR ADDING PRODUCT TO CATALOG:", err)
		return nil, fmt.Errorf("failed to add product to catalog: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var catalog entity.ProductCatalog
		if err := rows.Scan(&catalog.Id, &catalog.ProductID, &catalog.CatalogID, &catalog.DisplayOrder, &catalog.CreatedBy, &catalog.UpdatedBy, &catalog.CreatedAt, &catalog.UpdatedAt); err != nil {
			repo.log.Println("ERROR SCANNING PRODUCT CATALOG ROW:", err)
			return nil, fmt.Errorf("failed to scan product catalog row: %w", err)
		}
		catalogs = append(catalogs, catalog)
	}

	if err := rows.Err(); err != nil {
		repo.log.Println("ERROR ITERATING OVER PRODUCT CATALOG ROWS:", err)
		return nil, fmt.Errorf("failed to iterate over product catalog rows: %w", err)
	}

	repo.log.Println("PRODUCT ADDED TO CATALOG SUCCESSFULLY")
	return catalogs, nil
}

func (repo PsqlRepo) ArchiveProduct(productId uuid.UUID) error {
	_, err := repo.db.Exec(`
		UPDATE products SET status = 'archived', updated_at = NOW() WHERE id = $1
	`, productId)

	if err != nil {
		repo.log.Println("ERROR ARCHIVING PRODUCT:", err)
		return fmt.Errorf("failed to archive product: %w", err)
	}

	repo.log.Println("PRODUCT ARCHIVED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) RemoveProductFromCatalog(productId, catalogId uuid.UUID, userId uuid.UUID) ([]entity.ProductCatalog, error) {
	var productCatalogs []entity.ProductCatalog
	rows, err := repo.db.Query(`
		DELETE FROM erp.product_catalog WHERE product_id = $1 AND catalog_id = $2
		RETURNING id, product_id, catalog_id, display_order, created_at, updated_at, created_by, updated_by
	`, productId, catalogId)

	if err != nil {
		repo.log.Println("ERROR REMOVING PRODUCT FROM CATALOG:", err)
		return nil, fmt.Errorf("failed to remove product from catalog: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var catalog entity.ProductCatalog
		if err := rows.Scan(
			&catalog.Id,
			&catalog.ProductID,
			&catalog.CatalogID,
			&catalog.DisplayOrder,
			&catalog.CreatedAt,
			&catalog.UpdatedAt,
			&catalog.CreatedBy,
			&catalog.UpdatedBy); err != nil {
			repo.log.Println("ERROR SCANNING PRODUCT CATALOG ROW:", err)
			return nil, fmt.Errorf("failed to scan product catalog row: %w", err)
		}
		productCatalogs = append(productCatalogs, catalog)
	}
	if err := rows.Err(); err != nil {
		repo.log.Println("ERROR ITERATING OVER PRODUCT CATALOG ROWS:", err)
		return nil, fmt.Errorf("failed to iterate over product catalog rows: %w", err)
	}

	repo.log.Println("PRODUCT REMOVED FROM CATALOG SUCCESSFULLY")
	return productCatalogs, nil
}
