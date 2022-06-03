package test

import (
	"net/url"
	"testing"

	"github.com/ozonmp/act-device-api/test/config"
	"github.com/ozonmp/act-device-api/test/internal/helpers"
)

func TestMain(m *testing.M) {
	cfg, _ := config.GetConfig()

	helpers.IsAlive(url.URL{
		Scheme: "http",
		Host:   cfg.ApiHost,
		Path:   cfg.LivecheckURI,
	})
	m.Run()
}
