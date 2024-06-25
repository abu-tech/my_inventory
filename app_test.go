package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize()

	if err != nil {
		log.Fatal("Error occured while running the server")
	}

	m.Run()
}

func clearTable() {
	// Delete all rows
	_, err := a.Db.Exec("DELETE FROM products")
	if err != nil {
		log.Printf("Error clearing table: %v", err)
	}

	// Reset the sequence
	_, err = a.Db.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
	if err != nil {
		log.Printf("Error resetting sequence: %v", err)
	}
}

func addData(name string, quantity int, price int) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) values('%v', %v, %v)", name, quantity, price)
	_, err := a.Db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func checkStatus(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected code %v, Recieved %v", expected, actual)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addData("chair", 100, 300)
	req, err := http.NewRequest("GET", "/product/1", nil)

	if err != nil {
		t.Fatal(err)
	}

	response := sendRequest(req)
	checkStatus(t, http.StatusOK, response.Code)
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	// 1. Prepare the request payload
	payload := []byte(`{"name":"pen", "quantity":10, "price":15}`)

	// 2. Create a new request with the payload
	req, err := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	// 3. Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// 4. Send the request and get the response
	response := sendRequest(req)

	// 5. Check the response status code
	checkStatus(t, http.StatusCreated, response.Code)

	// 6. Decode the response body
	var m map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 7. Check the returned data
	if m["name"] != "pen" {
		t.Errorf("Expected product name to be 'pen'. Got '%v'", m["name"])
	}

	//and so on
}

func TestDeleteProuct(t *testing.T) {
	clearTable()
	addData("chair", 100, 300)

	//fetch
	req, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := sendRequest(req)
	checkStatus(t, http.StatusOK, response.Code)

	//delete
	req, err = http.NewRequest("DELETE", "/product/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response = sendRequest(req)
	checkStatus(t, http.StatusOK, response.Code)

	//again fetch
	req, err = http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response = sendRequest(req)
	checkStatus(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addData("phone", 100, 1000)

	//fetch
	req, err := http.NewRequest("GET", "/product/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	response := sendRequest(req)
	checkStatus(t, http.StatusOK, response.Code)

	//store the old response
	var oldValue map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &oldValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	//update product
	payload := []byte(`{"name":"pen", "quantity":10, "price":15}`)

	req, err = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)
	checkStatus(t, http.StatusOK, response.Code)

	//store the new response
	var newValue map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &newValue)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	//check the cases
	if oldValue["name"] == newValue["name"] {
		t.Errorf("Expected name %v, Recieved %v", oldValue["name"], newValue["name"])
	}

	//and so on
}
