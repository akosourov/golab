package storage

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	s := New()
	if lenS := s.Len(); lenS != 0 {
		t.Errorf("Storage must be empty. Got: '%s'", lenS)
	}
	if _, err := s.Get(1); err == nil {
		t.Error("Storage must return error due to empty")
	}

	driver := &Driver{
		ID: 1,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
	}
	s.Set(1, driver)
	if lenS := s.Len(); lenS != 1 {
		t.Errorf("Storage length must equal 1. Got: '%s'", lenS)
	}
	d, err := s.Get(1)
	if err != nil {
		t.Errorf("Must return ok. Got error :'%s'", err)
	}
	if d != driver {
		t.Errorf("Should return '%v'. Got: '%v'", driver, d)
	}

	err = s.Delete(1)
	if err != nil {
		t.Errorf("Must return ok. Got error :'%s'", err)
	}
	if lenS := s.Len(); lenS != 0 {
		t.Errorf("Storage must be empty. Got: '%s'", lenS)
	}
}

func TestNearest(t *testing.T) {
	s := New()
	s.Set(123, &Driver{
		ID: 123,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
	})
	s.Set(666, &Driver{
		ID: 666,
		LastLocation: Location{
			Lat: 42.875799,
			Lon: 74.588279,
		},
	})
	drivers := s.Nearest(1, 42.876420, 74.588332)
	assert.Equal(t, len(drivers), 1)
}

func BenchmarkNearest(b *testing.B) {
	s := New()
	for i := 1; i <= b.N; i++ {
		s.Set(i, &Driver{
			ID: i,
			LastLocation: Location{
				Lat: float64(i),
				Lon: float64(i),
			},
		})
	}
	for i := 1; i <= b.N; i++ {
		s.Nearest(1, 123, 123)
	}
}