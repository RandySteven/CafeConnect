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
		queries.CreateCafeFranchiseTable,
		queries.CreateCafeTable,
		queries.CreateProductCategoryTable,
		queries.CreateProductTable,
		queries.CreateCafeProductTable,
		queries.CreateReviewTable,
		queries.CreateCartTable,
		queries.CreateTransactionHeaderTable,
		queries.CreateTransactionDetailTable,
		queries.CreateMidtransTransactionTable,
		queries.CreateVerifyTokenTable,
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
