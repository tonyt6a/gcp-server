package main

import (
	"encoding/json"
	"log"
	"net/http"

	"gcp-server/internal/scheduler"

	"github.com/google/uuid"
)

type createJobRequest struct {
	Priority int    `json:"priority"`
	Payload  string `json:"payload"`
}

type nextJobRequest struct {
	WorkerID string `json:"worker_id"`
}

type completeJobRequest struct {
	WorkerID string `json:"worker_id"`
}

func main() {
	sched := scheduler.NewScheduler()

	mux := http.NewServeMux()

	// ----------------------------------
	// Health check
	// ----------------------------------
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// ----------------------------------
	// POST /jobs — Create a job
	// ----------------------------------
	mux.HandleFunc("POST /jobs", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req createJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.Priority == 0 {
			req.Priority = 10 // default priority
		}

		job := sched.CreateJob(req.Priority, req.Payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(job)
	})

	// ----------------------------------
	// GET /jobs/{id} — Fetch a job by ID
	// ----------------------------------
	mux.HandleFunc("GET /jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		jobID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
			return
		}

		job, ok := sched.GetJob(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	})

	// ----------------------------------
	// POST /jobs/next — Worker asks for job
	// ----------------------------------
	mux.HandleFunc("POST /jobs/next", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req nextJobRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		job := sched.NextJob(req.WorkerID)
		if job == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(job)
	})

	// ----------------------------------
	// POST /jobs/{id}/complete — Worker completes job
	// ----------------------------------
	mux.HandleFunc("POST /jobs/{id}/complete", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		jobID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()
		var req completeJobRequest
		json.NewDecoder(r.Body).Decode(&req) // ignore EOF

		_, ok := sched.CompleteJob(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// ----------------------------------
	// Start service
	// ----------------------------------
	addr := ":8080"
	log.Printf("scheduler listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}