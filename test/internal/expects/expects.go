package expects

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"

	"github.com/stretchr/testify/assert"
)

func ExpDeleteOrUpdateDevice(t provider.T, wantBody string, expectBody string) {
	t.Helper()
	assert.Equal(t, wantBody, expectBody)
}
