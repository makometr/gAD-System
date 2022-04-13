//go:build integration
// +build integration

package it_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type todoRequestStruct struct {
	Exprs []string `json:"exprs"`
}

type todoResponseStruct struct {
	Answers []string `json:"ans"`
}

func TestEndToEnd(t *testing.T) {
	url := "http://localhost:8080/calc"
	fmt.Println("URL:>", url)

	requestData := todoRequestStruct{
		Exprs: []string{"100+200", "500-600", "25*15", "100/30", "100%3"},
	}
	jsonStr, err := json.Marshal(requestData)
	require.Nil(t, err)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	require.Nil(t, err)

	var actualAns todoResponseStruct
	err = json.Unmarshal(body, &actualAns)
	require.Nil(t, err)
	fmt.Println("response Body:", string(body))

	require.Equal(t, resp.Status, "200 OK")

	expectedAns := todoResponseStruct{
		Answers: []string{"300", "-100", "375", "3", "1"},
	}
	require.Equal(t, expectedAns, actualAns, "answer should be correct")
}
