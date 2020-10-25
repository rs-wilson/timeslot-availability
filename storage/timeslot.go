package storage

import (
	"fmt"
	"time"
)

func NewTimeslot(s time.Time, d time.Duration) *Timeslot {
	return &Timeslot{
		start:    s,
		duration: d,
	}
}

type Timeslot struct {
	start    time.Time
	duration time.Duration
}

func (t Timeslot) Key() string {
	return fmt.Sprintf("%d|%f", t.start.Unix(), t.duration.Seconds())
}

func (t Timeslot) Overlaps(o *Timeslot) bool {
	tStart, tEnd := t.getEnds()
	oStart, oEnd := o.getEnds()

	if (tStart.After(oStart) || tStart.Equal(oStart)) && (tStart.Before(oEnd) || tStart.Equal(oEnd)) {
		return true
	}

	if (tEnd.After(oStart) || tEnd.Equal(oStart)) && (tEnd.Before(oEnd) || tEnd.Equal(oEnd)) {
		return true
	}

	return false
}

func (t Timeslot) getEnds() (time.Time, time.Time) {
	return t.start, t.start.Add(t.duration)
}
