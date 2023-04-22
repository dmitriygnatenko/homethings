package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.ThingNotificationRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	thingNotificationTableName = "thing_notification"
)

type thingNotificationRepository struct {
	db *sql.DB
}

func InitThingNotificationRepository(db *sql.DB) interfaces.ThingNotificationRepository {
	return thingNotificationRepository{db: db}
}

func (r thingNotificationRepository) Get(ctx context.Context, thingID int) (*models.ThingNotification, error) {
	query, args, err := sq.Select("thing_id", "notification_date", "created_at", "updated_at").
		From(thingNotificationTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": thingID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var res models.ThingNotification

	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(&res.ThingID, &res.NotificationDate, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r thingNotificationRepository) GetExpired(ctx context.Context) ([]models.ExtThingNotification, error) {
	return nil, nil // TODO
}

func (r thingNotificationRepository) Add(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(thingNotificationTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("thing_id", "notification_date").
		Values(req.ThingID, req.NotificationDate).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}

func (r thingNotificationRepository) Update(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) error {
	query, args, err := sq.Update(thingNotificationTableName).
		PlaceholderFormat(sq.Dollar).
		Set("notification_date", req.NotificationDate).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"thing_id": req.ThingID}).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}

func (r thingNotificationRepository) Delete(ctx context.Context, thingID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingNotificationTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": thingID}).
		ToSql()

	if err != nil {
		return err
	}

	if tx == nil {
		_, err = r.db.ExecContext(ctx, query, args...)
	} else {
		_, err = tx.ExecContext(ctx, query, args...)
	}

	return err
}
