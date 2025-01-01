package entity

import (
	"time"

	"github.com/google/uuid"
)


type CustomerDetails struct {
    CustomerID  uuid.UUID `json:"customer_id"`
    FirstName   string    `json:"first_name"`
    LastName    string    `json:"last_name"`
    SirName     string    `json:"sir_name"`
    FullName    string    `json:"full_name"` 
    Gender      string    `json:"gender"`
    DateOfBirth time.Time `json:"date_of_birth"`
    PhoneNumber string    `json:"phone_number"`
}

type Discount struct {
	Type        string  `json:"type"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
}

type Tax struct {
	Type  string  `json:"type"`
	Rate  float64 `json:"rate"`
	Value float64 `json:"value"`
}

type OrderItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	Category    string    `json:"category"`
	SKU         string    `json:"sku"`
	TotalPrice  float64   `json:"total_price"`
}

type OrderDetails struct {
	OrderTypeID   uuid.UUID  `json:"order_type_id"`
	TotalAmount   float64    `json:"total_amount"`
	Currency      string     `json:"currency"`
	Medium        string     `json:"medium"`
	Status        string     `json:"status"`
	PaymentStatus string     `json:"payment_status"`
	PaymentMethod string     `json:"payment_method"`
	PaymentRef    string     `json:"payment_reference"`
	ShippingAddr  string     `json:"shipping_address"`
	BillingAddr   string     `json:"billing_address"`
	Discounts     []Discount `json:"discounts"`
	Taxes         []Tax      `json:"taxes"`
	FinalAmount   float64    `json:"final_amount"`
}

type Metadata struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Notes     string    `json:"notes"`
}

type Tracking struct {
	Status               string `json:"status"`
	ExpectedDeliveryDate string `json:"expected_delivery_date"`
	ShipmentID           string `json:"shipment_id"`
}

type Order struct {
	ID              uuid.UUID       `json:"id"`
	CustomerDetails CustomerDetails `json:"customer_details"`
	OrderDetails    OrderDetails    `json:"order_details"`
	OrderItems      []OrderItem     `json:"order_items"`
	Metadata        Metadata        `json:"metadata"`
	Tracking        Tracking        `json:"tracking"`
}
