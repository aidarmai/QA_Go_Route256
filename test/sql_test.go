package test

import (
	"context"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ozonmp/act-device-api/internal/config"
	"github.com/ozonmp/act-device-api/test/internal/steps"
	routeclient "github.com/ozonmp/act-device-api/test/route_client"
	"github.com/ozonmp/act-device-api/test/sql_client"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func Test_Sql(t *testing.T) {
	// читаем данные из yml конфига
	if config.ReadConfigYML("../config.yml") != nil {
		t.Fatal("Failed init configuration")
	}
	cfg := config.GetConfigInstance()
	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode,
	)

	// получаем данные от самописного sql клиента
	db, err := sql_client.NewPostgres(dsn)
	if err != nil {
		t.Fatal(err)
	}

	// запускаем тесты
	t.Run("Method CREATE writes logs in device_events", func(t *testing.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		device := routeclient.CreateDeviceRequest{
			Platform: "Ubuntu",
			UserId:   "5",
		}
		oldMaxId, _ := db.EventIdCount()
		ctx := context.TODO()
		actualTime := time.Now().UTC().Format(time.UnixDate)
		deviceFromSwagger, _, _ := client.CreateDevice(ctx, device)
		deviceFromSql, err := db.ByDeviceId(ctx, deviceFromSwagger.DeviceId)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equalf(t, strconv.Itoa(deviceFromSwagger.DeviceId), strconv.FormatUint(deviceFromSql.DeviceId, 10),
			"не совпадают значения device_id из swagger и sql")

		newMaxId := deviceFromSql.ID - 1
		assert.Equalf(t, strconv.FormatUint(newMaxId, 10), oldMaxId, "id новой записи в device_events неверный")

		assert.Equalf(t, actualTime, deviceFromSql.CreatedAt.Format(time.UnixDate),
			"не совпадает время вызова метода CREATE и время создания записи в device_events")

	})

	runner.Run(t, "Method UPDATE writes logs in device_events", func(t provider.T) {
		tt := []struct {
			numDevice int
			platform  string
			userId    string
		}{
			{23, "Ios", "200"},
		}
		oldMaxId, _ := db.EventIdCount()

		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		urlDevice, b, err := steps.UpdateDevice(t, tt[0].numDevice, tt[0].platform, tt[0].userId)
		if err != nil {
			t.Log(err)
		}
		req, err := http.NewRequest(http.MethodPut, urlDevice, b)
		if err != nil {
			panic(err)
		}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				t.Log(err)
			}
		}(res.Body)

		ctx := context.TODO()
		actualTime := time.Now().UTC().Format(time.UnixDate)

		deviceFromSql, err := db.ByDeviceId(ctx, tt[0].numDevice)
		if err != nil {
			t.Log(err)
		}

		id, _ := db.ParsePayload(tt[0].numDevice, "id")
		userId, _ := db.ParsePayload(tt[0].numDevice, "user_id")
		pl, _ := db.ParsePayload(tt[0].numDevice, "platform")

		assert.Equalf(t, tt[0].platform, pl, "")
		t.Assert().GreaterOrEqual(tt[0].platform, pl)
		t.Assert().GreaterOrEqual(tt[0].userId, userId)
		t.Assert().GreaterOrEqual(strconv.Itoa(tt[0].numDevice), id)

		newMaxId := deviceFromSql.ID - 1
		assert.Equalf(t, strconv.FormatUint(newMaxId, 10), oldMaxId, "id новой записи в device_events неверный")

		assert.Equalf(t, actualTime, deviceFromSql.CreatedAt.Format(time.UnixDate),
			"не совпадает время вызова метода UPDATE и время создания записи в device_events")

	})

	runner.Run(t, "Method DELETE writes logs in device_events", func(t provider.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		numDevice := 23
		apiUrl, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices/%d", numDevice))
		oldMaxId, _ := db.EventIdCount()

		req, err := http.NewRequest(http.MethodDelete, apiUrl.String(), nil)
		if err != nil {
			panic(err)
		}
		res, err := client.Do(req)
		if err != nil {
			t.Log(err)
		}
		_, _ = ioutil.ReadAll(res.Body)

		ctx := context.TODO()
		actualTime := time.Now().UTC().Format(time.UnixDate)
		deviceFromSql, err := db.ByDeviceId(ctx, numDevice)

		checkNul, err := db.PayloadIsNull(int64(numDevice))
		if err != nil {
			t.Log(err)
		}
		assert.Equalf(t, true, checkNul, "запись не удалена")

		newMaxId := deviceFromSql.ID - 1
		assert.Equalf(t, strconv.FormatUint(newMaxId, 10), oldMaxId, "id новой записи в device_events неверный")

		assert.Equalf(t, actualTime, deviceFromSql.CreatedAt.Format(time.UnixDate),
			"не совпадает время вызова метода DELETE и время создания записи в device_events")
	})
}
