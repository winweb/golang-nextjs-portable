package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)
func TestHTTPGetUsers(t *testing.T) {
	t.Run("it should return httpCode 200", func(t *testing.T) {

		url := "http://localhost:8080/all"
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Println(string(body))

		if status := resp.StatusCode; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})
}