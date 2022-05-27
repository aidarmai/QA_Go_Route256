package sql_client

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/act-device-api/internal/database"
	"github.com/ozonmp/act-device-api/test/internal/models"
)

type EventStorage interface {
}

type storage struct {
	DB *sqlx.DB
}

func NewPostgres(dsn string) (*storage, error) {
	db, err := database.NewPostgres(dsn, "pgx")
	if err != nil {
		return nil, err
	}
	return &storage{DB: db}, nil
}

func (r storage) ByDeviceId(ctx context.Context, deviceID int) (*models.DeviceEvent, error) {
	var (
		event models.DeviceEvent
	)
	query := sq.Select("id", "device_id", "type", "status", "payload", "created_at", "updated_at").
		PlaceholderFormat(sq.Dollar).
		From("devices_events").
		Where(sq.Eq{"device_id": deviceID}).
		OrderBy("id DESC").
		Limit(1)

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.DB.GetContext(ctx, &event, s, args...)

	return &event, err
}

func (r storage) EventIdCount() (string, error) {

	maxCountRow := r.DB.QueryRow("SELECT max(id) FROM devices_events")
	var maxCount string
	err := maxCountRow.Scan(&maxCount)
	if err != nil {
		return "", err
	}
	return maxCount, nil
}

func (r storage) ParsePayload(numDevice int, s string) (string, error) {
	payloadRow := r.DB.QueryRow("SELECT payload ->> $1 FROM devices_events WHERE device_id = $2 ORDER BY id DESC LIMIT 1",
		s, numDevice)
	var payload string
	err := payloadRow.Scan(&payload)
	if err != nil {
		return "", err
	}
	return payload, nil
}

func (r storage) PayloadIsNull(id int64) (bool, error) {
	payloadRow := r.DB.QueryRow("SELECT payload FROM devices_events WHERE device_id = $1 ORDER BY id DESC LIMIT 1", id)
	var payload string
	err := payloadRow.Scan(&payload)
	if err != nil {
		return false, err
	}
	if payload == "null" {
		return true, nil
	}
	return false, err
}
