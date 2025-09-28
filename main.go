package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

// getRandomStatus returns a random evidence status (SUCCESS or FAILED)
func getRandomStatus() string {
    statuses := []string{"SUCCESS", "FAILED"}
    return statuses[rand.Intn(len(statuses))]
}

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
    AppID          string `json:"appId"`
}

// EvidenceTemplate represents the structure of our mock data
type EvidenceTemplate struct {
    EvidenceType string
    ControlID    string
}

// For demonstration, we'll use these templates to generate fresh evidences
var evidenceTemplates = []EvidenceTemplate{
    {EvidenceType: "dataDog", ControlID: "1234"},
    {EvidenceType: "sonar", ControlID: "5678"},
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

    // Generate a random number of evidences (10-26)
    numEvidences := rand.Intn(17) + 10  // 17 is the range (26-10+1), 10 is the minimum

    // Generate appIDs from A to Z
    appIDs := make([]string, 26)
    for i := 0; i < 26; i++ {
        appIDs[i] = string(rune('A' + i))
    }
    // Shuffle the appIDs to ensure random but unique assignment
    rand.Shuffle(len(appIDs), func(i, j int) {
        appIDs[i], appIDs[j] = appIDs[j], appIDs[i]
    })

    // Filter templates by controlId
    var matchingTemplates []EvidenceTemplate
    for _, template := range evidenceTemplates {
        if template.ControlID == controlID {
            matchingTemplates = append(matchingTemplates, template)
        }
    }

    // Generate fresh evidences
    var filteredEvidences []Evidence
    evidenceCount := 0

    // Keep adding evidences until we reach the required number
    for evidenceCount < numEvidences {
        // Cycle through the matching templates
        templateIndex := evidenceCount % len(matchingTemplates)
        template := matchingTemplates[templateIndex]

        // Only use appIDs if we haven't run out
        var appID string
        if evidenceCount < len(appIDs) {
            appID = appIDs[evidenceCount]
        } else {
            // If we need more evidences than available appIDs, don't add more
            break
        }

        evidence := Evidence{
            EvidenceID:     generateSysID(),
            EvidenceType:   template.EvidenceType,
            ControlID:      template.ControlID,
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
