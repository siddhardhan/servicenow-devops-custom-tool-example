package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// generateSysID generates a random system ID of 32 characters including "sys_" prefix
func generateSysID() string {
    const charset = "abcdef0123456789"
    rand.Seed(time.Now().UnixNano())
    sysID := make([]byte, 32) // 32 characters
    for i := range sysID {
        sysID[i] = charset[rand.Intn(len(charset))]
    }
    return fmt.Sprintf("sys_%s", string(sysID))
}

// Evidence represents the evidence data structure
type Evidence struct {
    EvidenceID     string `json:"evidenceId"`
    EvidenceType   string `json:"evidenceType"`
    ControlID      string `json:"controlId"`
    EvidenceStatus string `json:"evidenceStatus"`
}

// EvidenceTemplate represents the structure of our mock data
type EvidenceTemplate struct {
    EvidenceType string
    ControlID    string
}

// For demonstration, we'll use these templates to generate fresh evidences
var evidenceTemplates = []EvidenceTemplate{
    {EvidenceType: "document", ControlID: "1234"},
    {EvidenceType: "audit", ControlID: "1234"},
    {EvidenceType: "report", ControlID: "1234"},
    {EvidenceType: "compliance", ControlID: "1234"},
    {EvidenceType: "security", ControlID: "5678"},
}

// generateEvidence creates a new evidence from a template with a fresh sysID
func generateEvidence(template EvidenceTemplate) Evidence {
    return Evidence{
        EvidenceID:     generateSysID(),
        EvidenceType:   template.EvidenceType,
        ControlID:      template.ControlID,
        EvidenceStatus: "SUCCESS",
    }
}

// getEvidencesHandler handles the /getEvidences endpoint
func getEvidencesHandler(c *gin.Context) {
    // Get the controlId from query parameters
    controlID := c.Query("controlId")
    if controlID == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "controlId is required",
        })
        return
    }

    // Generate fresh evidences and filter by controlId
    var filteredEvidences []Evidence
    for _, template := range evidenceTemplates {
        if template.ControlID == controlID {
            evidence := generateEvidence(template)
            filteredEvidences = append(filteredEvidences, evidence)
        }
    }

    // Return the filtered evidences
    c.JSON(http.StatusOK, filteredEvidences)
}

func main() {
    // Create a new Gin router with default middleware
    r := gin.Default()

    // Register the handler for /getEvidences endpoint
    r.GET("/getEvidences", getEvidencesHandler)

    // Start the server
    fmt.Println("Starting server on :8080")
    if err := r.Run(":8080"); err != nil {
        fmt.Printf("Failed to start server: %v\n", err)
    }
}
