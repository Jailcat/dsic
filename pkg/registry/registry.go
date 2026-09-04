package registry

import (
	"encoding/json"
	"sync"
	"time"
)

type Record struct {
	Domain    string    `json:"domain"`
	Address   string    `json:"address"`
	Port      string    `json:"port"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Registry struct {
	mu      sync.RWMutex
	records map[string]Record
}

func New() *Registry {
	return &Registry{
		records: make(map[string]Record),
	}
}

func (r *Registry) Register(domain, address, port string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[domain] = Record{
		Domain:    domain,
		Address:   address,
		Port:      port,
		UpdatedAt: time.Now(),
	}
}

func (r *Registry) Lookup(domain string) (Record, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rec, ok := r.records[domain]
	return rec, ok
}

func (r *Registry) All() []Record {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Record, 0, len(r.records))
	for _, rec := range r.records {
		out = append(out, rec)
	}
	return out
}

func (r *Registry) FromJSON(data []byte) error {
	var records []Record
	if err := json.Unmarshal(data, &records); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, rec := range records {
		if existing, ok := r.records[rec.Domain]; !ok || rec.UpdatedAt.Before(existing.UpdatedAt) {
			r.records[rec.Domain] = rec
		}
	}
	return nil
}

func (r *Registry) ToJSON() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	records := make([]Record, 0, len(r.records))
	for _, rec := range r.records {
		records = append(records, rec)
	}
	return json.Marshal(records)
}
