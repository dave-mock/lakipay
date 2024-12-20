package rest

import (
	auth "auth/src/pkg/auth/adapter/controller/procedure"
	"auth/src/pkg/erp/usecase"
	"encoding/json"
	"log"
	"net/http"
)

// Controller struct
type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
}

func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth, sm: sm}
	// Merchant Management
	/* 	sm.HandleFunc("/merchants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateMerchant(w, r)
		}
	}) */
	/* 	sm.HandleFunc("/merchants/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListMerchants(w, r)
		case http.MethodPut:
			controller.UpdateMerchant(w, r)
		case http.MethodDelete:
			controller.DeactivateMerchant(w, r)
		}
	}) */

	sm.HandleFunc("/merchants/warehouses/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListMerchantWarehouses(w, r)
		}
	})
	sm.HandleFunc("/merchants/catalogs/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListMerchantCatalogs(w, r)
		}
	})
	/* 	sm.HandleFunc("/merchants/{id}/customers", func(w http.ResponseWriter, r *http.Request) {
	   		switch r.Method {
	   		case http.MethodGet:
	   			controller.ListMerchantCustomers(w, r)
	   		}
	   	})
	*/
	sm.HandleFunc("/merchants/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListMerchantOrders(w, r)
		}
	})

	sm.HandleFunc("/merchants/tatal/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.CountMerchantOrders(w, r)
		}
	})
	sm.HandleFunc("/merchants/invoice/fetch/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateMerchantInvoice(w, r)
		}
	})
	sm.HandleFunc("/merchants/total/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.CountMerchantProducts(w, r)
		}
	})

	// Warehouse Management
	sm.HandleFunc("/warehouses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateWarehouse(w, r)
		case http.MethodGet:
			controller.ListWarehouses(w, r)
		}
	})
	sm.HandleFunc("/warehouses/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetWarehouse(w, r)
		case http.MethodPut:
			controller.UpdateWarehouse(w, r)
		case http.MethodDelete:
			controller.DeleteWarehouse(w, r)
		}
	})

	// Catalog Management
	sm.HandleFunc("/merchants/catalogs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateCatalog(w, r)
		case http.MethodGet:
			controller.ListCatalogs(w, r)
		}
	})
	sm.HandleFunc("/catalogs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetCatalog(w, r)
		case http.MethodPut:
			controller.UpdateCatalog(w, r)
		case http.MethodDelete:
			controller.DeleteCatalog(w, r)
		}
	})

	// Catalog Management
	sm.HandleFunc("/merchants/sub-catalogs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateSubCatalog(w, r)
		case http.MethodGet:
			controller.ListSubCatalogs(w, r)
		}
	})
	sm.HandleFunc("/merchants/sub-catalogs/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetSubCatalog(w, r)
		case http.MethodPut:
			controller.UpdateSubCatalog(w, r)
		case http.MethodDelete:
			controller.DeleteSubCatalog(w, r)
		}
	})
	sm.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateProduct(w, r)
		case http.MethodGet:
			controller.ListProducts(w, r)
		}
	})
	sm.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetProduct(w, r)
		case http.MethodPut:
			controller.UpdateProduct(w, r)
		case http.MethodDelete:
			controller.DeleteProduct(w, r)
		}
	})
	sm.HandleFunc("/catalogs/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.AddProductToCatalog(w, r)
		case http.MethodDelete:
			controller.RemoveProductFromCatalog(w, r)
		}
	})

	// CRM
	sm.HandleFunc("/customers/count/all", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateCustomer(w, r)
		case http.MethodGet:
			controller.CountMerchantCustomers(w, r)
		}
	})
	sm.HandleFunc("/customers/list/all", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			controller.UpdateCustomer(w, r)
		case http.MethodGet:
			controller.GetMerchantCustomers(w, r)
		}
	})

	// Finance Management
	sm.HandleFunc("/payment-methods", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreatePaymentMethod(w, r)
		case http.MethodGet:
			controller.ListPaymentMethods(w, r)
		}
	})
	sm.HandleFunc("/payment-methods/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetPaymentMethod(w, r)
		case http.MethodPut:
			controller.UpdatePaymentMethod(w, r)
		case http.MethodDelete:
			controller.DeactivatePaymentMethod(w, r)
		}
	})

	// Order Management

	sm.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateOrder(w, r)
		case http.MethodGet:
			controller.ListOrders(w, r)
		}
	})
	sm.HandleFunc("/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetOrder(w, r)
		case http.MethodPut:
			controller.UpdateOrder(w, r)
		case http.MethodDelete:
			controller.CancelOrder(w, r)
		}
	})

	sm.HandleFunc("/cart/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateCartOrder(w, r)
		case http.MethodGet:
			controller.ListCartOrders(w, r)
		}
	})
	sm.HandleFunc("/cart/orders/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetCartOrder(w, r)
		case http.MethodPut:
			controller.UpdateCartOrder(w, r)
		case http.MethodDelete:
			controller.CancelCartOrder(w, r)
		}
	})
	sm.HandleFunc("/orders/{id}/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListOrderItems(w, r)
		}
	})

	sm.HandleFunc("/orders/{order_id}/items/{item_id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			controller.UpdateOrderItem(w, r)
		case http.MethodDelete:
			controller.RemoveOrderItem(w, r)
		}
	})

	/* // Invoices
	sm.HandleFunc("/invoices", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateInvoice(w, r)
		case http.MethodGet:
			controller.ListInvoices(w, r)
		}
	})
	sm.HandleFunc("/invoices/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetInvoice(w, r)
		case http.MethodPut:
			controller.UpdateInvoice(w, r)
		case http.MethodDelete:
			controller.CancelInvoice(w, r)
		}
	}) */

	return controller
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, data Response, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}
