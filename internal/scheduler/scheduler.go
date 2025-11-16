package scheduler

import (
	"container/heap"
	"sync"
	"time"

	"github.com/google/uuid"
)

type priorityJob struct {
	job   *Job
	index int
}

type priorityQueue []*priorityJob

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// Lower priority value = higher priority. Tie-breaker: older CreatedAt first.
	if pq[i].job.Priority == pq[j].job.Priority {
		return pq[i].job.CreatedAt.Before(pq[j].job.CreatedAt)
	}
	return pq[i].job.Priority < pq[j].job.Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*priorityJob)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

type Scheduler struct {
	mu   sync.Mutex
	jobs map[uuid.UUID]*Job
	pq   priorityQueue
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		jobs: make(map[uuid.UUID]*Job),
		pq:   make(priorityQueue, 0),
	}
}

func (s *Scheduler) CreateJob(priority int, payload string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	j := &Job{
		ID:        uuid.New(),
		Status:    StatusPending,
		Priority:  priority,
		Payload:   payload,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.jobs[j.ID] = j
	heap.Push(&s.pq, &priorityJob{job: j})
	return j
}

func (s *Scheduler) NextJob(workerID string) *Job {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.pq) == 0 {
		return nil
	}
	item := heap.Pop(&s.pq).(*priorityJob)
	j := item.job
	j.Status = StatusRunning
	j.UpdatedAt = time.Now()
	return j
}

func (s *Scheduler) CompleteJob(id uuid.UUID) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	j, ok := s.jobs[id]
	if !ok {
		return nil, false
	}
	j.Status = StatusCompleted
	j.UpdatedAt = time.Now()
	return j, true
}

func (s *Scheduler) GetJob(id uuid.UUID) (*Job, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	j, ok := s.jobs[id]
	return j, ok
}