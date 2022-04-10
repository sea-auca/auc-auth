package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sea-auca/auc-auth/user/service"
	"go.uber.org/zap"
)

type pgxUserRepository struct {
	db *pgxpool.Pool
}

func NewPgxUserRepo(db *pgxpool.Pool) service.UserRepository {
	return pgxUserRepository{db: db}
}

func (r pgxUserRepository) Create(ctx context.Context, u *service.User) (*service.User, error) {
	tx, err := r.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		zap.L().Error("PGX failed to start transaction", zap.Error(err))
		return nil, err
	}

	const query_create_user = `insert into users.users (id, email) VALUES ($1,$2) RETURNING created_at`
	err = tx.QueryRow(ctx, query_create_user, u.ID, u.Email).Scan(u.CreatedAt)
	if err != nil {
		zap.L().Error("PGX Query failed", zap.Error(err))
		return u, err
	}
	u.UpdatedAt = u.CreatedAt

	const query_create_settings = `insert into users.user_data (user_id) VALUES($1); insert into users.authentication_settings (user_id) VALUES ($1);`
	_, err = tx.Exec(ctx, query_create_settings, u.ID)
	if err != nil {
		zap.L().Error("PGX Query failed", zap.Error(err))
		return u, err
	}

	tx.Commit(ctx)
	return nil, nil
}

func (r pgxUserRepository) Update(ctx context.Context, u *service.User) error {
	const query = `update users.users set is_active=$1, is_validated=$2, updated_at=current_timestamp where id=$3`
	_, err := r.db.Exec(ctx, query, u.IsActive, u.IsValidated, u.ID)
	return err
}

func (r pgxUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*service.User, error) {
	const query = `select id, email, is_active, is_validated, created_at, updated_at from users.users where id=$1`
	var u service.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.IsActive, &u.IsValidated, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r pgxUserRepository) GetByEmail(ctx context.Context, email string) (*service.User, error) {
	const query = `select id, email, is_active, is_validated, created_at, updated_at from users.users where email=$1`
	var u service.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.IsActive, &u.IsValidated, &u.CreatedAt, &u.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r pgxUserRepository) PaginatedView(ctx context.Context, page, pageSize int) ([]*service.User, int, error) {
	const query = `select id, email, is_active, is_validated, created_at, updated_at from users.users offset $1 limit $2`
	var us []*service.User
	rows, err := r.db.Query(ctx, query, (page-1)*pageSize, pageSize)
	if err != nil {
		zap.L().Error("PGX Query failed", zap.Error(err))
		return nil, 0, err
	}
	for rows.Next() {
		var u service.User
		err := rows.Scan(&u.ID, &u.Email, &u.IsActive, &u.IsValidated, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			zap.L().Error("PGX Query failed", zap.Error(err))
			return nil, 0, err
		}
		us = append(us, &u)
	}
	if err := rows.Err(); err != nil {
		zap.L().Error("PGX Query failed", zap.Error(err))
		return nil, 0, err
	}
	return us, 0, nil
}

type pgxVerificationLinkRepo struct {
	db *pgxpool.Pool
}

func NewPgxVerificationLinkRepo(db *pgxpool.Pool) service.VerificationRepository {
	return pgxVerificationLinkRepo{db: db}
}

func (r pgxVerificationLinkRepo) Create(ctx context.Context, vl *service.VerificationLink) (*service.VerificationLink, error) {
	const query = `insert into users.validation_requests (id, user_id, expires_at) VALUES($1,$2,$3) RETURNING created_at`
	err := r.db.QueryRow(ctx, query, vl.ID, vl.UserID, vl.ExpiresAt).Scan(vl.CreatedAt)
	if err != nil {
		zap.L().Error("PGX Query failed", zap.Error(err))
		return vl, err
	}
	return vl, nil
}

func (r pgxVerificationLinkRepo) SearchByID(ctx context.Context, id uuid.UUID) (*service.VerificationLink, error) {
	const query = `select id, user_id, was_utilised, expires_at, created_at, updated_at from users.validation_requests where id=$1`
	var l service.VerificationLink
	err := r.db.QueryRow(ctx, query, id).Scan(&l.ID, &l.UserID, &l.WasUtilised, &l.ExpiresAt, &l.CreatedAt, &l.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r pgxVerificationLinkRepo) DeactivateLink(ctx context.Context, id uuid.UUID) error {
	const query = `update users.validation_requests set was_utilised='true' where id=$1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
