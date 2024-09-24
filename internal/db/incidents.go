package db

import (
	"context"
	"database/sql"
	"errors"
)

type IncidentRepository struct {
	db *sql.DB
}

func NewIncidentRepository(db *sql.DB) (IncidentRepository, error) {
	if db == nil {
		return IncidentRepository{}, errors.New("database is nil")
	}

	return IncidentRepository{db: db}, nil
}

func (repo *IncidentRepository) Create(ctx context.Context, incident Incident) (*Incident, error) {
	var res Incident
	err := repo.db.QueryRowContext(
		ctx,
		"INSERT INTO incidents (title, description) VALUES ($1,$2) RETURNING id, title, description, created_at",
		incident.Title,
		incident.Details,
	).Scan(&incident.Id, &incident.Title)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (repo *IncidentRepository) Update(ctx context.Context, incident Incident) error {
	res, err := repo.db.Exec(
		"UPDATE incidents SET title = $1, description = $2 WHERE id = $3",
		incident.Title,
		incident.Details,
		incident.Id,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("incident not found")
	}

	return nil
}

func (repo *IncidentRepository) GetById(ctx context.Context, id uint) (*Incident, error) {
	var incident Incident
	err := repo.db.QueryRowContext(
		ctx,
		"SELECT id, title, description, created_at FROM incidents WHERE id = $1",
		id,
	).Scan(&incident.Id, &incident.Title, &incident.Details, &incident.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("incident not found")
	} else if err != nil {
		return nil, err
	}

	return &incident, nil
}

func (repo *IncidentRepository) GetAll(ctx context.Context) ([]*Incident, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, title, description, created_at FROM incidents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*Incident

	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.Id, &incident.Title, &incident.Details, &incident.CreatedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, &incident)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}

func (repo *IncidentRepository) Delete(ctx context.Context, id uint) error {
	res, err := repo.db.Exec("DELETE FROM incidents WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("incident not found")
	}
	return nil
}

// find incidents for the specified system
func (repo *IncidentRepository) GetBySystem(ctx context.Context, systemId uint) ([]*Incident, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT i.id, i.title, i.description, i.created_at FROM system_incidents si INNER JOIN incidents i ON i.id = si.incident_id WHERE si.system_id = $1",
		systemId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []*Incident

	for rows.Next() {
		var incident Incident
		err := rows.Scan(&incident.Id, &incident.Title, &incident.Details, &incident.CreatedAt)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, &incident)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return incidents, nil
}
