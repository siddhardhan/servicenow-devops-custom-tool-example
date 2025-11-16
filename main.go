package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title Evidence Service API
// @version 1.0
// @description A service that provides evidence information based on control IDs.
// @host localhost:8080
// @BasePath /v1

// getRandomStatus returns a random evidence status (SUCCESS or FAILED)
func getRandomStatus() string {
	statuses := []string{"SUCCESS", "FAILED"}
	return statuses[mrand.Intn(len(statuses))]
}

// generateSysID generates a random system ID of 32 characters including "sys_" prefix
func generateSysID() string {
	// Generate 16 random bytes and return 32 hex chars prefixed with sys_
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to math/rand seeded by time if crypto/rand fails (very unlikely)
		mrand.Seed(time.Now().UnixNano())
		for i := range b {
			b[i] = byte(mrand.Intn(256))
		}
	}
	return fmt.Sprintf("sys_%s", hex.EncodeToString(b))
}

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Control ID is required"`
	Code    string `json:"code" example:"MISSING_CONTROL_ID"`
}

// Evidence represents the evidence data structure
type Evidence struct {
	EvidenceID     string `json:"evidenceId" example:"sys_1a2b3c4d5e6f7890abcdef0123456789"`
	EvidenceType   string `json:"evidenceType" example:"datadog" enums:"datadog,sonar,gitlab,practitest"`
	ControlID      string `json:"controlId" example:"1234" enums:"1234,5678"`
	EvidenceStatus string `json:"evidenceStatus" example:"SUCCESS" enums:"SUCCESS,FAILED"`
	AppID          string `json:"appId" example:"Corpsite"`
}

// EvidenceTemplate represents the structure of our mock data
type EvidenceTemplate struct {
	EvidenceType string
	ControlID    string
}

// For demonstration, we'll use these templates to generate fresh evidences
var evidenceTemplates = []EvidenceTemplate{
	{EvidenceType: "datadog", ControlID: "1234"},
	{EvidenceType: "sonar", ControlID: "5678"},
	{EvidenceType: "gitlab", ControlID: "9012"},
	{EvidenceType: "practitest", ControlID: "3456"},
}

// @Summary Get evidences by control ID
// @Description Returns a list of evidences filtered by the provided control ID
// @Tags evidence
// @Accept json
// @Produce json
// @Param controlId query string true "Control ID (1234 for datadog, 5678 for Sonar)"
// @Success 200 {array} Evidence
// @Failure 400 {object} ErrorResponse
// @Router /evidences [get]
func getEvidencesHandler(c *gin.Context) {
	// Get the controlId from query parameters
	controlID := c.Query("controlId")
	if controlID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Control ID is required",
			Code:    "MISSING_CONTROL_ID",
		})
		return
	}

	// Generate a random number of evidences (10-26)
	numEvidences := mrand.Intn(17) + 10 // 17 is the range (26-10+1), 10 is the minimum

	// Define specific AppIDs
	appIDs := []string{
		"Corpsite",
		"Cruz Bike Rentals",
		"Hotel Reservation System",
		"Portfolio",
		"Product - Corpsite & Cruz",
		"Smart Shopper",
	}
	// Shuffle the appIDs to ensure random but unique assignment
	mrand.Shuffle(len(appIDs), func(i, j int) {
		appIDs[i], appIDs[j] = appIDs[j], appIDs[i]
	})

	// Filter templates by controlId
	var matchingTemplates []EvidenceTemplate
	for _, template := range evidenceTemplates {
		if template.ControlID == controlID {
			matchingTemplates = append(matchingTemplates, template)
		}
	}

	// Return empty array if no matching templates found
	if len(matchingTemplates) == 0 {
		c.JSON(http.StatusOK, []Evidence{})
		return
	}

	// Generate fresh evidences while ensuring the (ControlID, EvidenceType)
	// pair is unique for each returned evidence. We also respect available
	// appIDs; if we run out of unique combinations or appIDs, we stop early.
	var filteredEvidences []Evidence

	seenPairs := make(map[string]bool)
	evidenceCount := 0

	// Build combinations by iterating appIDs and templates to maximize unique pairs
	// but stop when we hit numEvidences or when combinations are exhausted.
	for i := 0; i < numEvidences; i++ {
		// select appID for this index; if we run out, stop
		if i >= len(appIDs) {
			break
		}
		appID := appIDs[i]

		// Try to find a template that yields an unseen (control,evidenceType) pair
		found := false
		var chosenTemplate EvidenceTemplate
		for t := 0; t < len(matchingTemplates); t++ {
			template := matchingTemplates[(i+t)%len(matchingTemplates)]
			pairKey := template.ControlID + "|" + template.EvidenceType
			if !seenPairs[pairKey] {
				chosenTemplate = template
				seenPairs[pairKey] = true
				found = true
				break
			}
		}

		if !found {
			// No more unique template pairs available; stop generating.
			break
		}

		evidence := Evidence{
			EvidenceID:     generateSysID(),
			EvidenceType:   chosenTemplate.EvidenceType,
			ControlID:      chosenTemplate.ControlID,
			EvidenceStatus: getRandomStatus(),
			AppID:          appID,
		}
		filteredEvidences = append(filteredEvidences, evidence)
		evidenceCount++
	}

	// If we got no evidences, return an empty array instead of null
	if filteredEvidences == nil {
		filteredEvidences = []Evidence{}
	}

	// Return the filtered evidences
	c.JSON(http.StatusOK, filteredEvidences)
}

