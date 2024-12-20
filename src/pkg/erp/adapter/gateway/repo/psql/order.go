package psql

import (
	"auth/src/pkg/erp/core/entity"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateOrder(order entity.Order) (*entity.Order, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		repo.log.Println("BEGIN TRANSACTION ERROR:", err)
		return nil, err
	}
	defer tx.Rollback()

	orderID := uuid.New()
	createdAt := time.Now()
	createdAtStr := createdAt.Format(time.RFC3339)
	updatedAtStr := createdAtStr

	repo.log.Printf("Inserting order into erp.orders with orderID: %v, customerID: %v, orderTypeID: %v", orderID, order.CustomerDetails.CustomerID, order.OrderDetails.OrderTypeID)

	_, err = tx.Exec(`
        INSERT INTO erp.orders (
            id, customer_id, order_type_id, total_amount, currency, medium, 
            status, payment_status, payment_method, payment_reference, 
            shipping_address, billing_address, final_amount, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
        )
    `, orderID,
		order.CustomerDetails.CustomerID,
		order.OrderDetails.OrderTypeID,
		order.OrderDetails.TotalAmount,
		order.OrderDetails.Currency,
		order.OrderDetails.Medium,
		order.OrderDetails.Status,
		order.OrderDetails.PaymentStatus,
		order.OrderDetails.PaymentMethod,
		order.OrderDetails.PaymentRef,
		order.OrderDetails.ShippingAddr,
		order.OrderDetails.BillingAddr,
		order.OrderDetails.FinalAmount,
		createdAtStr,
		updatedAtStr)
	if err != nil {
		repo.log.Println("INSERT ORDER ERROR:", err)
		return nil, err
	}

	for _, item := range order.OrderItems {
		var productExists bool
		err := tx.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM erp.products WHERE id = $1)
		`, item.ProductID).Scan(&productExists)
		if err != nil {
			repo.log.Println("ERROR CHECKING PRODUCT EXISTENCE:", err)
			return nil, err
		}

		if !productExists {
			repo.log.Printf("Product with ID %v does not exist", item.ProductID)
			return nil, fmt.Errorf("Product with ID %v does not exist", item.ProductID)
		}

		repo.log.Printf("Inserting order item with productID: %v, productName: %v, quantity: %v", item.ProductID, item.ProductName, item.Quantity)
		_, err = tx.Exec(`
            INSERT INTO erp.order_items (id, order_id, product_id, product_name, quantity, unit_price, total_price, category, sku, merchant_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, uuid.New(), orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice, item.Category, item.SKU, item.MerchantID)

		if err != nil {
			repo.log.Println("INSERT ORDER ITEM ERROR:", err)
			return nil, err
		}
	}

	for _, discount := range order.OrderDetails.Discounts {
		repo.log.Printf("Inserting discount with type: %v, value: %v", discount.Type, discount.Value)
		_, err := tx.Exec(`
            INSERT INTO erp.discounts (id, order_id, type, value, description)
            VALUES ($1, $2, $3, $4, $5)
        `, uuid.New(), orderID, discount.Type, discount.Value, discount.Description)

		if err != nil {
			repo.log.Println("||||||   INSERTING DISCOUNT ERROR |||||:", err)
			return nil, err
		}
	}

	// Insert taxes
	for _, tax := range order.OrderDetails.Taxes {
		repo.log.Printf("Inserting tax with type: %v, rate: %v", tax.Type, tax.Rate)
		_, err := tx.Exec(`
            INSERT INTO erp.taxes (id, order_id, type, rate, value)
            VALUES ($1, $2, $3, $4, $5)
        `, uuid.New(), orderID, tax.Type, tax.Rate, tax.Value)

		if err != nil {
			repo.log.Println("INSERT TAX ERROR:", err)
			return nil, err
		}
	}

	if order.Tracking.ExpectedDeliveryDate == "" {
		order.Tracking.ExpectedDeliveryDate = time.Now().Format(time.RFC3339) // Set current time if empty
	}

	repo.log.Printf("Inserting tracking with status: %v, expectedDeliveryDate: %v", order.Tracking.Status, order.Tracking.ExpectedDeliveryDate)

	_, err = tx.Exec(`
		INSERT INTO erp.tracking (id, order_id, status, expected_delivery_date, shipment_id)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New(), orderID, order.Tracking.Status, order.Tracking.ExpectedDeliveryDate, order.Tracking.ShipmentID)

	if err != nil {
		repo.log.Println("INSERT TRACKING ERROR:", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		repo.log.Println("COMMIT TRANSACTION ERROR:", err)
		return nil, err
	}
	order.ID = orderID
	return &order, nil
}

func (repo PsqlRepo) UpdateOrder(order entity.Order) (entity.Order, error) {
	repo.log.Println("Updating order", "orderID", order.ID)
	_, err := repo.db.Exec(`
        UPDATE erp.orders
        SET 
            merchant_id = $1,
            customer_id = $2,
            order_type_id = $3,
            total_amount = $4,
            currency = $5,
            medium = $6,
            shipping_address = $7,
            billing_address = $8,
            updated_at = $9
        WHERE id = $10
    `,
		order.CustomerDetails.CustomerID,
		order.OrderDetails.OrderTypeID,
		order.OrderDetails.TotalAmount,
		order.OrderDetails.Currency,
		order.OrderDetails.Medium,
		order.OrderDetails.ShippingAddr,
		order.OrderDetails.BillingAddr,
		time.Now(),
		order.ID,
	)
	if err != nil {
		repo.log.Println("UPDATE ORDER ERROR:", err)
		return entity.Order{}, fmt.Errorf("failed to update order: %w", err)
	}

	var updatedOrder entity.Order
	err = repo.db.QueryRow(`
        SELECT id, 
            customer_id, 
            order_type_id, 
            total_amount, 
            currency, 
            medium, 
            shipping_address, 
            billing_address, 
            status, 
            created_at, 
            updated_at 
        FROM erp.orders 
        WHERE id = $1
    `, order.ID).Scan(
		&updatedOrder.ID,
		&updatedOrder.CustomerDetails.CustomerID,
		&updatedOrder.OrderDetails.OrderTypeID,
		&updatedOrder.OrderDetails.TotalAmount,
		&updatedOrder.OrderDetails.Currency,
		&updatedOrder.OrderDetails.Medium,
		&updatedOrder.OrderDetails.ShippingAddr,
		&updatedOrder.OrderDetails.BillingAddr,
		&updatedOrder.OrderDetails.Status,
		&updatedOrder.Metadata.CreatedAt,
		&updatedOrder.Metadata.UpdatedAt,
	)

	if err != nil {
		repo.log.Println("ERROR FETCHING UPDATED ORDER:", err)
		return entity.Order{}, fmt.Errorf("failed to fetch updated order: %w", err)
	}

	repo.log.Println("Order updated successfully", "orderID", updatedOrder.ID)
	return updatedOrder, nil
}

func (repo PsqlRepo) CreateCartOrder(order entity.Order) (*entity.Order, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		repo.log.Println("BEGIN TRANSACTION ERROR:", err)
		return nil, err
	}
	defer tx.Rollback()

	orderID := uuid.New()
	createdAt := time.Now()
	createdAtStr := createdAt.Format(time.RFC3339)
	updatedAtStr := createdAtStr

	repo.log.Printf("Inserting order into erp.cart_orders with orderID: %v, customerID: %v, orderTypeID: %v", orderID, order.CustomerDetails.CustomerID, order.OrderDetails.OrderTypeID)

	_, err = tx.Exec(`
        INSERT INTO erp.cart_orders (
            id, customer_id, order_type_id, total_amount, currency, medium, 
            status, payment_status, payment_method, payment_reference, 
            shipping_address, billing_address, final_amount, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
        )
    `, orderID,
		order.CustomerDetails.CustomerID,
		order.OrderDetails.OrderTypeID,
		order.OrderDetails.TotalAmount,
		order.OrderDetails.Currency,
		order.OrderDetails.Medium,
		order.OrderDetails.Status,
		order.OrderDetails.PaymentStatus,
		order.OrderDetails.PaymentMethod,
		order.OrderDetails.PaymentRef,
		order.OrderDetails.ShippingAddr,
		order.OrderDetails.BillingAddr,
		order.OrderDetails.FinalAmount,
		createdAtStr,
		updatedAtStr)
	if err != nil {
		repo.log.Println("INSERT ORDER ERROR:", err)
		return nil, err
	}

	for _, item := range order.OrderItems {
		var productExists bool
		err := tx.QueryRow(`
			SELECT EXISTS(SELECT 1 FROM erp.products WHERE id = $1)
		`, item.ProductID).Scan(&productExists)
		if err != nil {
			repo.log.Println("ERROR CHECKING PRODUCT EXISTENCE:", err)
			return nil, err
		}

		if !productExists {
			repo.log.Printf("Product with ID %v does not exist", item.ProductID)
			return nil, fmt.Errorf("Product with ID %v does not exist", item.ProductID)
		}

		repo.log.Printf("Inserting order item with productID: %v, productName: %v, quantity: %v", item.ProductID, item.ProductName, item.Quantity)
		_, err = tx.Exec(`
            INSERT INTO erp.cart_order_items (id, order_id, product_id, product_name, quantity, unit_price, total_price, category, sku, merchant_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        `, uuid.New(), orderID, item.ProductID, item.ProductName, item.Quantity, item.UnitPrice, item.TotalPrice, item.Category, item.SKU, item.MerchantID)

		if err != nil {
			repo.log.Println("INSERT ITEM INTO CART ERROR:", err)
			return nil, err
		}
	}

	for _, discount := range order.OrderDetails.Discounts {
		repo.log.Printf("Inserting discount with type: %v, value: %v", discount.Type, discount.Value)
		_, err := tx.Exec(`
            INSERT INTO erp.cart_discounts (id, order_id, type, value, description)
            VALUES ($1, $2, $3, $4, $5)
        `, uuid.New(), orderID, discount.Type, discount.Value, discount.Description)

		if err != nil {
			repo.log.Println("||||||   INSERTING DISCOUNT ERROR |||||:", err)
			return nil, err
		}
	}

	// Insert taxes
	for _, tax := range order.OrderDetails.Taxes {
		repo.log.Printf("Inserting tax with type: %v, rate: %v", tax.Type, tax.Rate)
		_, err := tx.Exec(`
            INSERT INTO erp.cart_taxes (id, order_id, type, rate, value)
            VALUES ($1, $2, $3, $4, $5)
        `, uuid.New(), orderID, tax.Type, tax.Rate, tax.Value)

		if err != nil {
			repo.log.Println("INSERT TAX ERROR:", err)
			return nil, err
		}
	}

	if order.Tracking.ExpectedDeliveryDate == "" {
		order.Tracking.ExpectedDeliveryDate = time.Now().Format(time.RFC3339)
	}

	repo.log.Printf("Inserting tracking with status: %v, expectedDeliveryDate: %v", order.Tracking.Status, order.Tracking.ExpectedDeliveryDate)

	_, err = tx.Exec(`
		INSERT INTO erp.cart_tracking (id, order_id, status, expected_delivery_date, shipment_id)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New(), orderID, order.Tracking.Status, order.Tracking.ExpectedDeliveryDate, order.Tracking.ShipmentID)

	if err != nil {
		repo.log.Println("INSERT TRACKING ERROR:", err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		repo.log.Println("COMMIT TRANSACTION ERROR:", err)
		return nil, err
	}
	order.ID = orderID
	return &order, nil
}

func (repo PsqlRepo) UpdateCartOrder(order entity.Order) (entity.Order, error) {
	repo.log.Println("Updating order", "orderID", order.ID)
	_, err := repo.db.Exec(`
        UPDATE erp.cart_orders
        SET 
            merchant_id = $1,
            customer_id = $2,
            order_type_id = $3,
            total_amount = $4,
            currency = $5,
            medium = $6,
            shipping_address = $7,
            billing_address = $8,
            updated_at = $9
        WHERE id = $10
    `,
		order.CustomerDetails.CustomerID,
		order.OrderDetails.OrderTypeID,
		order.OrderDetails.TotalAmount,
		order.OrderDetails.Currency,
		order.OrderDetails.Medium,
		order.OrderDetails.ShippingAddr,
		order.OrderDetails.BillingAddr,
		time.Now(),
		order.ID,
	)
	if err != nil {
		repo.log.Println("UPDATE CART ERROR:", err)
		return entity.Order{}, fmt.Errorf("failed to update cart: %w", err)
	}

	var updatedOrder entity.Order
	err = repo.db.QueryRow(`
        SELECT id, 
            customer_id, 
            order_type_id, 
            total_amount, 
            currency, 
            medium, 
            shipping_address, 
            billing_address, 
            status, 
            created_at, 
            updated_at 
        FROM erp.cart_orders 
        WHERE id = $1
    `, order.ID).Scan(
		&updatedOrder.ID,
		&updatedOrder.CustomerDetails.CustomerID,
		&updatedOrder.OrderDetails.OrderTypeID,
		&updatedOrder.OrderDetails.TotalAmount,
		&updatedOrder.OrderDetails.Currency,
		&updatedOrder.OrderDetails.Medium,
		&updatedOrder.OrderDetails.ShippingAddr,
		&updatedOrder.OrderDetails.BillingAddr,
		&updatedOrder.OrderDetails.Status,
		&updatedOrder.Metadata.CreatedAt,
		&updatedOrder.Metadata.UpdatedAt,
	)

	if err != nil {
		repo.log.Println("ERROR FETCHING UPDATED CART:", err)
		return entity.Order{}, fmt.Errorf("failed to fetch updated cart: %w", err)
	}

	repo.log.Println("Order updated successfully", "orderID", updatedOrder.ID)
	return updatedOrder, nil
}

func (repo PsqlRepo) CancelOrder(userID uuid.UUID, orderTypeID string) error {
	repo.log.Println("Canceling order for user ID:", userID, "Order Type ID:", orderTypeID)
	_, err := repo.db.Exec(`
		UPDATE erp.orders
		SET status = 'canceled'
		WHERE customer_id = $1 AND order_type_id = $2
	`, userID, orderTypeID)

	if err != nil {
		repo.log.Println("ERROR CANCELING ORDER:", err)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	repo.log.Println("Order canceled successfully", "userID", userID, "orderTypeID", orderTypeID)
	return nil
}

func (repo PsqlRepo) DeleteOrder(orderID uuid.UUID) error {
	_, err := repo.db.Exec(`
		DELETE FROM erp.orders
		WHERE id = $1
	`, orderID)
	if err != nil {
		repo.log.Println("DELETE ORDER ERROR:", err)
		return err
	}
	return nil
}

func (repo PsqlRepo) ListMerchantOrders(merchantID uuid.UUID) ([]entity.Order, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return nil, fmt.Errorf("invalid merchantID")
	}

	repo.log.Printf("Fetching orders for merchantID: %s", merchantID)
	query := `
        SELECT 
            o.id AS order_id, 
            o.customer_id, 
            o.order_type_id, 
            o.total_amount, 
            o.currency, 
            o.medium, 
            o.status, 
            o.shipping_address, 
            o.billing_address, 
            o.created_at, 
            o.updated_at
        FROM 
            erp.orders o
        JOIN 
            erp.order_items oi 
        ON 
            o.id = oi.order_id
        WHERE 
            oi.merchant_id = $1
        GROUP BY 
            o.id
    `

	rows, err := repo.db.Query(query, merchantID)
	if err != nil {
		repo.log.Printf("ERROR: Failed to fetch orders for merchantID %s: %v", merchantID, err)
		return nil, fmt.Errorf("failed to fetch orders for merchantID %s: %w", merchantID, err)
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID,
			&order.CustomerDetails.CustomerID,
			&order.OrderDetails.OrderTypeID,
			&order.OrderDetails.TotalAmount,
			&order.OrderDetails.Currency,
			&order.OrderDetails.Medium,
			&order.OrderDetails.Status,
			&order.OrderDetails.ShippingAddr,
			&order.OrderDetails.BillingAddr,
			&order.Metadata.CreatedAt,
			&order.Metadata.UpdatedAt,
		)
		if err != nil {
			repo.log.Printf("ERROR: Failed to scan order row for merchantID %s: %v", merchantID, err)
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		err = repo.fetchCustomerDetails(order.CustomerDetails.CustomerID, &order.CustomerDetails)
		if err != nil {
			repo.log.Printf("ERROR: Failed to fetch customer details for customerID %s: %v", order.CustomerDetails.CustomerID, err)
			return nil, fmt.Errorf("failed to fetch customer details: %w", err)
		}

		orderItems, err := repo.ListOrderItems(merchantID, order.ID.String())
		if err != nil {
			repo.log.Printf("ERROR: Failed to fetch order items for orderID %s: %v", order.ID, err)
			return nil, fmt.Errorf("failed to fetch order items for orderID %s: %w", order.ID, err)
		}
		order.OrderItems = orderItems

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		repo.log.Printf("ERROR: Row iteration error for merchantID %s: %v", merchantID, err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	repo.log.Printf("Successfully fetched %d orders for merchantID: %s", len(orders), merchantID)
	return orders, nil
}

func (repo PsqlRepo) fetchCustomerDetails(customerID uuid.UUID, details *entity.CustomerDetails) error {
	if customerID == uuid.Nil {
		repo.log.Println("WARNING: CustomerID is nil")
		return nil
	}

	userQuery := `
        SELECT first_name, last_name, sir_name
        FROM auth.users
        WHERE id = $1
    `
	var firstName, lastName, sirName string
	err := repo.db.QueryRow(userQuery, customerID).Scan(&firstName, &lastName, &sirName)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.log.Printf("No user found for Customer ID %s", customerID)
		} else {
			repo.log.Printf("Failed to fetch user for customer ID %s: %v", customerID, err)
			return err
		}
	}

	details.FirstName = firstName
	details.LastName = lastName
	details.SirName = sirName
	details.FullName = fmt.Sprintf("%s %s %s", sirName, firstName, lastName)
	return nil
}

func (repo PsqlRepo) ListOrderItems(merchantID uuid.UUID, orderID string) ([]entity.OrderItem, error) {
	// Validate input
	if merchantID == uuid.Nil || orderID == "" {
		repo.log.Println("ERROR: Invalid merchantID or orderID provided")
		return nil, fmt.Errorf("invalid merchantID or orderID")
	}

	repo.log.Printf("Fetching order items for orderID: %s and merchantID: %s", orderID, merchantID)

	// Adjust the query as necessary based on your actual schema
	query := `
        SELECT 
            product_id, 
            product_name, 
            quantity, 
            unit_price, 
            total_price, 
            category, 
            sku,
            merchant_id
        FROM 
            erp.order_items
        WHERE 
            order_id = $1 AND merchant_id = $2
    `

	// Query the database
	rows, err := repo.db.Query(query, orderID, merchantID)
	if err != nil {
		repo.log.Printf("ERROR: Failed to list order items for orderID %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to list order items for orderID %s: %w", orderID, err)
	}
	defer rows.Close()

	var orderItems []entity.OrderItem

	// Iterate over query results
	for rows.Next() {
		var orderItem entity.OrderItem
		err := rows.Scan(
			&orderItem.ProductID,
			&orderItem.ProductName,
			&orderItem.Quantity,
			&orderItem.UnitPrice,
			&orderItem.TotalPrice,
			&orderItem.Category,
			&orderItem.SKU,
			&orderItem.MerchantID, // Make sure this is added
		)
		if err != nil {
			repo.log.Printf("ERROR: Failed to scan order item for orderID %s: %v", orderID, err)
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		// Print orderItem to verify that it was scanned correctly
		repo.log.Printf("Fetched order item: %+v", orderItem)
		orderItems = append(orderItems, orderItem)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		repo.log.Printf("ERROR: Row iteration error for orderID %s: %v", orderID, err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	repo.log.Printf("Successfully fetched %d order items for orderID: %s", len(orderItems), orderID)
	return orderItems, nil
}

func (repo PsqlRepo) ListOrderTypes(merchantID uuid.UUID) ([]entity.OrderDetails, error) {
	var types []entity.OrderDetails

	rows, err := repo.db.Query(`
		SELECT order_type_id, total_amount, currency, status
		FROM erp.orders
		WHERE merchant_id = $1
	`, merchantID)
	if err != nil {
		repo.log.Println("LIST ORDER TYPES ERROR:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderDetails entity.OrderDetails
		if err := rows.Scan(&orderDetails.OrderTypeID, &orderDetails.TotalAmount, &orderDetails.Currency, &orderDetails.Status); err != nil {
			repo.log.Println("LIST ORDER TYPES SCAN ERROR:", err)
			return nil, err
		}
		types = append(types, orderDetails)
	}

	return types, nil
}

func (repo PsqlRepo) ListMerchantCustomers(merchantID uuid.UUID) ([]entity.CustomerDetails, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return nil, fmt.Errorf("invalid merchantID")
	}
	repo.log.Printf("Fetching customer details for merchantID: %s", merchantID)

	query := `
        SELECT 
            u.id,
            u.first_name,
            u.last_name,
            u.sir_name,
            u.gender,
            u.date_of_birth
        FROM 
            auth.users u
        WHERE 
            u.id = $1
    `

	rows, err := repo.db.Query(query, merchantID)
	if err != nil {
		repo.log.Printf("ERROR: Failed to fetch customer details for merchantID %s: %v", merchantID, err)
		return nil, fmt.Errorf("failed to fetch customer details for merchantID %s: %w", merchantID, err)
	}
	defer rows.Close()

	var customerDetailsList []entity.CustomerDetails
	for rows.Next() {
		var customer entity.CustomerDetails
		var id uuid.UUID
		var firstName, lastName, sirName sql.NullString
		var gender sql.NullString
		var dateOfBirth sql.NullTime

		err := rows.Scan(
			&id,
			&firstName,
			&lastName,
			&sirName,
			&gender,
			&dateOfBirth,
		)
		if err != nil {
			repo.log.Printf("ERROR: Failed to scan customer row for merchantID %s: %v", merchantID, err)
			return nil, fmt.Errorf("failed to scan customer row: %w", err)
		}

		customer.CustomerID = id
		if firstName.Valid {
			customer.FirstName = firstName.String
		}
		if lastName.Valid {
			customer.LastName = lastName.String
		}
		if sirName.Valid {
			customer.SirName = sirName.String
		}
		if gender.Valid {
			customer.Gender = gender.String
		}

		if dateOfBirth.Valid {
			customer.DateOfBirth = dateOfBirth.Time
		} else {
			customer.DateOfBirth = time.Time{}
		}

		customer.FullName = fmt.Sprintf("%s %s %s", customer.SirName, customer.FirstName, customer.LastName)

		customerDetailsList = append(customerDetailsList, customer)
	}

	if err := rows.Err(); err != nil {
		repo.log.Printf("ERROR: Row iteration error for merchantID %s: %v", merchantID, err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	repo.log.Printf("Successfully fetched %d customer details for merchantID: %s", len(customerDetailsList), merchantID)
	return customerDetailsList, nil
}

func (repo PsqlRepo) CountMerchantCustomers(merchantID uuid.UUID) (int, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return 0, fmt.Errorf("invalid merchantID")
	}

	repo.log.Printf("Counting unique customers for merchantID: %s", merchantID)
	query := `
        SELECT 
            COUNT(DISTINCT o.customer_id) AS total_customers
        FROM 
            erp.orders o
        JOIN 
            erp.order_items oi
        ON 
            o.id = oi.order_id
        WHERE 
            oi.merchant_id = $1
    `

	var totalCustomers int
	err := repo.db.QueryRow(query, merchantID).Scan(&totalCustomers)
	if err != nil {
		repo.log.Printf("ERROR: Failed to count unique customers for merchantID %s: %v", merchantID, err)
		return 0, fmt.Errorf("failed to count unique customers for merchantID %s: %w", merchantID, err)
	}

	repo.log.Printf("Successfully counted %d unique customers for merchantID: %s", totalCustomers, merchantID)
	return totalCustomers, nil
}

func (repo PsqlRepo) CountMerchantOrders(merchantID uuid.UUID) (int, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return 0, fmt.Errorf("invalid merchantID")
	}

	repo.log.Printf("Counting total orders for merchantID: %s", merchantID)

	query := `
        SELECT 
            COUNT(DISTINCT o.id) AS total_orders
        FROM 
            erp.orders o
        JOIN 
            erp.order_items oi 
        ON 
            o.id = oi.order_id
        WHERE 
            oi.merchant_id = $1
    `

	var totalOrders int
	err := repo.db.QueryRow(query, merchantID).Scan(&totalOrders)
	if err != nil {
		repo.log.Printf("ERROR: Failed to count total orders for merchantID %s: %v", merchantID, err)
		return 0, fmt.Errorf("failed to count total orders for merchantID %s: %w", merchantID, err)
	}

	repo.log.Printf("Successfully counted %d orders for merchantID: %s", totalOrders, merchantID)
	return totalOrders, nil
}

func (repo PsqlRepo) CountMerchantProducts(merchantID uuid.UUID) (int, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return 0, fmt.Errorf("invalid merchantID")
	}

	repo.log.Printf("Counting total products for merchantID: %s", merchantID)
	query := `
        SELECT 
            COUNT(*) AS total_products
        FROM 
            erp.products
        WHERE 
            merchant_id = $1
    `

	var totalProducts int
	err := repo.db.QueryRow(query, merchantID).Scan(&totalProducts)
	if err != nil {
		repo.log.Printf("ERROR: Failed to count total products for merchantID %s: %v", merchantID, err)
		return 0, fmt.Errorf("failed to count total products for merchantID %s: %w", merchantID, err)
	}

	repo.log.Printf("Successfully counted %d products for merchantID: %s", totalProducts, merchantID)
	return totalProducts, nil
}

func (repo PsqlRepo) GetOrder(orderID string, userId uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order

	rows, err := repo.db.Query(`
			SELECT id, customer_id, order_type_id, total_amount, currency, status, payment_status, payment_method, shipping_address, billing_address, created_at, updated_at
			FROM erp.orders
			WHERE id = $1 AND customer_id = $2
		`, orderID, userId)
	if err != nil {
		repo.log.Println("GET ORDER ERROR:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.CustomerDetails.CustomerID, &order.OrderDetails.OrderTypeID, &order.OrderDetails.TotalAmount,
			&order.OrderDetails.Currency, &order.OrderDetails.Status, &order.OrderDetails.PaymentStatus, &order.OrderDetails.PaymentMethod,
			&order.OrderDetails.ShippingAddr, &order.OrderDetails.BillingAddr, &order.Metadata.CreatedAt, &order.Metadata.UpdatedAt); err != nil {
			repo.log.Println("GET ORDER SCAN ERROR:", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo PsqlRepo) ListOrders() ([]entity.Order, error) {
	var orders []entity.Order

	rows, err := repo.db.Query(`
		SELECT id, customer_id, order_type_id, total_amount, currency, status, payment_status, payment_method, shipping_address, billing_address, created_at, updated_at
		FROM erp.orders
	`)
	if err != nil {
		repo.log.Println("LIST ORDERS ERROR:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.ID, &order.CustomerDetails.CustomerID, &order.OrderDetails.OrderTypeID, &order.OrderDetails.TotalAmount,
			&order.OrderDetails.Currency, &order.OrderDetails.Status, &order.OrderDetails.PaymentStatus, &order.OrderDetails.PaymentMethod,
			&order.OrderDetails.ShippingAddr, &order.OrderDetails.BillingAddr, &order.Metadata.CreatedAt, &order.Metadata.UpdatedAt); err != nil {
			repo.log.Println("LIST ORDERS SCAN ERROR:", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (repo PsqlRepo) CreateOrderType(typeName string, description string, merchantID uuid.UUID, createdBy uuid.UUID) (*entity.OrderType, error) {
	orderTypeID := uuid.New()
	createdAt := time.Now().Format(time.RFC3339)
	updatedAt := createdAt

	_, err := repo.db.Exec(`
        INSERT INTO erp.order_types (id, merchant_id, type_name, description, created_at, updated_at, created_by, updated_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, orderTypeID, merchantID, typeName, description, createdAt, updatedAt, createdBy, createdBy)
	if err != nil {
		repo.log.Println("CREATE ORDER TYPE ERROR:", err)
		return nil, err
	}

	return &entity.OrderType{
		ID:          orderTypeID,
		TypeName:    typeName,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}, nil
}

func (repo PsqlRepo) UpdateOrderType(
	orderTypeID uuid.UUID,
	name string,
	status string,
	createdBy uuid.UUID,
	updatedBy uuid.UUID,
) error {
	_, err := repo.db.Exec(`
		UPDATE erp.order_types
		SET type_name = $1, status = $2, updated_at = $3, created_by = $4, updated_by = $5
		WHERE id = $6
	`, name, status, time.Now(), createdBy, updatedBy, orderTypeID)

	if err != nil {
		repo.log.Println("UPDATE ORDER TYPE ERROR:", err)
		return err
	}
	return nil
}

func (repo PsqlRepo) DeactivateOrderType(userId uuid.UUID, orderTypeID string) error {
	const ErrFailedToDeactivateOrderType = "FAILED_TO_DEACTIVATE_ORDER_TYPE"

	_, err := repo.db.Exec(`
		UPDATE erp.order_types
		SET status = $1, updated_by = $2
		WHERE id = $3
	`, "inactive", userId, orderTypeID)
	if err != nil {
		repo.log.Println("ERROR DEACTIVATING ORDER TYPE:", err)
		return fmt.Errorf("%s: %v", ErrFailedToDeactivateOrderType, err)
	}

	repo.log.Println("ORDER TYPE DEACTIVATED SUCCESSFULLY")
	return nil
}

func (repo PsqlRepo) UpdateOrderItem(orderID uuid.UUID, itemID uuid.UUID, userId uuid.UUID, quantity int, price float64, discount float64, tax float64) error {
	_, err := repo.db.Exec(`
		UPDATE erp.order_items
		SET quantity = $1, unit_price = $2, total_price = $3, category = $4, sku = $5, updated_at = $6
		WHERE order_id = $7 AND id = $8 AND user_id = $9
	`, quantity, price, discount, tax, time.Now(), orderID, itemID, userId)
	if err != nil {
		repo.log.Println("UPDATE ORDER ITEM ERROR:", err)
		return err
	}

	return nil
}

func (repo PsqlRepo) DeleteOrderItem(orderID uuid.UUID, itemID uuid.UUID) error {
	_, err := repo.db.Exec(`
		DELETE FROM erp.order_items
		WHERE order_id = $1 AND id = $2
	`, orderID, itemID)
	if err != nil {
		repo.log.Println("DELETE ORDER ITEM ERROR:", err)
		return err
	}

	return nil
}
func (repo PsqlRepo) RemoveOrderItem(orderID, itemID, userID uuid.UUID) error {
	repo.log.Println("Removing order item", "orderID", orderID, "itemID", itemID, "userID", userID)
	_, err := repo.db.Exec(`
        DELETE FROM erp.order_items
        WHERE order_id = $1 AND item_id = $2 AND user_id = $3
    `, orderID, itemID, userID)

	if err != nil {
		repo.log.Println("REMOVE ORDER ITEM ERROR:", err)
		return fmt.Errorf("failed to remove order item: %w", err)
	}

	repo.log.Println("Order item removed successfully", "orderID", orderID, "itemID", itemID, "userID", userID)
	return nil
}
