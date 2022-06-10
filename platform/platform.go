package platform

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ozonmp/act-device-api/internal/database"
)

type DeviceModel interface {
	CountPlatform() (int64, error)
	CountPlatformCertainType(string) (int64, error)
}

type DeviceDB struct {
	*sqlx.DB
}

func NewPostgres(dsn string) (*DeviceDB, error) {
	db, err := database.NewPostgres(dsn, "pgx")
	if err != nil {
		return nil, err
	}
	return &DeviceDB{db}, nil
}

// CountPlatform метод выводит кол-во всех устройств
func (sdb *DeviceDB) CountPlatform() (int64, error) {
	var count int64
	err := sdb.QueryRow("SELECT count(*) FROM devices WHERE remove = 'false'").Scan(&count)
	return count, err
}

// CountPlatformCertainType метод выводит кол-во устройств определенной платформы
func (sdb *DeviceDB) CountPlatformCertainType(platform string) (int64, error) {
	var count int64
	err := sdb.QueryRow("SELECT count(*) FROM devices WHERE platform = $1 AND remove = 'false'", platform).Scan(&count)
	return count, err
}

// PercentagePlatformCertainType  функция возвращает кол-во устройств определенной платформы
func PercentagePlatformCertainType(sdb DeviceModel) (string, error) {
	allCount, err := sdb.CountPlatform()
	if err != nil {
		return "", err
	}

	countCertainType, err := sdb.CountPlatformCertainType("ios")
	if err != nil {
		return "", err
	}

	rate := float64(countCertainType * 100 / allCount)
	return fmt.Sprintf("%.1f %s", rate, "%"), nil
}
