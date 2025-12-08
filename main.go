package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"net/http"
	"strings"
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
	Version        string `json:"version,omitempty" example:"v1.2.3"`
}

// Allowed evidence types
var allowedEvidenceTypes = []string{"datadog", "sonar", "gitlab", "practitest"}

// requestTracker keeps track of requests by app_id and control_ids combination
// key format: "appid|controlid1,controlid2,..."
var requestTracker = make(map[string]int)

// getRandomEvidenceType returns a random evidence type from the allowed list
func getRandomEvidenceType() string {
	return allowedEvidenceTypes[mrand.Intn(len(allowedEvidenceTypes))]
}

// @Summary Get evidences by control ID
// @Description Returns a list of evidences filtered by the provided control ID
// @Tags evidence
// @Accept json
// @Produce json
// @Param control_ids query string true "Control ID(s). Accepts a single ID or a comma-separated list (e.g. 1234 or 1234,5678)"
// @Success 200 {array} Evidence
// @Failure 400 {object} ErrorResponse
// @Router /evidences [get]
func getEvidencesHandler(c *gin.Context) {
	// Get the controlId from query parameters
	controlIDsParam := c.Query("control_ids")
	if controlIDsParam == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "control_ids is required (single or comma-separated list)",
			Code:    "MISSING_CONTROL_ID",
		})
		return
	}

	// support comma-separated control IDs
	requested := make(map[string]bool)
	for _, id := range strings.Split(controlIDsParam, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			requested[id] = true
		}
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

	// Return empty array if no control IDs requested
	if len(requested) == 0 {
		c.JSON(http.StatusOK, []Evidence{})
		return
	}

	// Generate fresh evidences with random evidence types
	var filteredEvidences []Evidence

	// Build combinations by iterating appIDs and control IDs
	for i := 0; i < numEvidences; i++ {
		// select appID for this index; if we run out, stop
		if i >= len(appIDs) {
			break
		}
		appID := appIDs[i]

		// Rotate through requested control IDs
		cidList := make([]string, 0, len(requested))
		for cid := range requested {
			cidList = append(cidList, cid)
		}
		controlID := cidList[i%len(cidList)]

		evidence := Evidence{
			EvidenceID:     generateSysID(),
			EvidenceType:   getRandomEvidenceType(),
			ControlID:      controlID,
			EvidenceStatus: getRandomStatus(),
			AppID:          appID,
		}
		filteredEvidences = append(filteredEvidences, evidence)
	}

	// If we got no evidences, return an empty array instead of null
	if filteredEvidences == nil {
		filteredEvidences = []Evidence{}
	}

	// Return the filtered evidences
	c.JSON(http.StatusOK, filteredEvidences)
}

// @Summary Get evidences by app ID
// @Description Returns a set of evidences for a given app_id. Supports "Hogan" and "SinglePoint".
// @Tags evidence
// @Accept json
// @Produce json
// @Param app_id query string true "App ID (Hogan or SinglePoint)"
// @Param control_ids query string true "Comma-separated list of control IDs to include (e.g. 1234,5678)"
// @Param version query string false "Optional version string to include in the evidence (e.g. v1.2.3)"
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

	// New: controlIds is a mandatory, comma-separated list of control IDs to include
	controlIdsParam := c.Query("control_ids")
	if controlIdsParam == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "control_ids is required (comma-separated list)",
			Code:    "MISSING_CONTROL_IDS",
		})
		return
	}

	// Optional version parameter: will be echoed into the returned Evidence objects
	versionParam := c.Query("version")

	// Parse requested control IDs and normalize them (trim and sort for consistent tracking key)
	var controlIDsList []string
	for _, id := range strings.Split(controlIdsParam, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			controlIDsList = append(controlIDsList, id)
		}
	}

	if len(controlIDsList) == 0 {
		c.JSON(http.StatusOK, []Evidence{})
		return
	}

	// Create a tracking key to identify unique combinations of app_id, control_ids, and version
	trackingKey := appID + "|" + strings.Join(controlIDsList, ",") + "|" + versionParam

	// Increment request count for this combination
	requestTracker[trackingKey]++
	requestCount := requestTracker[trackingKey]

	// Generate ONE evidence per control ID with random evidence type
	var evidences []Evidence

	for _, controlID := range controlIDsList {
		// Determine status based on request count
		// Odd requests (1st, 3rd, 5th...): at least one FAILED
		// Even requests (2nd, 4th, 6th...): all SUCCESS
		status := "SUCCESS"
		if requestCount%2 == 1 {
			// Odd request: make at least one evidence FAILED
			// We'll make the first one FAILED
			if len(evidences) == 0 {
				status = "FAILED"
			}
		}

		ev := Evidence{
			EvidenceID:     generateSysID(),
			EvidenceType:   getRandomEvidenceType(),
			ControlID:      controlID,
			EvidenceStatus: status,
			AppID:          appID,
		}
		if versionParam != "" {
			ev.Version = versionParam
		}

		evidences = append(evidences, ev)
	}

	// Apply version-specific overrides if needed
	switch versionParam {
	case "R22-2.0":
		// Force all evidences to SUCCESS for this version
		for i := range evidences {
			evidences[i].EvidenceStatus = "SUCCESS"
		}
	case "R22-1.0":
		// Ensure at least one FAILED evidence exists for the requested control IDs
		hasFailed := false
		for _, ev := range evidences {
			if ev.EvidenceStatus == "FAILED" {
				hasFailed = true
				break
			}
		}
		if !hasFailed && len(evidences) > 0 {
			// If none failed yet, mark the first evidence as FAILED to satisfy requirement
			evidences[0].EvidenceStatus = "FAILED"
		}
	case "":
		// No version provided: use the request count logic (already applied above)
		// No additional changes needed
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
