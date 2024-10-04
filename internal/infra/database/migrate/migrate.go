package migrate

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrateService struct {
	pathMigration string
	DB            *sql.DB
}

type MigrateInterface interface {
	Up() error
}

func NewMigrateService(db *sql.DB) *migrateService {
	return &migrateService{
		pathMigration: "internal/infra/database/migrate/migrations",
		DB:            db,
	}
}

func (ms *migrateService) Up() error {
	driver, err := mysql.WithInstance(ms.DB, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("error gerar drive: %s", err.Error())
	}

	// Criar uma nova instância de migração
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", ms.pathMigration), // Caminho para as migrações
		"mysql", // Nome do banco de dados
		driver,
	)
	if err != nil {
		return fmt.Errorf("error ao ler migration: %s", err.Error())
	}

	// Aplicar as migrações
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error ao rodar migration: %s", err.Error())
	}
	return nil
}
