package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"url-shortener/models"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) FindByOgURL(ctx context.Context, OriginalURL string) (*models.URL, error) {
	query := `
			SELECT id, original_url, short_code, created_at , expires_at
			FROM urls
			WHERE original_url = $1;
	`

	row := r.db.QueryRowContext(ctx, query, OriginalURL)

	var ExistingUrl models.URL

	err := row.Scan(
		&ExistingUrl.ID,
		&ExistingUrl.OriginalURL,
		&ExistingUrl.ShortCode,
		&ExistingUrl.CreatedAt,
		&ExistingUrl.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err

	}
	return &ExistingUrl, nil
}

func (r *URLRepository) CreateURL(ctx context.Context, urlRecord *models.URL) (*models.URL, error) {

	query := `
		INSERT INTO urls (original_url , short_code , expires_at)
		VALUES ($1 , $2 , $3)
		ON CONFLICT (original_url)
		DO NOTHING 
		RETURNING id,original_url , short_code, created_at, expires_at;
	`

	var createdURL models.URL

	err := r.db.QueryRowContext(ctx, query, urlRecord.OriginalURL, urlRecord.ShortCode, urlRecord.ExpiresAt).Scan(
		&createdURL.ID,
		&createdURL.OriginalURL,
		&createdURL.ShortCode,
		&createdURL.CreatedAt,
		&createdURL.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r.FindByOgURL(ctx, urlRecord.OriginalURL)
		}
		return nil, err

	}

	return &createdURL, nil

}

func (r *URLRepository) RenewUrl(ctx context.Context, urlRecord *models.URL) (*models.URL, error) {
	query := `
		UPDATE urls
		SET short_code = $1 , expires_at = $2
		WHERE id = $3
			AND expires_at <= NOW()
		RETURNING id, original_url , short_code , created_at , expires_at;
	`
	var RenewedURL models.URL

	err := r.db.QueryRowContext(ctx, query, urlRecord.ShortCode, urlRecord.ExpiresAt, urlRecord.ID).Scan(
		&RenewedURL.ID,
		&RenewedURL.OriginalURL,
		&RenewedURL.ShortCode,
		&RenewedURL.CreatedAt,
		&RenewedURL.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err , sql.ErrNoRows){
			return r.FindByOgURL(ctx , urlRecord.OriginalURL)
		}
		return nil, err 
	}

	return &RenewedURL, nil

}

func (r *URLRepository) GetURLByShortCode(ctx context.Context, shortCode string) (string,time.Time, error)  {
	query := `
		SELECT  original_url,expires_at
		FROM urls
		WHERE short_code = $1;
	`
	
	var originalURL string
	var expiresAt time.Time

	err := r.db.QueryRowContext(ctx, query, shortCode).Scan(&originalURL, &expiresAt)

	if err != nil{
		if errors.Is(err , sql.ErrNoRows){
			return "", time.Time{} , nil
		}
		return "", time.Time{} ,err
	}
	return originalURL , expiresAt , nil


}