// @Summary Get evidences by app ID
// @Description Returns a fixed set of evidences for a given app_id. Supports "Hogan" and "SinglePoint".
// @Tags evidence
// @Accept json
// @Produce json
// @Param app_id query string true "App ID (Hogan or SinglePoint)"
// @Success 200 {array} Evidence
// @Failure 400 {object} ErrorResponse
// @Router /evidences/by-app [get]
func getEvidencesByAppIDHandler(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "app_id is required",
			Code:    "MISSING_APP_ID",
		})
		return
	}

	var evidences []Evidence

	// Determine desired counts per app
	var desiredCount int
	switch appID {
	case "Hogan":
		desiredCount = 2
	case "SinglePoint":
		desiredCount = 4
	default:
		// Unknown app_id -> return empty list (200)
		c.JSON(http.StatusOK, []Evidence{})
		return
	}

	// Build unique (ControlID, EvidenceType) pairs. If templates are exhausted
	// create deterministic variations to ensure uniqueness up to desiredCount.
	seen := make(map[string]bool)
	idx := 0
	for len(evidences) < desiredCount {
		base := evidenceTemplates[idx%len(evidenceTemplates)]
		pairKey := base.ControlID + "|" + base.EvidenceType

		tpl := base
		// We now have four allowed templates (dataDog, sonar, gitlab, practitest).
		// Use the next template when available to ensure unique (ControlID, EvidenceType)
		// pairs. Since desiredCount for current callers is <= 4, this will be enough
		// to produce unique entries without creating synthetic types.
		if seen[pairKey] {
			// pick next template in sequence to avoid repeating the same pair
			next := evidenceTemplates[(idx+1)%len(evidenceTemplates)]
			tpl = next
			pairKey = tpl.ControlID + "|" + tpl.EvidenceType
			if seen[pairKey] {
				// fallback: continue to next idx iteration
				idx++
				continue
			}
		}

		seen[pairKey] = true

		status := "SUCCESS"
		// For SinglePoint, make the 3rd evidence FAILED to match requirement
		if appID == "SinglePoint" && len(evidences) == 2 {
			status = "FAILED"
		}

		evidences = append(evidences, Evidence{
			EvidenceID:     generateSysID(),
			EvidenceType:   tpl.EvidenceType,
			ControlID:      tpl.ControlID,
			EvidenceStatus: status,
			AppID:          appID,
		})

		idx++
		// Safety: avoid infinite loop (shouldn't happen)
		if idx > desiredCount*10 {
			break
		}
	}

	c.JSON(http.StatusOK, evidences)
}

func main() {
	// Create a new Gin router with default middleware
	r := gin.Default()

	// Configure Swagger
	config := &ginSwagger.Config{
		URL: "/docs/swagger.json",
	}
	r.GET("/swagger/*any", ginSwagger.CustomWrapHandler(config, swaggerFiles.Handler))
	r.Static("/docs", "./docs")

	// Create v1 route group
	v1 := r.Group("/v1")
	{
		// Register the handler for /v1/evidences endpoint
		v1.GET("/evidences", getEvidencesHandler)
		// Register the handler for /v1/evidences/by-app endpoint
		v1.GET("/evidences/by-app", getEvidencesByAppIDHandler)
	}

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
