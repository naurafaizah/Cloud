package tests

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestPickupEndpoint(t *testing.T) {
	reqBody := PickupRequest{
		OrderID:       "ORD1",
		PaymentStatus: "paid",
		Weight:        2,
	}

	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/pickup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	pickupHandler(w, req)

	var res PickupResponse
	json.NewDecoder(w.Body).Decode(&res)

	if res.Status != "scheduled" {
		t.Errorf("Expected scheduled, got %s", res.Status)
	}
}
