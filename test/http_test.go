package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/ozonmp/act-device-api/test/internal/expects"
	"github.com/ozonmp/act-device-api/test/internal/steps"
	routeclient "github.com/ozonmp/act-device-api/test/route_client"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
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

func Test_HttpServer(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	runner.Run(t, "GET on list return 200", func(t provider.T) {
		response, err := http.Get("http://127.0.0.1:8080/api/v1/devices?page=1&perPage=1")
		if err != nil {
			panic(err)
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Got %v, but want %v", response.StatusCode, http.StatusOK)
		}
	})

	runner.Run(t, "GET on list return devices list", func(t provider.T) {
		page := 1
		countOfItems := 10
		urlStr := fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices?page=%d&perPage=%d", page, countOfItems)
		response, err := http.Get(urlStr)
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

		t.WithNewAttachment("request", allure.Text, []byte(urlStr))
		t.WithNewAttachment("response", allure.Text, data)

		if len(list.Items) != countOfItems {
			t.Errorf("Want %d, get %d items", countOfItems, len(list.Items))
		}
	})

	runner.Run(t, "GET on list return devices list if zeroed", func(t provider.T) {
		page := 1
		countOfItems := 0
		urlStr := fmt.Sprintf("http://127.0.0.1:8080/api/v1/devices?page=%d&perPage=%d", page, countOfItems)
		response, err := http.Get(urlStr)
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

		t.WithNewAttachment("request", allure.Text, []byte(urlStr))
		t.WithNewAttachment("response", allure.Text, data)

		if len(list.Items) != countOfItems {
			t.Errorf("Want %d, get %d items", countOfItems, len(list.Items))
		}
	})

	runner.Run(t, "POST on creating device", func(t provider.T) {
		data := []byte(`{"platform": "Android", "userId": "123456"}`)
		r := bytes.NewReader(data)
		contentType := "application/json"

		_, err := http.Post("http://127.0.0.1:8080/api/v1/devices", contentType, r)
		if err != nil {
			panic(err)
		}

		payload := ItemRequest{Platform: "Android", UserID: "123456"}
		payloadJSON, _ := json.Marshal(payload)

		t.WithNewAttachment("request", allure.Text, payloadJSON)
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
		runner.Run(t, tt.testName+" POST with client", func(t provider.T) {
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

			t.NewStep("response status check")
			if res.StatusCode != http.StatusOK {
				t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
			}

			t.NewStep("checking for non-empty response")
			data, _ := ioutil.ReadAll(res.Body)
			t.WithNewAttachment("response", allure.Text, data)
			if len(data) != 0 {
				t.Log(string(data))
			}
		})
	}
	runner.Run(t, "Create device via client API", func(t provider.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		device := routeclient.CreateDeviceRequest{
			Platform: "Ubuntu",
			UserId:   "701",
		}
		ctx := context.TODO()
		id, _, _ := client.CreateDevice(ctx, device)
		t.Logf("New device is %d", id.DeviceId)
		t.Assert().GreaterOrEqual(id.DeviceId, 0)
	})

	runner.Run(t, "List devices via client API", func(t provider.T) {
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		opts := url.Values{}
		opts.Add("page", "1")
		opts.Add("perPage", "100")
		ctx := context.TODO()
		items, _, _ := client.ListDevices(ctx, opts)
		t.Assert().GreaterOrEqual(len(items.Items), 1)
	})

	runner.Run(t, "Delete device via client API", func(t provider.T) {
		// arrange
		client := routeclient.NewHTTPClient("http://127.0.0.1:8080", 5, 1*time.Second)
		numDevice := 78
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
		t.NewStep("response status check")
		if res.StatusCode != http.StatusOK {
			t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
		}

		data, _ := ioutil.ReadAll(res.Body)
		removeTrue, _ := json.Marshal(RemoveDeviceResponse{Found: true})

		t.NewStep("comparing the result with the expected")
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

		runner.Run(t, tt.testName+" Update device via client API", func(t provider.T) {
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
			t.NewStep("response status check")
			if res.StatusCode != http.StatusOK {
				t.Errorf("Got %v, but want %v", res.StatusCode, http.StatusOK)
			}
			data, _ := ioutil.ReadAll(res.Body)
			createTrue, _ := json.Marshal(CreateDeviceResponse{Success: true})

			t.NewStep("comparing the result with the expected")
			expects.ExpDeleteOrUpdateDevice(t, string(createTrue), string(data))

		})
	}

}
