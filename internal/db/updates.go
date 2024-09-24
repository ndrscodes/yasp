package db

import (
	"context"
	"database/sql"
	"errors"
)

type UpdateRepository struct {
	db *sql.DB
}

func NewUpdateRepository(db *sql.DB) (UpdateRepository, error) {
	if db == nil {
		return UpdateRepository{}, errors.New("database is nil")
	}

	return UpdateRepository{db: db}, nil
}

func (repo *UpdateRepository) Create(ctx context.Context, incident uint, update Update) (*Update, error) {
	var res Update
	dbs, err := FromStatus(update.Status)
	if err != nil {
		return nil, err
	}

	var status string
	err = repo.db.QueryRowContext(
		ctx,
		"INSERT INTO updates (title, details, incident_id, new_status) VALUES ($1, $2, $3, $4) RETURNING id, title, description, created_at, new_status",
		update.Title,
		update.Details,
		incident,
		dbs,
	).Scan(&update.Id, &update.Details, &update.CreatedAt, &status)

	if err != nil {
		return nil, err
	}

	if res.Status, err = ToStatus(status); err != nil {
		return nil, err
	}

	return &res, nil
}

func (repo *UpdateRepository) Update(ctx context.Context, update Update) error {
	dbs, err := FromStatus(update.Status)
	if err != nil {
		return err
	}

	res, err := repo.db.Exec(
		"UPDATE updates SET title = $1, details = $2, new_status = $3 WHERE id = $4",
		update.Title,
		update.Details,
		dbs,
		update.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("update not found")
	}

	return nil
}

func (repo *UpdateRepository) GetById(ctx context.Context, id uint) (*Update, error) {
	var update Update
	var status string
	err := repo.db.QueryRowContext(
		ctx,
		"SELECT id, title, description, created_at, new_status FROM updates WHERE id = $1",
		id,
	).Scan(&update.Id, &update.Title, &update.Details, &update.CreatedAt, &status)

	if err == sql.ErrNoRows {
		return nil, errors.New("update not found")
	} else if err != nil {
		return nil, err
	}

	if update.Status, err = ToStatus(status); err != nil {
		return nil, err
	}

	return &update, nil
}

func (repo *UpdateRepository) GetAll(ctx context.Context) ([]*Update, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, title, description, created_at, new_status FROM updates",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []*Update

	for rows.Next() {
		var update Update
		var status string
		err := rows.Scan(&update.Id, &update.Title, &update.Details, &update.CreatedAt, &status)
		if err != nil {
			return nil, err
		}

		if update.Status, err = ToStatus(status); err != nil {
			return nil, err
		}

		updates = append(updates, &update)
	}

	return updates, nil
}

func (repo *UpdateRepository) Delete(ctx context.Context, id uint) error {
	res, err := repo.db.Exec(
		"DELETE FROM updates WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("update not found")
	}

	return nil
}

// find all updates for the specified incident
func (repo *UpdateRepository) GetByIncident(ctx context.Context, incident uint) ([]*Update, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, title, description, created_at, new_status FROM updates WHERE incident_id = $1",
		incident,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []*Update

	for rows.Next() {
		var update Update
		var status string
		err := rows.Scan(&update.Id, &update.Title, &update.Details, &update.CreatedAt, &status)
		if err != nil {
			return nil, err
		}

		if update.Status, err = ToStatus(status); err != nil {
			return nil, err
		}

		updates = append(updates, &update)
	}

	return updates, nil
}
