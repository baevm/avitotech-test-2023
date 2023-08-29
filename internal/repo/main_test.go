package repo

import (
	"os"
	"testing"

	"github.com/dezzerlol/avitotech-test-2023/pkg/testhelper"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testDbInstance *pgxpool.Pool

func TestMain(m *testing.M) {
	testDB := testhelper.SetupTestDatabase()

	testDbInstance = testDB.DbInstance
	defer testDB.TearDown()

	os.Exit(m.Run())
}
