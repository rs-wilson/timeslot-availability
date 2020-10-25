package storage

import (
	"errors"
	"time"
)

func NewInMemoryTimeslotStore() *InMemoryTimeslotStore {
	return &InMemoryTimeslotStore{}
}

type InMemoryTimeslotStore struct {
}

func (ts *InMemoryTimeslotStore) IsAvailable(time.Time, time.Duration) bool {
	return false
}

func (ts *InMemoryTimeslotStore) Reserve(time.Time, time.Duration) error {
	return errors.New("not implemented")
}

func (ts *InMemoryTimeslotStore) Delete(time.Time, time.Duration) (bool, error) {
	return false, errors.New("not implemented")
}
