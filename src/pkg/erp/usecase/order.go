package usecase

import (
	"auth/src/pkg/erp/core/entity"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ErrFailedToListMerchantOrders = "Failed to list merchant orders"
	ErrFailedToCreateOrderType    = "Failed to create order type"
	ErrFailedToListOrderTypes     = "Failed to list order types"
	ErrFailedToCreateOrder        = "Failed to create order"
	ErrFailedToUpdateOrder        = "Failed to update order"
	ErrFailedToListOrderItems     = "Failed to list order items"
	ErrFailedToGetOrderType       = "Failed to get order type"
	ErrFailedToUpdateOrderType    = "Failed to update order type"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func (uc Usecase) CreateOrder(
	merchantID uuid.UUID,
	customerID string,
	orderTypeID string,
	totalAmount float64,
	currency string,
	medium string,
	shippingAddress string,
	billingAddress string,
	orderItems []entity.OrderItem,
	discounts []entity.Discount,
	taxes []entity.Tax,
) (*entity.Order, error) {

	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		uc.log.Println("Error: Invalid customerID:", err)
		return nil, fmt.Errorf("invalid customerID: %w", err)
	}

	parsedOrderTypeID, err := uuid.Parse(orderTypeID)
	if err != nil {
		uc.log.Println("Error: Invalid orderTypeID:", err)
		return nil, fmt.Errorf("invalid orderTypeID: %w", err)
	}

	if err := validateOrderInput(merchantID, parsedCustomerID, parsedOrderTypeID, totalAmount, currency, medium, shippingAddress, billingAddress); err != nil {
		uc.log.Println("Error: Failed to validate order input:", err)
		return nil, fmt.Errorf("failed to validate order input: %w", err)
	}

	order := &entity.Order{
		ID: uuid.New(),
		CustomerDetails: entity.CustomerDetails{
			CustomerID: parsedCustomerID,
		},
		OrderDetails: entity.OrderDetails{
			OrderTypeID:  parsedOrderTypeID,
			TotalAmount:  totalAmount,
			Currency:     currency,
			Medium:       medium,
			ShippingAddr: shippingAddress,
			BillingAddr:  billingAddress,
			Discounts:    discounts,
			Taxes:        taxes,
			FinalAmount:  totalAmount,
		},
		OrderItems: orderItems,
		Metadata: entity.Metadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdOrder, err := uc.repo.CreateOrder(*order)
	if err != nil {
		uc.log.Println("Error: Failed to create order in repository:", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	uc.log.Println("Order created successfully:", createdOrder.ID)
	return createdOrder, nil
}

func (uc Usecase) UpdateOrder(
	orderID string,
	merchantID uuid.UUID,
	customerID uuid.UUID,
	orderTypeID string,
	totalAmount float64,
	currency string,
	medium string,
	shippingAddress string,
	billingAddress string,
	orderItems []entity.OrderItem,
) (*entity.Order, error) {
	uc.log.Println("Updating order", "orderID", orderID)
	order := entity.Order{
		ID: uuid.MustParse(orderID),
		CustomerDetails: entity.CustomerDetails{
			CustomerID: customerID,
		},
		OrderDetails: entity.OrderDetails{
			OrderTypeID:  uuid.MustParse(orderTypeID),
			TotalAmount:  totalAmount,
			Currency:     currency,
			Medium:       medium,
			ShippingAddr: shippingAddress,
			BillingAddr:  billingAddress,
			FinalAmount:  totalAmount,
		},
		OrderItems: orderItems,
		Metadata: entity.Metadata{
			UpdatedAt: time.Now(),
		},
	}

	updatedOrder, err := uc.repo.UpdateOrder(order)
	if err != nil {
		uc.log.Println("Error updating order", "error", err)
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	uc.log.Println("Order updated successfully", "orderID", updatedOrder.ID)
	return &updatedOrder, nil
}

func (uc Usecase) CreateCartOrder(
	merchantID uuid.UUID,
	customerID string,
	orderTypeID string,
	totalAmount float64,
	currency string,
	medium string,
	shippingAddress string,
	billingAddress string,
	orderItems []entity.OrderItem,
	discounts []entity.Discount,
	taxes []entity.Tax,
) (*entity.Order, error) {

	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		uc.log.Println("Error: Invalid customerID:", err)
		return nil, fmt.Errorf("invalid customerID: %w", err)
	}

	parsedOrderTypeID, err := uuid.Parse(orderTypeID)
	if err != nil {
		uc.log.Println("Error: Invalid orderTypeID:", err)
		return nil, fmt.Errorf("invalid orderTypeID: %w", err)
	}

	if err := validateOrderInput(merchantID, parsedCustomerID, parsedOrderTypeID, totalAmount, currency, medium, shippingAddress, billingAddress); err != nil {
		uc.log.Println("Error: Failed to validate order input:", err)
		return nil, fmt.Errorf("failed to validate order input: %w", err)
	}

	order := &entity.Order{
		ID: uuid.New(),
		CustomerDetails: entity.CustomerDetails{
			CustomerID: parsedCustomerID,
		},
		OrderDetails: entity.OrderDetails{
			OrderTypeID:  parsedOrderTypeID,
			TotalAmount:  totalAmount,
			Currency:     currency,
			Medium:       medium,
			ShippingAddr: shippingAddress,
			BillingAddr:  billingAddress,
			Discounts:    discounts,
			Taxes:        taxes,
			FinalAmount:  totalAmount,
		},
		OrderItems: orderItems,
		Metadata: entity.Metadata{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	createdOrder, err := uc.repo.CreateCartOrder(*order)
	if err != nil {
		uc.log.Println("Error: Failed to create order in repository:", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	uc.log.Println("|||||    Order added to cart |||||||", createdOrder.ID)
	return createdOrder, nil
}

func (uc Usecase) UpdateCartOrder(
	orderID string,
	merchantID uuid.UUID,
	customerID uuid.UUID,
	orderTypeID string,
	totalAmount float64,
	currency string,
	medium string,
	shippingAddress string,
	billingAddress string,
	orderItems []entity.OrderItem,
) (*entity.Order, error) {
	uc.log.Println("Updating order", "orderID", orderID)
	order := entity.Order{
		ID: uuid.MustParse(orderID),
		CustomerDetails: entity.CustomerDetails{
			CustomerID: customerID,
		},
		OrderDetails: entity.OrderDetails{
			OrderTypeID:  uuid.MustParse(orderTypeID),
			TotalAmount:  totalAmount,
			Currency:     currency,
			Medium:       medium,
			ShippingAddr: shippingAddress,
			BillingAddr:  billingAddress,
			FinalAmount:  totalAmount,
		},
		OrderItems: orderItems,
		Metadata: entity.Metadata{
			UpdatedAt: time.Now(),
		},
	}

	updatedOrder, err := uc.repo.UpdateCartOrder(order)
	if err != nil {
		uc.log.Println("Error updating order", "error", err)
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	uc.log.Println("Cart updated ||||||| ||", "orderID", updatedOrder.ID)
	return &updatedOrder, nil
}

func (uc Usecase) ListMerchantOrders(userId uuid.UUID) ([]entity.Order, error) {
	uc.log.Println("Fetching list of merchant orders", "userId", userId)

	orders, err := uc.repo.ListMerchantOrders(userId)
	if err != nil {
		uc.log.Println(ErrFailedToListMerchantOrders, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListMerchantOrders, err)
	}

	uc.log.Println("Merchant orders fetched successfully", "userId", userId)
	return orders, nil
}

func (uc Usecase) CountMerchantOrders(userId uuid.UUID) (int, error) {
	uc.log.Println("Fetching total count of merchant orders", "userId", userId)

	totalOrders, err := uc.repo.CountMerchantOrders(userId)
	if err != nil {
		uc.log.Println(ErrFailedToListMerchantOrders, err)
		return 0, fmt.Errorf("%s: %w", ErrFailedToListMerchantOrders, err)
	}

	uc.log.Println("Total Merchant orders fetched successfully", "userId", userId)
	return totalOrders, nil
}
func (uc Usecase) CountMerchantCustomers(userId uuid.UUID) (int, error) {
	uc.log.Println("Fetching total count of merchant orders", "userId", userId)

	totalOrders, err := uc.repo.CountMerchantCustomers(userId)
	if err != nil {
		uc.log.Println(ErrFailedToListMerchantOrders, err)
		return 0, fmt.Errorf("%s: %w", ErrFailedToListMerchantOrders, err)
	}

	uc.log.Println("Total Merchant orders fetched successfully", "userId", userId)
	return totalOrders, nil
}

func (uc Usecase) ListMerchantCustomers(merchantID uuid.UUID) ([]entity.CustomerDetails, error) {
	uc.log.Println("Fetching customer details for merchant", "merchantID", merchantID)

	customerDetails, err := uc.repo.ListMerchantCustomers(merchantID)
	if err != nil {
		uc.log.Println("Failed to fetch customer details", err)
		return nil, fmt.Errorf("failed to fetch customer details: %w", err)
	}

	if len(customerDetails) == 0 {
		uc.log.Println("No customer details found")
		return nil, nil
	}

	uc.log.Println("Customer details fetched successfully", "merchantID", merchantID)
	return customerDetails, nil
}

func (uc Usecase) CountMerchantProducts(userId uuid.UUID) (int, error) {
	uc.log.Println("Fetching total count of merchant orders", "userId", userId)

	totalOrders, err := uc.repo.CountMerchantProducts(userId)
	if err != nil {
		uc.log.Println(ErrFailedToListMerchantOrders, err)
		return 0, fmt.Errorf("%s: %w", ErrFailedToListMerchantOrders, err)
	}

	uc.log.Println("Total Merchant orders fetched successfully", "userId", userId)
	return totalOrders, nil
}

func (uc Usecase) ListOrderItems(userId uuid.UUID, orderID string) ([]entity.OrderItem, error) {
	uc.log.Println("Listing order items")
	orderItems, err := uc.repo.ListOrderItems(userId, orderID)
	if err != nil {
		uc.log.Println(ErrFailedToListOrderItems, err)
		return nil, fmt.Errorf("%s: %w", ErrFailedToListOrderItems, err)
	}

	uc.log.Println("Order items listed successfully")
	return orderItems, nil
}

func validateInput(name string, description string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if description == "" {
		return errors.New("description cannot be empty")
	}
	return nil
}
func validateOrderInput(merchantID, customerID uuid.UUID, orderTypeID uuid.UUID, totalAmount float64, currency, medium, shippingAddress, billingAddress string) error {
	if orderTypeID == uuid.Nil {
		return errors.New("order type ID is required")
	}
	if totalAmount <= 0 {
		return errors.New("total amount must be greater than zero")
	}
	if currency == "" {
		return errors.New("currency is required")
	}
	if medium == "" {
		return errors.New("medium is required")
	}
	return nil
}

func (uc Usecase) CancelOrder(userId uuid.UUID, orderTypeID string) error {
	uc.log.Println("Canceling order", "orderTypeID", orderTypeID, "userId", userId)
	err := uc.repo.CancelOrder(
		userId,
		orderTypeID)
	if err != nil {
		uc.log.Println("Error canceling order", "error", err)
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	uc.log.Println("Order canceled successfully", "orderTypeID", orderTypeID)
	return nil
}
func (uc Usecase) ListOrders() ([]entity.Order, error) {
	uc.log.Println("Listing all orders")
	orders, err := uc.repo.ListOrders()
	if err != nil {
		uc.log.Println("Error fetching orders", "error", err)
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	uc.log.Println("Orders fetched successfully")
	return orders, nil
}
func (uc Usecase) GetOrder(orderID string, userId uuid.UUID) ([]entity.Order, error) {
	uc.log.Println("Fetching order details", "orderID", orderID, "userId", userId)
	order, err := uc.repo.GetOrder(
		orderID,
		userId)
	if err != nil {
		uc.log.Println("Error fetching order", "error", err)
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}

	uc.log.Println("Order fetched successfully", "orderID", orderID)
	return order, nil
}

//   ||||||||||||

func (uc Usecase) UpdateOrderItem(orderID uuid.UUID, itemID uuid.UUID, userId uuid.UUID, quantity int, price float64, discount float64, tax float64) error {
	uc.log.Println("Updating order item", "orderID", orderID, "itemID", itemID)

	if err := validateOrderItemInput(
		userId, itemID,
		"",
		quantity,
		price,
		discount,
		tax); err != nil {
		uc.log.Println("Validation failed for order item update", "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}

	err := uc.repo.UpdateOrderItem(
		orderID,
		itemID,
		userId,
		quantity,
		price,
		discount, tax)
	if err != nil {
		uc.log.Println("Error updating order item", "error", err)
		return fmt.Errorf("failed to update order item: %w", err)
	}

	uc.log.Println("Order item updated successfully", "orderID", orderID, "itemID", itemID)
	return nil
}
func (uc Usecase) RemoveOrderItem(orderID uuid.UUID, itemID uuid.UUID, userId uuid.UUID) error {
	uc.log.Println("Removing order item", "orderID", orderID, "itemID", itemID)
	if err := validateOrderItemInput(
		userId,
		itemID,
		"",
		0,
		0.0,
		0.0,
		0.0); err != nil {
		uc.log.Println("Validation failed for order item removal", "error", err)
		return fmt.Errorf("validation failed: %w", err)
	}
	err := uc.repo.RemoveOrderItem(orderID, itemID, userId)
	if err != nil {
		uc.log.Println("Error removing order item", "error", err)
		return fmt.Errorf("failed to remove order item: %w", err)
	}

	uc.log.Println("Order item removed successfully", "orderID", orderID, "itemID", itemID)
	return nil
}

func validateUpdateOrderInput(orderID uuid.UUID, merchantID, customerID uuid.UUID, orderTypeID uuid.UUID, totalAmount float64, currency, medium, shippingAddress, billingAddress string) error {
	if orderTypeID == uuid.Nil {
		return errors.New("order type ID is required")
	}
	if totalAmount <= 0 {
		return errors.New("total amount must be greater than zero")
	}
	if currency == "" {
		return errors.New("currency is required")
	}
	if medium == "" {
		return errors.New("medium is required")
	}
	return nil
}
func validateOrderItemInput(userId uuid.UUID, itemOrderID uuid.UUID, productID string, quantity int, price, discount, tax float64) error {
	if userId == uuid.Nil {
		return errors.New("user ID is required")
	}

	if productID == "" {
		return errors.New("product ID is required")
	}
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if discount < 0 {
		return errors.New("discount cannot be negative")
	}
	if tax < 0 {
		return errors.New("tax cannot be negative")
	}
	return nil
}
