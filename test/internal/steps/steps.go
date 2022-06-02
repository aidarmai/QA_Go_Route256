package steps

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func UpdateDevice(t *testing.T, numDevice int, platform string, userId string) (string, *bytes.Buffer, error) {
	t.Helper()
	type ItemRequest struct {
		Platform string `json:"platform"`
		UserID   string `json:"userId"`
	}
	t.Logf("UpdatedUser id is %s", userId)
	payload := ItemRequest{Platform: platform, UserID: userId}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(payload)
	if err != nil {
		return "", b, err
	}
	urlDevice := fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices/%d", numDevice)
	return urlDevice, b, nil
}

func CreateDevice(t *testing.T, platform string, userId string) (*bytes.Buffer, error) {
	t.Helper()
	type ItemRequest struct {
		Platform string `json:"platform"`
		UserID   string `json:"userId"`
	}
	payload := ItemRequest{Platform: platform, UserID: userId}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(payload)
	if err != nil {
		return b, err
	}
	return b, nil
}
