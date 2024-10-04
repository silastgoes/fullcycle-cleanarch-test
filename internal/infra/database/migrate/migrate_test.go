package migrate

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrate_UP(t *testing.T) {
	assert := assert.New(t)
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/orders")
	assert.NoError(err)
	m := NewMigrateService(db)
	m.pathMigration = "./migrate/migrations"
	assert.NoError(m.Up())
}
