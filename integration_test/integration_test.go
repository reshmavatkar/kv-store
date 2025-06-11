package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://localhost:8080/store"

func TestIntegrationPutGetDelete(t *testing.T) {
	key := "mykey"
	value := "myvalue"

	// ----- PUT -----
	putPayload := map[string]string{"key": key, "value": value}
	body, _ := json.Marshal(putPayload)
	req, _ := http.NewRequest(http.MethodPut, baseURL, bytes.NewBuffer(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("PUT failed: %v, status: %d", err, resp.StatusCode)
	}
	resp.Body.Close()

	// ----- GET -----
	resp, err = http.Get(baseURL + "/" + key)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET failed: %v, status: %d", err, resp.StatusCode)
	}
	getResp, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Failed to read GET response: %v", err)
	}
	var getBody map[string]string
	if err := json.Unmarshal(getResp, &getBody); err != nil {
		t.Fatalf("Failed to unmarshal GET response: %v", err)
	}
	if getBody["value"] != value {
		t.Fatalf("Expected value %s, got %s", value, getBody["value"])
	}

	// ----- DELETE -----
	req, _ = http.NewRequest("DELETE", baseURL+"/"+key, nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("DELETE failed: %v, status: %d", err, resp.StatusCode)
	}
	resp.Body.Close()

	// ----- GET after DELETE -----
	resp, err = http.Get(baseURL + "/" + key)
	if err != nil {
		t.Fatalf("GET after DELETE failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("Expected 404 after DELETE, got %d", resp.StatusCode)
	}
}
