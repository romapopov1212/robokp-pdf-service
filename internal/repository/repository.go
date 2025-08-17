package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type PdfRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (*PdfRepository, error) {
	repo := &PdfRepository{db: db}
	if err := repo.CreateTable(); err != nil {
		return nil, fmt.Errorf("failed to create table")
	}
	return repo, nil
}

func (p *PdfRepository) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS pdf_kp(
		id BIGSERIAL PRIMARY KEY,
		id_user BIGINT,
		id_cart BIGINT,
		id_publication BIGINT,
		publication_url TEXT,
		logo JSONB,
		executor_parameters JSONB,
		presentation_parameters JSONB,
		style_template JSONB,
		count INT,
		save_required BOOLEAN,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	    );
	`
	_, err := p.db.Exec(query)
	return err
}

func (p *PdfRepository) Save(
	ctx context.Context,
	userId int64,
	cartId int64,
	publicationId int64,
	logo json.RawMessage,
	executorParameters json.RawMessage,
	presentationParameters json.RawMessage,
	styleTemplate json.RawMessage,
	count int,
) error {
	const op = "repository.Save"
	query := `
	INSERT INTO pdf_kp (
	    id_user,
		id_cart,
		id_publication,
		logo,
		executor_parameters,
		presentation_parameters,
		style_template,
		count,
		save_required,
		created_at,
		updated_at
	) VALUES (
	    $1, $2, $3, $4, $5, $6, $7, $8, true, now(), now()
	)
	`
	_, err := p.db.ExecContext(ctx, query, userId, cartId, publicationId, logo, executorParameters, presentationParameters, styleTemplate, count)
	if err != nil {
		return fmt.Errorf("ошибка добавления записи: %s: %v", op, err)
	}
	
	return nil
}
