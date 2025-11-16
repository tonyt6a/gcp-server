package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Job matches what scheduler returns.
type Job struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Priority  int    `json:"priority"`
	Payload   string `json:"payload"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {

	// ---------------------------------------------------------
	// Environment Variables
	// ---------------------------------------------------------
	schedulerURL := os.Getenv("SCHEDULER_URL")
	if schedulerURL == "" {
		schedulerURL = "http://localhost:8080" // local default
	}

	workerID := os.Getenv("WORKER_ID")
	if workerID == "" {
		workerID = "worker-local"
	}

	log.Printf("[worker %s] Starting. scheduler=%s", workerID, schedulerURL)

	// ---------------------------------------------------------
	// Poll loop
	// ---------------------------------------------------------
	for {
		job, err := fetchNextJob(schedulerURL, workerID)
		if err != nil {
			log.Printf("[worker %s] Error fetching job: %v", workerID, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// No job available
		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("[worker %s] Got job %s (priority=%d payload=%q)",
			workerID, job.ID, job.Priority, job.Payload)

		// -----------------------------------------------------
		// Simulate doing work
		// -----------------------------------------------------
		time.Sleep(2 * time.Second)

		// -----------------------------------------------------
		// Notify scheduler that work is complete
		// -----------------------------------------------------
		if err := completeJob(schedulerURL, workerID, job.ID); err != nil {
			log.Printf("[worker %s] Error completing job %s: %v", workerID, job.ID, err)
		} else {
			log.Printf("[worker %s] Completed job %s", workerID, job.ID)
		}
	}
}

//
// ---------------------------------------------------------
// Helper: Worker -> Scheduler GET next job
// ---------------------------------------------------------
func fetchNextJob(baseURL, workerID string) (*Job, error) {
	body, _ := json.Marshal(map[string]string{
		"worker_id": workerID,
	})

	resp, err := http.Post(baseURL+"/jobs/next", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// No job available
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		// Unexpected, but not fatal
		return nil, nil
	}

	// Decode job response
	var job Job
	if err := json.NewDecoder(resp.Body).Decode(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

//
// ---------------------------------------------------------
// Helper: Worker -> Scheduler POST complete job
// ---------------------------------------------------------
func completeJob(baseURL, workerID, jobID string) error {
	body, _ := json.Marshal(map[string]string{
		"worker_id": workerID,
	})

	req, _ := http.NewRequest("POST", baseURL+"/jobs/"+jobID+"/complete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
