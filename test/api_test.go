package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

var token string

func TestAddUser(t *testing.T) {
	data := map[string]string{
		"name":     "test",
		"username": "test_user",
		"password": "test_123456",
	}
	jsonStr, err := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/user", strings.NewReader(string(jsonStr)))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("response code is not success: ", response.StatusCode)
	}
}

func TestLoginUser(t *testing.T) {
	data := map[string]string{
		"username": "test_user",
		"password": "test_123456",
	}
	jsonStr, err := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/auth", strings.NewReader(string(jsonStr)))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("response code is not success: ")
	}
	body, err := io.ReadAll(response.Body)

	type ResponseData struct {
		Msg  string `json:"msg"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	var responseData ResponseData
	errParse := json.Unmarshal(body, &responseData)
	if errParse != nil {
		t.Fatal(errParse.Error())
	}
	if responseData.Data.Token == "" {
		t.Fatal("token is empty")
	}
	token = responseData.Data.Token
}

func TestAddData(t *testing.T) {
	data := map[string]interface{}{
		"device": "esp1",
		"result": 2.33,
	}
	jsonStr, err := json.Marshal(data)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/data", strings.NewReader(string(jsonStr)))
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != http.StatusCreated {
		t.Fatal("response code is not success: ")
	}
}

func TestFetchData(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/data?page=1&size=10&device=esp1", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	if response.StatusCode != http.StatusOK {
		t.Fatal("response code is not success: ")
	}
}
