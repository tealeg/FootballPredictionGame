package competition

import (
	"bytes"
	"encoding/json"
	"fmt"

	bolt "github.com/coreos/bbolt"
)

type Season struct {
	ID        string
	StartYear uint16
	EndYear   uint16
}

// makeSeasonId returns a compound identifier that can be used to
// uniquely identify a season, and supports prefix scanning in bbolt.
func makeSeasonId(leagueId uint64, startYear, endYear uint16) string {
	return fmt.Sprintf("%04d-%d-%d", leagueId, startYear, endYear)
}

func NewSeason(leagueID uint64, startYear, endYear uint16) *Season {
	return &Season{
		ID:        makeSeasonId(leagueID, startYear, endYear),
		StartYear: startYear,
		EndYear:   endYear,
	}
}

func (db *DB) CreateSeason(s *Season) (string, error) {

	err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(seasonBN)
		buf, err := json.Marshal(s)
		if err != nil {
			return err
		}

		return b.Put([]byte(s.ID), buf)
	})
	if err != nil {
		return "", err
	}
	return s.ID, nil
}

func (db *DB) GetSeason(id string) (*Season, error) {
	s := &Season{}

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(seasonBN)
		v := b.Get([]byte(id))
		err := json.Unmarshal(v, s)
		if err != nil {
			return err
		}
		return nil
	})
	return s, err
}

func (db *DB) GetAllLeagueSeasons(leagueID uint64) ([]Season, error) {
	seasons := make([]Season, 0)

	err := db.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(seasonBN)
		c := b.Cursor()
		prefix := []byte(fmt.Sprintf("%04d", leagueID))
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			s := &Season{}
			err := json.Unmarshal(v, s)
			if err != nil {
				return err
			}
			seasons = append(seasons, *s)
		}

		return nil
	})
	return seasons, err
}
