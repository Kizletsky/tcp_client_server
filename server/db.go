package main

import "time"

type model struct {
	key       string
	value     float64
	timestamp time.Time
}

type db interface {
	put(model) error
	get(string) ([]model, error)
	getAll() ([]model, error)
}
