package mysql_client

import (
	"context"
	"github.com/RandySteven/CafeConnect/be/queries"
)

func initMigrations() []queries.MigrationQuery {
	return []queries.MigrationQuery{
		queries.CreateUserTable,
		queries.CreatePointTable,
		queries.CreateRoleTable,
		queries.CreateRoleUserTable,
		queries.CreateAddressTable,
		queries.CreateAddressUserTable,
		queries.CreateReferralTable,
	}
}

func (m *mysqlClient) Migration(ctx context.Context) error {
	migrations := initMigrations()

	for _, query := range migrations {
		_, err := m.db.ExecContext(ctx, query.String())
		if err != nil {
			return err
		}
	}
	return nil
}
