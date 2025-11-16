package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	// "strconv"
	"gcp-server/internal/algorithms"
)

type shortestPathRequest struct {
	Graph algorithms.Graph `json:"graph"`
	Src   int              `json:"src"`
}

type shortestPathResponse struct {
	Dist []float64 `json:"dist"`
}

func shortestPathHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "use POST", http.StatusMethodNotAllowed)
		return
	}
	var req shortestPathRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}

	res := algorithms.Dijkstra(req.Graph, req.Src)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shortestPathResponse{Dist: res.Dist})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/shortest-path", shortestPathHandler)

	log.Println("listening on :" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}