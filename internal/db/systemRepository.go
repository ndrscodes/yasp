package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/georgysavva/scany/v2/sqlscan"
)

var ctx = context.Background()

type SystemRepository struct {
	db *sql.DB
}

func NewSystemRepository(db *sql.DB) (SystemRepository, error) {
	if db == nil {
		return SystemRepository{}, errors.New("database is nil")
	}

	return SystemRepository{db: db}, nil
}

func (repo *SystemRepository) Create(sys System) (System, error) {
	var system System
	err := sqlscan.Get(ctx, repo.db, &system, "INSERT INTO systems (name) VALUES ($1) RETURNING id", sys.Name)
	return system, err
}

func (repo *SystemRepository) Update(sys System) error {
	res, err := repo.db.Exec("UPDATE systems SET name = $1 WHERE id = $2", sys.Name, sys.Id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		return errors.New("system not found")
	}
	return nil
}

func (repo *SystemRepository) GetById(id uint) (*System, error) {
	var system System
	err := sqlscan.Get(ctx, repo.db, &system, "SELECT id, name FROM systems WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &system, nil
}

func (repo *SystemRepository) GetAll() ([]*System, error) {
	var systems []*System
	err := sqlscan.Select(ctx, repo.db, &systems, "SELECT * FROM systems")
	if err != nil {
		return nil, err
	}
	return systems, nil
}

func (repo *SystemRepository) Delete(id uint) (System, error) {
	var system System
	err := sqlscan.Get(ctx, repo.db, &system, "DELETE FROM systems WHERE id = $1 RETURNING id, name", id)
	if err != nil {
		return System{}, err
	}
	return system, nil
}
