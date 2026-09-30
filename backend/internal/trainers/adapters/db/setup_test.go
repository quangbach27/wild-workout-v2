//go:build integration

package db_test

import (
	"embed"
	"os"
	"testing"

	"github.com/quangbach27/golang-common/testutils"

	"backend/internal/configs"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var (
	config = configs.NewConfig()
	pgxDB  = testutils.NewDB(config.DB.Dsn())
)

func TestMain(m *testing.M) {
	testutils.RunMigrations("trainers", embedMigrations, "migrations", config.DB.Dsn())

	code := m.Run()

	pgxDB.Close()

	os.Exit(code)
}
