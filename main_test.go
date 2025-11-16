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
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=Hogan", nil)
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

	for _, e := range resp {
		if e.EvidenceStatus != "SUCCESS" {
			t.Fatalf("expected evidence status SUCCESS for Hogan, got %s", e.EvidenceStatus)
		}
		// ensure evidence type is allowed
		if e.EvidenceType != "datadog" && e.EvidenceType != "sonar" && e.EvidenceType != "gitlab" && e.EvidenceType != "practitest" {
			t.Fatalf("unexpected evidence type: %s", e.EvidenceType)
		}
		if len(e.ControlID) != 4 {
			t.Fatalf("controlId must be 4 digits, got %s", e.ControlID)
		}
	}
}

func TestGetEvidencesByAppID_SinglePoint(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/v1/evidences/by-app?app_id=SinglePoint", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp []Evidence
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp) != 4 {
		t.Fatalf("expected 4 evidences for SinglePoint, got %d", len(resp))
	}

	failedCount := 0
	for _, e := range resp {
		if e.EvidenceStatus == "FAILED" {
			failedCount++
		}
		if e.EvidenceType != "datadog" && e.EvidenceType != "sonar" && e.EvidenceType != "gitlab" && e.EvidenceType != "practitest" {
			t.Fatalf("unexpected evidence type: %s", e.EvidenceType)
		}
		if len(e.ControlID) != 4 {
			t.Fatalf("controlId must be 4 digits, got %s", e.ControlID)
		}
	}

	if failedCount != 1 {
		t.Fatalf("expected exactly 1 FAILED evidence for SinglePoint, got %d", failedCount)
	}
}
