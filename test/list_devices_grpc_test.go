//go:build gRPCtest
// +build gRPCtest

package test

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	act_device_api "github.com/ozonmp/act-device-api/pkg/act-device-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDevice(t *testing.T) {
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

	t.Run("LD1 ListDevice display device on page with a given number entry", func(t *testing.T) {

		page := uint64(1)
		perPages := uint64(5)
		request := &act_device_api.ListDevicesV1Request{
			Page:    page,
			PerPage: perPages,
		}

		listDevicesV1Response, err := act_device_api.NewActDeviceApiServiceClient(conn).ListDevicesV1(ctx, request)
		assert.Equal(t, codes.OK.String(), status.Code(err).String())
		require.NoError(t, err)
		require.NotNil(t, listDevicesV1Response)
		assert.Equal(t, len(listDevicesV1Response.Items), int(perPages)) // кол-во записей = perPage
		for _, value := range listDevicesV1Response.Items {              //перебираем тела ответов
			require.NotEmpty(t, value.Platform) // поле платформа не пустое
			require.NotEmpty(t, value.UserId)   // поле идентификатор пользователя не пустое
		}

	})

	t.Run("LD2 ListDevice incorrect request", func(t *testing.T) {
		var x uint64 // пустое значение
		page := []uint64{1, x, x}
		perPages := []uint64{x, 1, x}
		for i := 0; i < len(page); i++ {
			request := &act_device_api.ListDevicesV1Request{
				Page:    page[i],
				PerPage: perPages[i],
			}

			listDevicesV1Response, _ := act_device_api.NewActDeviceApiServiceClient(conn).ListDevicesV1(ctx, request)
			assert.Emptyf(t, listDevicesV1Response.GetItems(), "listDevicesV1Response.Items - не пустой")

		}
	})
}
