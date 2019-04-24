package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestRouterIndex(t *testing.T) {
	router := newRouter()

	server := httptest.NewServer(router)

	resp, err := http.Get(server.URL + "/")

	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Status code incorrect: expected %v, got %v", http.StatusOK, status)
	}

	const expected = "text/html; charset=utf-8"
	actual := resp.Header.Get("Content-Type")

	if actual != expected {
		t.Errorf("Responses did not match: expected %q, got %q", expected, actual)
	}
}

func TestRouterGetSubmit(t *testing.T) {
	router := newRouter()

	server := httptest.NewServer(router)

	resp, err := http.Get(server.URL + "/submit")

	if err != nil {
		t.Fatal(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		t.Errorf("Status code incorrect: expected %v, got %v", http.StatusOK, status)
	}

	const expected = "text/html; charset=utf-8"
	actual := resp.Header.Get("Content-Type")

	if actual != expected {
		t.Errorf("Responses did not match: expected %q, got %q", expected, actual)
	}
}

func TestRouterPostSubmit(t *testing.T) {

	names = map[string]int{
		"test": 1,
	}

	form := url.Values{}
	form.Set("name", "test")

	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(submitHandler)

	hf.ServeHTTP(recorder, req)

	// TODO: change this to the standard
	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := 2

	actual := names["test"]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRouterInvalidRoute(t *testing.T) {
	r := newRouter()
	server := httptest.NewServer(r)

	resp, err := http.Get(server.URL + "/yeet")

	if err != nil {
		t.Fatal(err)
	}

	const expectedStatusCode = http.StatusNotFound
	if status := resp.StatusCode; status != expectedStatusCode {
		t.Errorf("Incorrect status code: expected %v, got %v", expectedStatusCode, status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	expected := "404 page not found\n"
	actual := string(body)

	if actual != expected {
		t.Errorf("Responses did not match: expected %q, got %q", expected, actual)
	}
}

func TestStaticFileServer(t *testing.T) {
	router := newRouter()
	server := httptest.NewServer(router)

	resp, err := http.Get(server.URL + "/assets/")
	if err != nil {
		t.Fatal(err)
	}

	const expectedStatusCode = http.StatusOK
	if actualStatusCode := resp.StatusCode; actualStatusCode != expectedStatusCode {
		t.Errorf("Mismatched status code: expected %v, got %v", expectedStatusCode, actualStatusCode)
	}

	const expectedContentType = "text/html; charset=utf-8"
	actualContentType := resp.Header.Get("Content-Type")

	if actualContentType != expectedContentType {
		t.Errorf("Mismatched content type: expected %v, got %v", expectedContentType, actualContentType)
	}
}
