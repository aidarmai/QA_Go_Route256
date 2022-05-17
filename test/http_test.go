package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ozonmp/act-device-api/test/internal/expects"
	"github.com/ozonmp/act-device-api/test/internal/steps"
	routeclient "github.com/ozonmp/act-device-api/test/route_client"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type ListOfItemsResponse struct {
	Items []struct {
		ID        string     `json:"id"`
		Platform  string     `json:"platform"`
		UserID    string     `json:"userId"`
		EnteredAt *time.Time `json:"enteredAt"`
	} `json:"items"`
}

type ItemRequest struct {
	Platform string `json:"platform"`
	UserID   string `json:"userId"`
}

type RemoveDeviceResponse struct {
	Found bool `json:"found"`
}

type CreateDeviceResponse struct {
	Success bool `json:"success"`
}

func Test_HttpServer_List(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("GET on list return 200", func(t *testing.T) {
		response, err := http.Get("http://127.0.0.1:8080/api/v1/devices?page=1&perPage=1")
		if err != nil {
			panic(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Got %v, but want %v", response.StatusCode, http.StatusOK)
		}
	})

	t.Run("GET on list return devices list", func(t *testing.T) {
		countOfItems := 10
		response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices?page=1&perPage=%d", countOfItems))
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		list := new(ListOfItemsResponse)
		err = json.Unmarshal(data, &list)
		if err != nil {
			panic(err)
		}

		if len(list.Items) != countOfItems {
			t.Errorf("Want %d, get %d items", countOfItems, len(list.Items))
		}
	})

	t.Run("GET on list return devices list if zeroed", func(t *testing.T) {
		countOfItems := 0
		response, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices?page=1&perPage=%d", countOfItems))
		if err != nil {
			panic(err)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		list := new(ListOfItemsResponse)
		err = json.Unmarshal(data, &list)
		if err != nil {
			panic(err)
		}

		if len(list.Items) != countOfItems {
			t.Errorf("Want %d, get %d items", countOfItems, len(list.Items))
		}
	})

	t.Run("POST on creating device", func(t *testing.T) {
		data := []byte(`{"platform": "Android", "userId": "123456"}`)
		r := bytes.NewReader(data)
		contentType := "application/json"

		_, err := http.Post("http://127.0.0.1:8080/api/v1/devices", contentType, r)
		if err != nil {
			panic(err)
		}

		payload := ItemRequest{Platform: "Android", UserID: "123456"}
		payloadJSON, _ := json.Marshal(payload)

		_, err = http.Post("http://127.0.0.1:8080/api/v1/devices", contentType, bytes.NewBuffer(payloadJSON))
		if err != nil {
			panic(err)
		}
	})

	createTableTest := []struct {
		testName string
		platform string
		userId   string
	}{
		{"CD1(1)", "Ios", "200"},
		{"CD1(2)", "Android", "201"},
		{"CD1(3)", "Ubuntu", "202"},
	}

	for _, tt := range createTableTest {
		t.Run(tt.testName+" POST with client", func(t *testing.T) {
			// arrange
			b, _ := steps.CreateDevice(t, tt.platform, tt.userId)
			client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
			// action
			req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/api/v1/devices", b)
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
			//assert

			if res.StatusCode != http.StatusOK {
				t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
			}
			data, _ := ioutil.ReadAll(res.Body)
			if len(data) != 0 {
				t.Log(string(data))
			}
		})
	}
	t.Run("Create device via client API", func(t *testing.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		device := routeclient.CreateDeviceRequest{
			Platform: "Ubuntu",
			UserId:   "701",
		}
		ctx := context.TODO()
		id, _, _ := client.CreateDevice(ctx, device)
		t.Logf("New device is %d", id.DeviceId)
		assert.GreaterOrEqual(t, id.DeviceId, 0)
	})

	t.Run("List devices via client API", func(t *testing.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		opts := url.Values{}
		opts.Add("page", "1")
		opts.Add("perPage", "100")
		ctx := context.TODO()
		items, _, _ := client.ListDevices(ctx, opts)
		assert.GreaterOrEqual(t, len(items.Items), 1)
	})

	t.Run("Delete device via client API", func(t *testing.T) {
		// arrange
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		numDevice := 76
		apiUrl, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices/%d", numDevice))

		// action
		req, err := http.NewRequest(http.MethodDelete, apiUrl.String(), nil)
		if err != nil {
			panic(err)
		}
		res, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		//assert
		if res.StatusCode != http.StatusOK {
			t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
		}

		data, _ := ioutil.ReadAll(res.Body)
		removeTrue, _ := json.Marshal(RemoveDeviceResponse{Found: true})

		expects.ExpDeleteOrUpdateDevice(t, string(removeTrue), string(data))

	})

	updateTableTest := []struct {
		testName  string
		numDevice int
		platform  string
		userId    string
	}{
		{"UD1(1)", 48, "Ios", "100"},
		{"UD1(2)", 49, "Android", "101"},
		{"UD1(3)", 50, "Ubuntu", "102"},
	}

	for _, tt := range updateTableTest {

		t.Run(tt.testName+" Update device via client API", func(t *testing.T) {
			tt := tt
			t.Parallel()
			client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
			urlDevice, b, err := steps.UpdateDevice(t, tt.numDevice, tt.platform, tt.userId)
			if err != nil {
				panic(err)
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

			//assert
			if res.StatusCode != http.StatusOK {
				t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
			}
			data, _ := ioutil.ReadAll(res.Body)
			createTrue, _ := json.Marshal(CreateDeviceResponse{Success: true})

			expects.ExpDeleteOrUpdateDevice(t, string(createTrue), string(data))

		})
	}

}
