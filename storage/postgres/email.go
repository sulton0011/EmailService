package postgres

import (
	"time"

	"github.com/home_work/TaskUserService/EmailService/storage/repo"
	"github.com/jmoiron/sqlx"
)

type sendRepo struct {
	db *sqlx.DB
}

// NewSendRepo ...
func NewSendRepo(db *sqlx.DB) repo.SendStorageI {
	return &sendRepo{db: db}
}

// MakeSent ...
func (s *sendRepo) CreatEmailText(id, subject, body string) (string, error) {
	query := `
		INSERT INTO email_text(id, subject, body, created_at)
		VALUES($1, $2, $3, $4)
		RETURNING id;
	`
	err := s.db.DB.QueryRow(
		query,
		id,
		subject,
		body,
		time.Now().UTC().Format("2006-01-02 15:04:05"),
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, err
}

func (s *sendRepo) CreatEmail(id, emailTextId, email string,  status bool) error {
	query := `
		INSERT INTO email_send_email(
			id,
			email_text_id,
			email,
			send_status,
			created_at
		)
		VALUES($1, $2, $3, $4, $5);
	`

	err := s.db.DB.QueryRow(
		query,
		id,
		emailTextId,
		email,
		status,
		time.Now().UTC().Format("2006-01-02 15:04:05"),
	).Err()
	if err != nil {
		return err
	}
	return nil
}
