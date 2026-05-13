package repository

import (
	"database/sql"

	"github.com/kidus-yoseph1/ecommerce-api/internal/domain"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) AddProduct(product domain.Product) error {
	_, err := r.db.Exec(`
	INSERT INTO products (name, description, category, price, stock_quantity) 
	VALUES ($1, $2, $3, $4, $5)`,
		product.Name,
		product.Description,
		product.Category,
		product.Price,
		product.StockQuantity,
	)
	return err
}

func (r *ProductRepo) GetProductByID(id string) (*domain.Product, error) {
	product := &domain.Product{}

	err := r.db.QueryRow(`
		SELECT id, name,  description, category, price, stock_quantity, deleted_at, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(
		&product.Id,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.Price,
		&product.StockQuantity,
		&product.DeletedAt,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepo) ListProducts(category string) ([]domain.Product, error) {
	products := []domain.Product{}

	var rows *sql.Rows
	var err error

	if category != "" {
		rows, err = r.db.Query(`
			SELECT id, name, description, category, price, stock_quantity, deleted_at, created_at, updated_at
			FROM products
			WHERE category = $1 AND deleted_at IS NULL
		`, category)
	} else {
		rows, err = r.db.Query(`
			SELECT id, name, description, category, price, stock_quantity, deleted_at, created_at, updated_at
			FROM products
			WHERE deleted_at IS NULL
		`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Product
		rows.Scan(
			&p.Id,
			&p.Name,
			&p.Description,
			&p.Category,
			&p.Price,
			&p.StockQuantity,
			&p.DeletedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepo) UpdateProduct(product domain.Product) error {
	_, err := r.db.Exec(`
        UPDATE products 
        SET name=$1, description=$2, category=$3, price=$4, stock_quantity=$5, updated_at=now()
        WHERE id=$6
    `,
		product.Name,
		product.Description,
		product.Category,
		product.Price,
		product.StockQuantity,
		product.Id,
	)
	return err
}

func (r *ProductRepo) DeleteProduct(id string) error {
	_, err := r.db.Exec(`
        UPDATE products SET deleted_at=now() WHERE id=$1
    `, id)
	return err
}
