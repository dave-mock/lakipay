package usecase

import (
	"auth/src/pkg/erp/core/entity"

	"github.com/google/uuid"
)

type Interactor interface {
	// Catalogs
	CreateCatalog(userId uuid.UUID, name string, description string, status string) (*entity.Catalog, error)
	ListMerchantCatalogs(merchantID string, userId uuid.UUID) ([]entity.Catalog, error)
	ListCatalogs(userId uuid.UUID) ([]entity.Catalog, error)
	GetCatalog(catalogID string, userId uuid.UUID) (*entity.Catalog, error)
	ArchiveCatalog(catalogID string, userId uuid.UUID) error
	UpdateCatalog(userId uuid.UUID, catalogID, name, description, status string) (*entity.Catalog, error)
	DeleteCatalog(catalogID string, userId uuid.UUID) error

	// Sub Catalogs

	CreateSubCatalog(subCatalog *entity.SubCatalog) error
	ListSubCatalogs(userId uuid.UUID) ([]entity.SubCatalog, error)
	GetSubCatalog(SubcatalogID uuid.UUID, userId uuid.UUID) (*entity.SubCatalog, error)
	ArchiveSubCatalog(SubcatalogID uuid.UUID, userId uuid.UUID) error
	UpdateSubCatalog(uuid.UUID, *entity.SubCatalog) (*entity.SubCatalog, error)
	DeleteSubCatalog(SubcatalogID uuid.UUID, userId uuid.UUID) error

	// Customer Management
	CreateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, createdBy, merchantID uuid.UUID) (*entity.Customer, error)
	UpdateCustomer(customerID uuid.UUID, name, email, phone, address string, loyaltyPoints int, dateOfBirth, status string, updatedBy, merchantID uuid.UUID) (*entity.Customer, error)

	// Product Management
	CreateProduct(name string, description string, price float64, currency string, SKU string, weight float64, dimensions string, imageURL string, userId uuid.UUID) error
	GetProduct(productId uuid.UUID, userId uuid.UUID) (*entity.Product, error)
	UpdateProduct(
		productIDParsed uuid.UUID,
		name string,
		description string,
		price float64,
		currency string,
		SKU string,
		weight float64,
		dimensions string,
		imageURL string,
		status string,
		userId uuid.UUID,
	) error
	DeleteProduct(UpdateproductID uuid.UUID, userId uuid.UUID) error
	ListProducts(userId uuid.UUID, userType string) ([]entity.Product, error)
	ListAllProducts() ([]entity.Product, error)
	ArchiveProduct(productId uuid.UUID) error
	AddProductToCatalog(NewproductID uuid.UUID, NewCatalogID uuid.UUID, displayOrder int, userId uuid.UUID) ([]entity.ProductCatalog, error)
	RemoveProductFromCatalog(productID uuid.UUID, catalogID uuid.UUID, userId uuid.UUID) ([]entity.ProductCatalog, error)

	// Warehouse Management
	CreateWarehouse(name string, location string, capacity int, description string, userId uuid.UUID) (*entity.Warehouse, error)
	GetWarehouse(GetWarehouseID uuid.UUID, userId uuid.UUID) (*entity.Warehouse, error)
	UpdateWarehouse(
		UpdateWarehouseID uuid.UUID,
		name string,
		location string,
		capacity int,
		description string,
		userId uuid.UUID,
	) (*entity.Warehouse, error)
	DeleteWarehouse(WarehouseID uuid.UUID, userId uuid.UUID) error
	ListWarehouses(userId uuid.UUID) ([]entity.Warehouse, error)
	ListMerchantWarehouses(userId uuid.UUID) ([]entity.Warehouse, error)
	DeactivateWarehouse(DeactivateWarehouseID uuid.UUID, userId uuid.UUID) (*entity.Warehouse, error)
	CreatePaymentMethod(name string, methodType string, commission float64, details string, isActive bool, userId uuid.UUID) error
	UpdatePaymentMethod(id uuid.UUID, name string, methodType string, commission float64, details string, isActive *bool, userId uuid.UUID) error
	ListPaymentMethods(userId uuid.UUID) ([]entity.PaymentMethod, error)
	GetPaymentMethod(id uuid.UUID, userId uuid.UUID) (*entity.PaymentMethod, error)
	DeactivatePaymentMethod(id string, userId uuid.UUID) error
	CreateMerchantInvoice(merchantID uuid.UUID, orderID uuid.UUID) (*entity.Order, error)
	ListMerchantOrders(userId uuid.UUID) ([]entity.Order, error)
	CountMerchantCustomers(userId uuid.UUID) (int, error)
	ListMerchantCustomers(uuid.UUID) ([]entity.CustomerDetails, error)
	CountMerchantOrders(userId uuid.UUID) (int, error)
	CountMerchantProducts(userId uuid.UUID) (int, error)

	// Order Management
	CreateOrder(
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
	) (*entity.Order, error)

	UpdateOrder(
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
	) (*entity.Order, error)
	CreateCartOrder(
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
	) (*entity.Order, error)

	UpdateCartOrder(
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
	) (*entity.Order, error)
	ListOrderItems(userId uuid.UUID, orderID string) ([]entity.OrderItem, error)
	CancelOrder(userId uuid.UUID, orderTypeID string) error
	ListOrders() ([]entity.Order, error)
	GetOrder(orderID string, userId uuid.UUID) ([]entity.Order, error)
	UpdateOrderItem(orderID uuid.UUID, itemID uuid.UUID, userId uuid.UUID, quantity int, price, discount, tax float64) error
	RemoveOrderItem(orderID uuid.UUID, itemID uuid.UUID, userId uuid.UUID) error
}
