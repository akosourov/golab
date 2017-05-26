package storage

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/dhconnelly/rtreego"
)

func TestNew(t *testing.T) {
	s := New(10)
	assert.Equal(t, s.Len(), 0)
	d, err := s.Get(0)
	assert.Error(t, err)
	assert.Nil(t, d)

	driver := &Driver{
		ID: 1,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	}
	s.Set(driver)
	assert.Equal(t, s.Len(), 1)
	d, err = s.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, d, driver)

	err = s.Delete(1)
	assert.Equal(t, s.Len(), 0)
}

func TestDriverStorage_Nearest(t *testing.T) {
	s := New(10)
	s.Set(&Driver{
		ID: 1,
		LastLocation: Location{
			Lat: 42.875799,
			Lon: 74.588279,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 2,
		LastLocation: Location{
			Lat: 42.875508,
			Lon: 74.588107,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 3,
		LastLocation: Location{
			Lat: 42.876106,
			Lon: 74.588204,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 4,
		LastLocation: Location{
			Lat: 42.874942,
			Lon: 74.585908,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})
	s.Set(&Driver{
		ID: 5,
		LastLocation: Location{
			Lat: 42.875744,
			Lon: 74.584503,
		},
		Expiration: time.Now().Add(15).UnixNano(),
	})

	drivers := s.Nearest(rtreego.Point{42.876420, 74.588332}, 3)
	assert.Equal(t, len(drivers), 3)
	//assert.Equal(t, drivers[0].ID, 1)
	//assert.Equal(t, drivers[1].ID, 2)
	//assert.Equal(t, drivers[2].ID, 3)
}

func TestDriverStorage_DeleteExpired(t *testing.T) {
	s := New(10)
	s.Set(&Driver{
		ID: 1,
		LastLocation: Location{
			Lat: 1,
			Lon: 1,
		},
		Expiration: time.Now().Add(-20).UnixNano(),
	})
	s.DeleteExpired()
	assert.Equal(t, s.Len(), 0)
}

func BenchmarkNearest(b *testing.B) {
	s := New(10)
	for i := 1; i <= b.N; i++ {
		s.Set(&Driver{
			ID: i,
			LastLocation: Location{
				Lat: float64(i),
				Lon: float64(i),
			},
			Expiration: time.Now().Add(20).UnixNano(),
		})
	}
	for i := 1; i <= b.N; i++ {
		s.Nearest(rtreego.Point{1,1}, 1)
	}
}