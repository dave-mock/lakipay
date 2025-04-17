package psql

import "auth/src/pkg/checkout/core/entity"

func (repo CheckoutPSQLRepo) FindGateways() ([]entity.Gateway, error) {
	var gateways []entity.Gateway = make([]entity.Gateway, 0)

	rows, err := repo.db.Query(`
	SELECT id, key, name, acronym, icon, type, can_process, can_settle, created_at, updated_at
	FROM checkout.gateways;
	`)

	if err != nil {
		return gateways, err
	}

	defer rows.Close()

	for rows.Next() {
		var gateway entity.Gateway

		err = rows.Scan(
			&gateway.Id,
			&gateway.Key,
			&gateway.Name,
			&gateway.Acronym,
			&gateway.Icon,
			&gateway.Type,
			&gateway.CanProcess,
			&gateway.CanSettle,
			&gateway.CreatedAt,
			&gateway.UpdatedAt,
		)
		if err != nil {
			return gateways, err
		}

		gateways = append(gateways, gateway)
	}

	return gateways, nil
}

func (repo CheckoutPSQLRepo) FindGatewayByKey(key string) (*entity.Gateway, error) {
	var gateway entity.Gateway

	err := repo.db.QueryRow(`
	SELECT id, key, name, acronym, icon, type, can_process, can_settle, created_at, updated_at
	FROM checkout.gateways
	WHERE key = $1;
	`, key).Scan(
		&gateway.Id,
		&gateway.Key,
		&gateway.Name,
		&gateway.Acronym,
		&gateway.Icon,
		&gateway.Type,
		&gateway.CanProcess,
		&gateway.CanSettle,
		&gateway.CreatedAt,
		&gateway.UpdatedAt,
	)

	if err != nil {
		return &gateway, err
	}

	return &gateway, nil
}
