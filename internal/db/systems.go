package db

import (
	"context"
	"database/sql"
	"errors"
)

type SystemRepository struct {
	db *sql.DB
}

func NewSystemRepository(db *sql.DB) (SystemRepository, error) {
	if db == nil {
		return SystemRepository{}, errors.New("database is nil")
	}

	return SystemRepository{db: db}, nil
}

func (repo *SystemRepository) Create(ctx context.Context, system System) (*System, error) {
	var res System
	var status string
	err := repo.db.QueryRowContext(ctx, "INSERT INTO systems (name) VALUES ($1) RETURNING id, name, status", system.Name).Scan(&system.Id, &system.Name, &status)

	return &res, err
}

func (repo *SystemRepository) Update(sys System) error {
	res, err := repo.db.Exec("UPDATE systems SET name = $1, current_status = $2 WHERE id = $3", sys.Name, sys.Status, sys.Id)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil || rows == 0 {
		return errors.New("system not found")
	}
	return nil
}

func (repo *SystemRepository) GetById(ctx context.Context, id uint) (*System, error) {
	var system System
	var status string
	err := repo.db.QueryRowContext(ctx, "SELECT id, name, current_status FROM systems WHERE id = $1", id).Scan(&system.Id, &system.Name, &status)
	if err != nil {
		return nil, err
	}

	if system.Status, err = ToStatus(status); err != nil {
		return nil, err
	}

	return &system, nil
}

func (repo *SystemRepository) GetAll(ctx context.Context) ([]*System, error) {
	var systems []*System
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM systems")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s System
		var status string
		err := rows.Scan(&s.Id, &s.Name, &status)
		if err != nil {
			return nil, err
		}

		if s.Status, err = ToStatus(status); err != nil {
			return nil, err
		}

		systems = append(systems, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return systems, nil
}

func (repo *SystemRepository) Delete(ctx context.Context, id uint) (System, error) {
	var system System
	var status string
	err := repo.db.QueryRowContext(ctx, "DELETE FROM systems WHERE id = $1 RETURNING id, name, status", id).Scan(&system.Id, &system.Name, &status)
	if err != nil {
		return System{}, err
	}

	if system.Status, err = ToStatus(status); err != nil {
		return System{}, err
	}

	return system, nil
}

// find all systems affected by the specified incident
func (repo *SystemRepository) GetAffectedByIncident(ctx context.Context, incident uint) ([]*System, error) {
	var systems []*System
	rows, err := repo.db.QueryContext(ctx, "SELECT DISTINCT s.id, s.name, s.status FROM systems s INNER JOIN system_incidents si ON si.system_id = s.id AND si.incident_id = $1", incident)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s System
		var stat string
		err := rows.Scan(&s.Id, &s.Name, &stat)
		if err != nil {
			return nil, err
		}
		if s.Status, err = ToStatus(stat); err != nil {
			return nil, err
		}
		systems = append(systems, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return systems, nil
}
