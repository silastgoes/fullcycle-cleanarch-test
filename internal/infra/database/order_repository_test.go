package database

import (
	"database/sql"
	"testing"

	"github.com/silastgoes/fullcycle-cleanarch-test/internal/entity"
	"github.com/stretchr/testify/suite"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())

	order1, err := entity.NewOrder("1234", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)

	err = repo.Save(order)
	suite.NoError(err)

	err = repo.Save(order1)
	suite.NoError(err)

	list, err := repo.List()
	suite.NoError(err)

	suite.Equal(order.ID, list[0].ID)
	suite.Equal(order.Price, list[0].Price)
	suite.Equal(order.Tax, list[0].Tax)
	suite.Equal(order.FinalPrice, list[0].FinalPrice)

	suite.Equal(order1.ID, list[1].ID)
	suite.Equal(order1.Price, list[1].Price)
	suite.Equal(order1.Tax, list[1].Tax)
	suite.Equal(order1.FinalPrice, list[1].FinalPrice)
}
