package database

import (
	"database/sql"

	"github.com/silastgoes/fullcycle-cleanarch-test/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List() ([]*entity.Order, error) {
	list := make([]*entity.Order, 0)

	rows, err := r.Db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var i entity.Order
		if err := rows.Scan(&i.ID, &i.Price, &i.Tax, &i.FinalPrice); err != nil {
			return nil, err
		}
		list = append(list, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	err := r.Db.QueryRow("Select count(*) from orders").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
