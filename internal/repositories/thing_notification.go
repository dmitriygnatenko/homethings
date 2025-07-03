package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
)

const thingNotificationTableName = "thing_notification"

type ThingNotificationRepository struct {
	db DB
}

func InitThingNotificationRepository(db DB) *ThingNotificationRepository {
	return &ThingNotificationRepository{db: db}
}

func (r ThingNotificationRepository) Get(ctx context.Context, id uint64) (*models.ThingNotification, error) {
	q, v, err := sq.Select("thing_id", "notification_date", "created_at", "updated_at").
		From(thingNotificationTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res models.ThingNotification

	err = r.db.GetContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	return &res, nil
}

func (r ThingNotificationRepository) GetExpired(ctx context.Context) ([]models.ExtThingNotification, error) {
	var res []models.ExtThingNotification

	q, v, err := sq.Select(
		"n.thing_id AS thing_id",
		"n.notification_date",
		"n.created_at",
		"n.updated_at",
		"t.title AS thing_title",
		"p.id AS place_id",
		"p.title AS place_title",
	).
		From(thingNotificationTableName + " n").
		Join(thingTableName + " t ON t.id = n.thing_id").
		LeftJoin(placeThingTableName + " pt ON pt.thing_id = n.thing_id").
		LeftJoin(placeTableName + " p ON p.id = pt.place_id").
		PlaceholderFormat(placeholder).
		Where(sq.Lt{"n.notification_date": time.Now()}).
		OrderBy("n.notification_date DESC").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	err = r.db.SelectContext(ctx, &res, q, v...)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}

	return res, nil
}

func (r ThingNotificationRepository) Add(ctx context.Context, req models.AddThingNotificationRequest) error {
	q, v, err := sq.Insert(thingNotificationTableName).
		PlaceholderFormat(placeholder).
		Columns("thing_id", "notification_date").
		Values(req.ThingID, req.NotificationDate).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r ThingNotificationRepository) Update(ctx context.Context, req models.UpdateThingNotificationRequest) error {
	q, v, err := sq.Update(thingNotificationTableName).
		PlaceholderFormat(placeholder).
		Set("notification_date", req.NotificationDate).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"thing_id": req.ThingID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}

func (r ThingNotificationRepository) Delete(ctx context.Context, id uint64) error {
	q, v, err := sq.Delete(thingNotificationTableName).
		PlaceholderFormat(placeholder).
		Where(sq.Eq{"thing_id": id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, q, v...)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
