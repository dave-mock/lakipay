package basic

import "auth/src/pkg/checkout/core/entity"

func (service BasicCheckoutService) GetGateways() ([]entity.Gateway, error) {
	var gateways []entity.Gateway = make([]entity.Gateway, 0)

	gateways, err := service.repo.FindGateways()

	return gateways, err
}

func (service BasicCheckoutService) FindGatewayByKey(key string) (*entity.Gateway, error) {
	return service.repo.FindGatewayByKey(key)
}
