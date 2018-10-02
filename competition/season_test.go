package competition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSeasonId(t *testing.T) {
	id := makeSeasonId(1, 2018, 2019)
	expected := "0001-2018-2019"
	assert.Equal(t, expected, id)
}

func TestNewSeason(t *testing.T) {
	s := NewSeason(1, 2018, 2019)
	assert.NotNil(t, s)

	expected := "0001-2018-2019"
	assert.Equal(t, expected, s.ID)
	assert.Equal(t, uint16(2018), s.StartYear)
	assert.Equal(t, uint16(2019), s.EndYear)
}

func TestCreateAndGetSeason(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	s := NewSeason(1, 2018, 2019)
	id, err := db.CreateSeason(s)
	assert.NoError(t, err)
	s2, err := db.GetSeason(id)
	assert.NoError(t, err)
	assert.NotNil(t, s2)
	assert.Equal(t, s.ID, s2.ID)
	assert.Equal(t, s.StartYear, s2.StartYear)
	assert.Equal(t, s.EndYear, s2.EndYear)
}

func TestGetAllLeagueSesasons(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	l := &League{Name: "English Premier League"}
	_, err := db.CreateLeague(l)
	assert.NoError(t, err)

	l2 := &League{Name: "English Championship"}
	_, err = db.CreateLeague(l2)
	assert.NoError(t, err)

	for start := uint16(2016); start < 2020; start++ {
		s := NewSeason(l.ID, start, start+1)
		_, err := db.CreateSeason(s)
		assert.NoError(t, err)

		s2 := NewSeason(l2.ID, start, start+1)
		_, err = db.CreateSeason(s2)
		assert.NoError(t, err)
	}

	// We should only get EPL seasons, not championship
	seasons, err := db.GetAllLeagueSeasons(l.ID)
	assert.NoError(t, err)
	assert.Len(t, seasons, 4)
}
