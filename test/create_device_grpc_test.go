package test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	act_device_api "github.com/ozonmp/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
	"testing"
)

func TestCreateDevice(t *testing.T) {
	ctx := context.Background()
	host := "localhost:8082"
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc.Dial err:%v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Logf("conn.Close err:%v", err)
		}
	}(conn)

	t.Run("CD1 CreateDevice create device with correct data", func(t *testing.T) {

		platform := "Win10"
		userId := 111222333
		request := &act_device_api.CreateDeviceV1Request{ // передаем номер устройства
			Platform: platform,
			UserId:   uint64(userId),
		}

		createDevicesV1Response, err := act_device_api.NewActDeviceApiServiceClient(conn).CreateDeviceV1(ctx, request)
		assert.Equal(t, codes.OK.String(), status.Code(err).String()) //  проверяем статус код
		require.NoError(t, err)                                       // отсутствуют ошибок в ответе
		require.NotNil(t, createDevicesV1Response)                    // ответ не пустой
		// сравним полученный ответ с ожидаемым
		assert.NotNil(t, createDevicesV1Response.DeviceId) // не пусто
		t.Logf("DeviceId %v", createDevicesV1Response.DeviceId)
	})

	t.Run("CD2 CreateDevice create device without platform", func(t *testing.T) {

		platform := ""
		userId := 222333
		request := &act_device_api.CreateDeviceV1Request{ // передаем номер устройства
			Platform: platform,
			UserId:   uint64(userId),
		}

		_, err := act_device_api.NewActDeviceApiServiceClient(conn).CreateDeviceV1(ctx, request)
		assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String()) //  проверяем статус код

	})

	t.Run("CD3 CreateDevice create device without userId", func(t *testing.T) {

		platform := "Win10"
		var userId uint64
		request := &act_device_api.CreateDeviceV1Request{ // передаем номер устройства
			Platform: platform,
			UserId:   userId,
		}

		_, err := act_device_api.NewActDeviceApiServiceClient(conn).CreateDeviceV1(ctx, request)
		assert.Equal(t, codes.InvalidArgument.String(), status.Code(err).String()) //  проверяем статус код

	})
}
