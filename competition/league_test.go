package competition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAndGetLeague(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	l := &League{Name: "English Premier League"}

	id, err := db.CreateLeague(l)
	assert.NoError(t, err)
	l2, err := db.GetLeague(id)
	assert.NoError(t, err)
	assert.Equal(t, l.Name, l2.Name)
}

func TestGetAllLeagues(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	leagues := []*League{
		&League{Name: "English Premier League"},
		&League{Name: "English Championship"},
		&League{Name: "English FA League 1"},
		&League{Name: "English FA League 2"},
	}
	for _, l := range leagues {
		_, err := db.CreateLeague(l)
		assert.NoError(t, err)
	}

	results, err := db.GetAllLeagues()
	assert.NoError(t, err)
	assert.Len(t, results, 4)
}
