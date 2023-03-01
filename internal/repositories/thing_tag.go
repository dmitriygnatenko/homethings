package repositories

//go:generate mkdir -p mocks
//go:generate rm -rf ./mocks/*_minimock.go
//go:generate minimock -i git.dmitriygnatenko.ru/dima/homethings/internal/interfaces.ThingTagRepository -o ./mocks/ -s "_minimock.go"

import (
	"context"
	"database/sql"
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/models"
	sq "github.com/Masterminds/squirrel"
)

const (
	thingTagTableName = "thing_tag"
)

type thingTagRepository struct {
	db *sql.DB
}

func InitThingTagRepository(db *sql.DB) interfaces.ThingTagRepository {
	return thingTagRepository{db: db}
}

func (r thingTagRepository) BeginTx(ctx context.Context, level sql.IsolationLevel) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
}

func (r thingTagRepository) CommitTx(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("empty transaction")
	}

	return tx.Commit()
}

func (r thingTagRepository) Add(ctx context.Context, req models.AddThingTagRequest, tx *sql.Tx) error {
	query, args, err := sq.Insert(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Columns("thing_id", "tag_id").
		Values(req.ThingID, req.TagID).
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

func (r thingTagRepository) Delete(ctx context.Context, req models.DeleteThingTagRequest, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"thing_id": req.ThingID}).
		Where(sq.Eq{"tag_id": req.TagID}).
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

func (r thingTagRepository) DeleteByTagID(ctx context.Context, tagID int, tx *sql.Tx) error {
	query, args, err := sq.Delete(thingTagTableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"tag_id": tagID}).
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

func (r thingTagRepository) GetByPlaceID(ctx context.Context, placeID int) ([]models.ThingTag, error) {
	var res []models.ThingTag

	query := "WITH RECURSIVE cte (id, parent_id) AS (" +
		"SELECT id, parent_id " +
		"FROM " + placeTableName + " " +
		"WHERE id = $1 " +
		"UNION ALL " +
		"SELECT p.id, p.parent_id " +
		"FROM " + placeTableName + " p " +
		"INNER JOIN cte ON p.parent_id = cte.id " +
		")" +
		"SELECT t.id, t.title, t.style, t.created_at, t.updated_at, tt.thing_id " +
		"FROM cte, " + placeThingTableName + " pt, " + thingTagTableName + " tt, " + tagTableName + " t " +
		"WHERE pt.place_id = cte.id AND tt.thing_id = pt.thing_id AND tt.tag_id = t.id " +
		"ORDER BY t.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, placeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resRow := models.ThingTag{}

		err = rows.Scan(
			&resRow.ID,
			&resRow.Title,
			&resRow.Style,
			&resRow.CreatedAt,
			&resRow.UpdatedAt,
			&resRow.ThingID,
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
