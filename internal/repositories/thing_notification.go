package repositories

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const (
	thingNotificationTableName = "thing_notification"
)

type ThingNotificationRepository struct {
	db *sql.DB
}

func InitThingNotificationRepository(db *sql.DB) *ThingNotificationRepository {
	return &ThingNotificationRepository{db: db}
}

func (r ThingNotificationRepository) Get(ctx context.Context, thingID int) (*models.ThingNotification, error) {
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

func (r ThingNotificationRepository) GetExpired(ctx context.Context) ([]models.ExtThingNotification, error) {
	var res []models.ExtThingNotification

	query, args, err := sq.Select("n.thing_id", "n.notification_date", "n.created_at", "n.updated_at", "t.title", "p.id", "p.title").
		From(thingNotificationTableName + " n").
		Join(thingTableName + " t ON t.id = n.thing_id").
		LeftJoin(placeThingTableName + " pt ON pt.thing_id = n.thing_id").
		LeftJoin(placeTableName + " p ON p.id = pt.place_id").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Lt{"n.notification_date": "NOW()"}).
		OrderBy("n.notification_date DESC").
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resRow := models.ExtThingNotification{}

		err = rows.Scan(
			&resRow.ThingID,
			&resRow.NotificationDate,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
			&resRow.ThingTitle,
			&resRow.PlaceID,
			&resRow.PlaceTitle,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, resRow)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r ThingNotificationRepository) Add(ctx context.Context, req models.AddThingNotificationRequest, tx *sql.Tx) error {
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

func (r ThingNotificationRepository) Update(ctx context.Context, req models.UpdateThingNotificationRequest, tx *sql.Tx) error {
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

func (r ThingNotificationRepository) Delete(ctx context.Context, thingID int, tx *sql.Tx) error {
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
