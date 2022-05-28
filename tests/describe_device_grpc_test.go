package tests

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

func TestDescribeDevice(t *testing.T) {
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

	t.Run("DD1 DescribeDevice  requesting data about an existing device", func(t *testing.T) {

		deviceId := int64(10)
		request := &act_device_api.DescribeDeviceV1Request{ // передаем запрос номер устройства
			DeviceId: uint64(deviceId),
		}

		describeDevicesV1Response, err := act_device_api.NewActDeviceApiServiceClient(conn).DescribeDeviceV1(ctx, request)
		assert.Equal(t, codes.OK.String(), status.Code(err).String()) //  проверяем статус код
		require.NoError(t, err)                                       // отсутствуют ошибки в ответе
		require.NotNil(t, describeDevicesV1Response)                  // ответ не пустой
		// сравним полученный ответ с ожидаемым
		assert.Equal(t, int64(describeDevicesV1Response.Value.Id), deviceId) // получили запрашиваемый идентификатор
		assert.NotNil(t, describeDevicesV1Response.Value.Platform)           // не пусто
		assert.NotNil(t, describeDevicesV1Response.Value.UserId)             // не пусто
	})

	t.Run("DD2 DescribeDevice requesting data about an NON-existing device", func(t *testing.T) {

		nonDeviceId := int64(10000)                         // передаем не существующий идентификатор
		request := &act_device_api.DescribeDeviceV1Request{ // передаем запрос номер устройства
			DeviceId: uint64(nonDeviceId),
		}

		_, err := act_device_api.NewActDeviceApiServiceClient(conn).DescribeDeviceV1(ctx, request)
		assert.Equal(t, codes.NotFound.String(), status.Code(err).String()) //  проверяем статус код
	})

	t.Run("DD3 DescribeDevice requesting data about a remove device", func(t *testing.T) {

		removeDeviceId := int64(4)                          // передаем идентификатор удаленного объекта
		request := &act_device_api.DescribeDeviceV1Request{ // передаем запрос номер устройства
			DeviceId: uint64(removeDeviceId),
		}

		_, err := act_device_api.NewActDeviceApiServiceClient(conn).DescribeDeviceV1(ctx, request)
		assert.Equal(t, codes.NotFound.String(), status.Code(err).String()) //  проверяем статус код
	})

}
