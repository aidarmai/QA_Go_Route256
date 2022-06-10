package test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/act-device-api/platform"
	"github.com/stretchr/testify/assert"
)

type mockDeviceDB struct{}

func (m *mockDeviceDB) CountPlatform() (int64, error) {
	return 1000, nil
}

func (m *mockDeviceDB) CountPlatformCertainType(platform string) (int64, error) {
	return 600, nil
}

func TestPercentagePlatformCertainType(t *testing.T) {
	//заглушка
	m := &mockDeviceDB{}
	// используем заглушку
	resultPercentage, _ := platform.PercentagePlatformCertainType(m)

	assert.Equal(t, "60.0 %", resultPercentage)

	// для использования с моками
	PlatformTable := []struct {
		pmName     string
		pTypeCount int64
		pAllCount  int64
	}{
		{"ios", 3000, 17605},
		{"android", 4256, 5321},
		{"windows", 7369, 8000},
	}

	//используем моки
	for _, tt := range PlatformTable {
		t.Run("Test PercentagePlatformCertainType using mocks", func(t *testing.T) {
			tt := tt
			t.Parallel()
			ctrl := gomock.NewController(t)
			mock := NewMockDeviceModel(ctrl)

			mock.EXPECT().CountPlatformCertainType("ios").Return(tt.pTypeCount, nil).AnyTimes()
			mock.EXPECT().CountPlatform().Return(tt.pAllCount, nil).AnyTimes()
			res, _ := platform.PercentagePlatformCertainType(mock)
			expRes := float64(tt.pTypeCount * 100 / tt.pAllCount)

			assert.Equal(t, (fmt.Sprintf("%.1f %s", expRes, "%")), res)
		})
	}
}
