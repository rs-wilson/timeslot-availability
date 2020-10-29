package storage

import (
	"log"
	"sync"
	"time"
)

func NewInMemoryTimeslotStore() *InMemoryTimeslotStore {
	return &InMemoryTimeslotStore{
		timeslots: make(map[string]*Timeslot),
	}
}

type InMemoryTimeslotStore struct {
	timeslots map[string]*Timeslot
	rw        sync.RWMutex //ensure our storage can be accessed concurrently safely
}

// This public function locks our storage
func (i *InMemoryTimeslotStore) IsAvailable(slot time.Time, dur time.Duration) bool {
	i.rw.RLock()
	defer i.rw.RUnlock()
	return i.isAvailable(slot, dur)
}

// This private function can be used in any internal function, because it doesn't lock
func (i *InMemoryTimeslotStore) isAvailable(slot time.Time, dur time.Duration) bool {
	ts := NewTimeslot(slot, dur)
	log.Printf("checking availability for timeslot T=%d D=%f", slot.Unix(), dur.Seconds())
	for _, used := range i.timeslots {
		log.Printf("\tagainst T=%d D=%f", used.start.Unix(), used.duration.Seconds())
		if ts.Overlaps(used) {
			return false
		}
	}
	return true
}

// Reserve checks availability and returns false if called on an unavailable slot
func (i *InMemoryTimeslotStore) Reserve(slot time.Time, dur time.Duration) (bool, error) {
	i.rw.Lock()
	defer i.rw.Unlock()

	//Check if the slot is available
	available := i.isAvailable(slot, dur)
	if !available {
		log.Printf("timeslot is not available T=%d D=%f", slot.Unix(), dur.Seconds())
		return false, nil
	}

	// If it is, reserve it
	ts := NewTimeslot(slot, dur)
	i.timeslots[ts.Key()] = ts
	log.Printf("reserved timeslot T=%d D=%f", slot.Unix(), dur.Seconds())
	return true, nil
}

func (i *InMemoryTimeslotStore) Delete(slot time.Time, dur time.Duration) (bool, error) {
	i.rw.Lock()
	defer i.rw.Unlock()

	ts := NewTimeslot(slot, dur)
	_, ok := i.timeslots[ts.Key()]
	if !ok {
		log.Printf("timeslot does not exist T=%d D=%f", slot.Unix(), dur.Seconds())
		return false, nil
	}
	delete(i.timeslots, ts.Key())

	log.Printf("freed timeslot T=%d D=%f", slot.Unix(), dur.Seconds())

	return true, nil
}
