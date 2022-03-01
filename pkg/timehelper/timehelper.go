package timehelper

import (
	"errors"
	"time"
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

var (
	ErrInvalidHour = errors.New("TimeCond - invalid hour")
	ErrInvalidMin  = errors.New("TimeCond - invalid min")
	ErrInvalidSec  = errors.New("TimeCond - invalid sec")
	ErrNilTimeCond = errors.New("nil time cond")
)

type TimeCond struct {
	hour int
	min  int
	sec  int
}

func NewTimeCond(hour, min, sec int) (*TimeCond, error) {
	if hour < 0 || hour > 23 {
		return nil, ErrInvalidHour
	}

	if min < 0 || min > 59 {
		return nil, ErrInvalidMin
	}

	if sec < 0 || sec > 59 {
		return nil, ErrInvalidSec
	}

	return &TimeCond{hour, min, sec}, nil
}

func UntilNextTimeCond(t time.Time, timeCond *TimeCond) (time.Duration, error) {
	if timeCond == nil {
		return 0, ErrNilTimeCond
	}

	// seconds from the beginning of the day
	now := t.Hour()*60*60 + t.Minute()*60 + t.Second()

	next := timeCond.hour*60*60 + timeCond.min*60 + timeCond.sec

	var dur = next - now

	// add one more day
	if dur < 0 {
		dur += 24 * 60 * 60
	}

	return time.Duration(dur) * time.Second, nil
}

func ShortestUntilNextTimeCond(t time.Time, timeConds ...*TimeCond) (time.Duration, error) {
	var shortest time.Duration

	for _, timeCond := range timeConds {
		d, err := UntilNextTimeCond(t, timeCond)
		if err != nil {
			return 0, err
		}

		if shortest == 0 || d < shortest {
			shortest = d
		}
	}

	return shortest, nil
}
