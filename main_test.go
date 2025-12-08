package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.New()
	v1 := r.Group("/v1")
	{
		v1.GET("/evidences", getEvidencesHandler)
		v1.GET("/evidences/by-app", getEvidencesByAppIDHandler)
	}
	return r
}

func TestGetEvidencesByAppID_Hogan(t *testing.T) {
	// Clear the request tracker before this test
	requestTracker = make(map[string]int)

	r := setupRouter()

	// First request (odd) - should have at least one FAILED
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=Hogan&control_ids=1234,5678", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []Evidence
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 2 {
		t.Fatalf("expected 2 evidences for Hogan, got %d", len(resp))
	}

	failedCount := 0
	for _, e := range resp {
		if e.EvidenceStatus == "FAILED" {
			failedCount++
		}
		if e.AppID != "Hogan" {
			t.Fatalf("expected app_id Hogan, got %s", e.AppID)
		}
		if e.EvidenceType != "datadog" && e.EvidenceType != "sonar" && e.EvidenceType != "gitlab" && e.EvidenceType != "practitest" {
			t.Fatalf("unexpected evidence type: %s", e.EvidenceType)
		}
		if len(e.ControlID) != 4 {
			t.Fatalf("controlId must be 4 digits, got %s", e.ControlID)
		}
	}

	if failedCount < 1 {
		t.Fatalf("expected at least 1 FAILED evidence on first (odd) request, got %d", failedCount)
	}

	// Second request (even) - should all be SUCCESS
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=Hogan&control_ids=1234,5678", nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	var resp2 []Evidence
	if err := json.Unmarshal(w2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp2) != 2 {
		t.Fatalf("expected 2 evidences for Hogan on second request, got %d", len(resp2))
	}

	for _, e := range resp2 {
		if e.EvidenceStatus != "SUCCESS" {
			t.Fatalf("expected all SUCCESS on second (even) request, got %s", e.EvidenceStatus)
		}
	}

	// Third request (odd) - should have at least one FAILED again
	w3 := httptest.NewRecorder()
	req3 := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=Hogan&control_ids=1234,5678", nil)
	r.ServeHTTP(w3, req3)

	if w3.Code != 200 {
		t.Fatalf("expected 200, got %d", w3.Code)
	}

	var resp3 []Evidence
	if err := json.Unmarshal(w3.Body.Bytes(), &resp3); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	failedCount3 := 0
	for _, e := range resp3 {
		if e.EvidenceStatus == "FAILED" {
			failedCount3++
		}
	}

	if failedCount3 < 1 {
		t.Fatalf("expected at least 1 FAILED evidence on third (odd) request, got %d", failedCount3)
	}
}

func TestGetEvidencesByAppID_SinglePoint(t *testing.T) {
	// Clear the request tracker before this test
	requestTracker = make(map[string]int)

	r := setupRouter()

	// First request (odd) - should have at least one FAILED
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=SinglePoint&control_ids=1234,5678,9012,3456", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []Evidence
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 4 {
		t.Fatalf("expected 4 evidences for SinglePoint (one per control ID), got %d", len(resp))
	}

	failedCount := 0
	for _, e := range resp {
		if e.EvidenceStatus == "FAILED" {
			failedCount++
		}
		if e.AppID != "SinglePoint" {
			t.Fatalf("expected app_id SinglePoint, got %s", e.AppID)
		}
		if e.EvidenceType != "datadog" && e.EvidenceType != "sonar" && e.EvidenceType != "gitlab" && e.EvidenceType != "practitest" {
			t.Fatalf("unexpected evidence type: %s", e.EvidenceType)
		}
		if len(e.ControlID) != 4 {
			t.Fatalf("controlId must be 4 digits, got %s", e.ControlID)
		}
	}

	if failedCount < 1 {
		t.Fatalf("expected at least 1 FAILED evidence on first (odd) request, got %d", failedCount)
	}

	// Second request (even) - should all be SUCCESS
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=SinglePoint&control_ids=1234,5678,9012,3456", nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

	var resp2 []Evidence
	if err := json.Unmarshal(w2.Body.Bytes(), &resp2); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp2) != 4 {
		t.Fatalf("expected 4 evidences for SinglePoint on second request, got %d", len(resp2))
	}

	for _, e := range resp2 {
		if e.EvidenceStatus != "SUCCESS" {
			t.Fatalf("expected all SUCCESS on second (even) request, got %s", e.EvidenceStatus)
		}
	}
}

func TestVersion_R22_2_0_AllSuccess(t *testing.T) {
	// Clear the request tracker before this test
	requestTracker = make(map[string]int)

	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=SinglePoint&control_ids=1234,5678,9012,3456&version=R22-2.0", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []Evidence
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) == 0 {
		t.Fatalf("expected at least 1 evidence for R22-2.0, got 0")
	}

	for _, e := range resp {
		if e.EvidenceStatus != "SUCCESS" {
			t.Fatalf("expected all SUCCESS for R22-2.0, got %s", e.EvidenceStatus)
		}
	}
}

func TestVersion_R22_1_0_AtLeastOneFailed(t *testing.T) {
	// Clear the request tracker before this test
	requestTracker = make(map[string]int)

	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=SinglePoint&control_ids=1234,5678,9012,3456&version=R22-1.0", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []Evidence
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) == 0 {
		t.Fatalf("expected at least 1 evidence for R22-1.0, got 0")
	}

	failedCount := 0
	for _, e := range resp {
		if e.EvidenceStatus == "FAILED" {
			failedCount++
		}
	}

	if failedCount < 1 {
		t.Fatalf("expected at least 1 FAILED evidence for R22-1.0, got %d", failedCount)
	}
}
