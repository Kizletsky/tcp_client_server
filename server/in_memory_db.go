package main

import "sync"

type inMemoryDb struct {
	mu     sync.Mutex
	models []model
}

func newInMemoryDb() *inMemoryDb {
	return &inMemoryDb{}
}

func (db *inMemoryDb) put(m model) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, pm := range db.models {
		if pm.key == m.key && pm.timestamp == m.timestamp {
			db.models[i] = m
			return
		}
	}

	db.models = append(db.models, m)

	return
}

func (db *inMemoryDb) get(key string) ([]model, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	var res []model

	for _, m := range db.models {
		if m.key == key {
			res = append(res, m)
		}
	}

	return res, nil
}

func (db *inMemoryDb) getAll() ([]model, error) {
	return db.models, nil
}
