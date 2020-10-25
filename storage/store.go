package storage

import (
	"log"
	"time"
)

func NewInMemoryTimeslotStore() *InMemoryTimeslotStore {
	return &InMemoryTimeslotStore{
		timeslots: make(map[string]*Timeslot),
	}
}

type InMemoryTimeslotStore struct {
	timeslots map[string]*Timeslot
}

func (i *InMemoryTimeslotStore) IsAvailable(slot time.Time, dur time.Duration) bool {
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

func (i *InMemoryTimeslotStore) Reserve(slot time.Time, dur time.Duration) error {
	ts := NewTimeslot(slot, dur)
	i.timeslots[ts.Key()] = ts
	log.Printf("reserved timeslot T=%d D=%f", slot.Unix(), dur.Seconds())
	return nil
}

func (i *InMemoryTimeslotStore) Delete(slot time.Time, dur time.Duration) (bool, error) {
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
