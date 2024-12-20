package psql

import (
	"auth/src/pkg/erp/core/entity"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func (repo PsqlRepo) CreateMerchantInvoice(merchantID uuid.UUID, orderID uuid.UUID) (*entity.Order, error) {
	if merchantID == uuid.Nil {
		repo.log.Println("ERROR: Invalid merchantID provided")
		return nil, fmt.Errorf("invalid merchantID")
	}
	if orderID == uuid.Nil {
		repo.log.Println("ERROR: Invalid orderID provided")
		return nil, fmt.Errorf("invalid orderID")
	}

	repo.log.Printf("Creating Invoice for orderID: %s and merchantID: %s", orderID, merchantID)
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
        WHERE 
            o.id = $1
    `

	rows, err := repo.db.Query(query, orderID)
	if err != nil {
		repo.log.Printf("ERROR: Failed to fetch order data for orderID %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to fetch order data: %w", err)
	}
	defer rows.Close()

	var order *entity.Order
	if rows.Next() {
		order = &entity.Order{}
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
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Fetch additional customer details and order items
		err = repo.InvoiceCustomerDetails(order.CustomerDetails.CustomerID, &order.CustomerDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch customer details: %w", err)
		}

		orderItems, err := repo.InvoiceOrderItems(merchantID, order.ID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order items: %w", err)
		}
		order.OrderItems = orderItems
	} else {
		repo.log.Printf("INFO: No order found for orderID %s", orderID)
	}

	if order == nil {
		return nil, fmt.Errorf("no orders found for orderID %s", orderID)
	}

	repo.log.Println("Merchant Invoice created successfully", "orderID", orderID)
	return order, nil
}

func (repo PsqlRepo) InvoiceCustomerDetails(customerID uuid.UUID, details *entity.CustomerDetails) error {
	if customerID == uuid.Nil {
		repo.log.Println("WARNING: CustomerID is nil")
		return nil
	}

	userQuery := `
        SELECT first_name, last_name, sir_name
        FROM auth.users
        WHERE id = $1
    `
	var firstName, lastName, sirName sql.NullString
	err := repo.db.QueryRow(userQuery, customerID).Scan(&firstName, &lastName, &sirName)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.log.Printf("INFO: No user found for Customer ID %s", customerID)
		} else {
			repo.log.Printf("ERROR: Failed to fetch user for customer ID %s: %v", customerID, err)
			return fmt.Errorf("failed to fetch user for customer ID %s: %w", customerID, err)
		}
	}

	details.FirstName = firstName.String
	details.LastName = lastName.String
	details.SirName = sirName.String
	details.FullName = fmt.Sprintf("%s %s %s", details.SirName, details.FirstName, details.LastName) // Corrected order

	phoneQuery := `
        SELECT prefix, number
        FROM auth.phones
        WHERE id = $1
    `
	var prefix, number string
	err = repo.db.QueryRow(phoneQuery, customerID).Scan(&prefix, &number)
	if err != nil {
		if err == sql.ErrNoRows {
			repo.log.Printf("INFO: No phone found for customerID %s", customerID)
		} else {
			repo.log.Printf("ERROR: Failed to fetch phone for customerID %s: %v", customerID, err)
			return fmt.Errorf("failed to fetch phone for customerID %s: %w", customerID, err)
		}
	}
	if prefix != "" && number != "" {
		details.PhoneNumber = fmt.Sprintf("%s%s", prefix, number)
	}

	return nil
}

func (repo PsqlRepo) InvoiceOrderItems(merchantID uuid.UUID, orderID string) ([]entity.OrderItem, error) {
	if merchantID == uuid.Nil || orderID == "" {
		repo.log.Println("ERROR: Invalid merchantID or orderID provided")
		return nil, fmt.Errorf("invalid merchantID or orderID")
	}

	repo.log.Printf("Fetching order items for orderID: %s and merchantID: %s", orderID, merchantID)

	query := `
        SELECT 
            product_id, 
            product_name, 
            quantity, 
            unit_price, 
            total_price, 
            category, 
            sku
        FROM 
            erp.order_items
        WHERE 
            order_id = $1 AND merchant_id = $2
    `

	rows, err := repo.db.Query(query, orderID, merchantID)
	if err != nil {
		repo.log.Printf("ERROR: Failed to list order items for orderID %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to list order items for orderID %s: %w", orderID, err)
	}
	defer rows.Close()

	var orderItems []entity.OrderItem

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
		)
		if err != nil {
			repo.log.Printf("ERROR: Failed to scan order item for orderID %s: %v", orderID, err)
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		repo.log.Printf("ERROR: Row iteration error for orderID %s: %v", orderID, err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if len(orderItems) == 0 {
		repo.log.Printf("INFO: No order items found for orderID: %s", orderID)
	}

	repo.log.Printf("Successfully fetched %d order items for orderID: %s", len(orderItems), orderID)
	return orderItems, nil
}
