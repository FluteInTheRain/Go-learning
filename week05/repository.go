package main

import (
	"database/sql"
	"fmt"
)

type ProductRepo struct {
	db *sql.DB
}

// Create inserts a new product and returns its generated id.
func (r *ProductRepo) Create(p *Product) error {
	// RETURNING id lets Postgres give back the new auto-generated id.
	query := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, p.Name, p.Price).Scan(&p.ID)
}

// GetAll returns every product in the table.
func (r *ProductRepo) GetAll(limit, offset int) ([]Product, error) {
	rows, err := r.db.Query(`SELECT id, name, price FROM products ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // ALWAYS close rows to avoid connection leaks

	var list []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// GetByID returns one product, or sql.ErrNoRows if not found.
func (r *ProductRepo) GetByID(id int64) (*Product, error) {
	var p Product
	query := `SELECT id, name, price FROM products WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err // could be sql.ErrNoRows
	}
	return &p, nil
}

// Update modifies an existing product; returns rows affected.
func (r *ProductRepo) Update(p *Product) (int64, error) {
	query := `UPDATE products SET name = $1, price = $2 WHERE id = $3`
	res, err := r.db.Exec(query, p.Name, p.Price, p.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Delete removes a product by id; returns rows affected.
func (r *ProductRepo) Delete(id int64) (int64, error) {
	res, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// Get total of products
func (r *ProductRepo) Count() (int64, error) {
	query := `SELECT COUNT(*) FROM products`
	var count int64
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepo) PartialUpdate(id int64, up *UpdateProduct) (int64, error) {
	query := `UPDATE products SET `
	args := []any{}
	paramNum := 1
	if up.Name != nil {
		query += fmt.Sprintf("name = $%d, ", paramNum)
		args = append(args, *up.Name)
		paramNum++
	}
	if up.Price != nil {
		query += fmt.Sprintf("price = $%d, ", paramNum)
		args = append(args, *up.Price)
		paramNum++
	}
	if len(args) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", paramNum)
	res, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
