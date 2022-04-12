//go:build integration
// +build integration

package it_test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

type todoAnswetStruct struct {
	Answers []string `json:"ans"`
}

func TestEndToEnd(t *testing.T) {
	url := "http://localhost:8080/calc"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{
		"exprs": [
			"100+200",
			"500-600",
			"24*15",
			"100/30",
			"100%3"
		]
	}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	require.Equal(t, resp.Status, "200 OK")

}
