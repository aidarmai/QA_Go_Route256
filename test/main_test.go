package test

import (
	"github.com/ozonmp/act-device-api/test/config"
	"github.com/ozonmp/act-device-api/test/internal/helpers"
	"net/url"
	"testing"
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
