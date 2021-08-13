package tests

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDB returns a *gorm.DB instance point to a postgres or a sqlite
// depends on the envvar `POSTGRES_HOST`
// it's usefull for unittest in both, local and CI
func NewDB() (*gorm.DB, error) {
	hostName := os.Getenv("POSTGRES_HOST")
	if hostName != "" {
		return newPostgres(hostName)
	}
	return gorm.Open(sqlite.Open("fair.db"), &gorm.Config{})
}

func newPostgres(hostName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=fair password=fair dbname=fair sslmode=disable connect_timeout=3",
		hostName,
	)
	return gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)
}
