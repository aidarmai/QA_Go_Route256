package helpers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func IsAlive(checkurl url.URL) {
	response, err := http.Get(checkurl.String())
	if err != nil {
		log.Fatalf("Service is not ready: %v", err.Error())
	} else {
		fmt.Println("Сервис жив")
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("Service status code is not ok, as my WiFi: %v", response.StatusCode)
	}
}
