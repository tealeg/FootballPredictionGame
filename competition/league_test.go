package competition

import "testing"

func TestCreateAndGetLeague(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	l := &League{Name: "English Premier League"}

	id, err := db.CreateLeague(l)
	if err != nil {
		t.Fatalf("unexpected error in CreateLeague: %s", err.Error())
	}

	l2, err := db.GetLeague(id)
	if err != nil {
		t.Fatalf("unexpected error GetLeague")
	}
	if l2.Name != l.Name {
		t.Errorf("Expected l2.Name == %q, but got %q", l.Name, l2.Name)
	}

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
	for i, l := range leagues {
		_, err := db.CreateLeague(l)
		if err != nil {
			t.Fatalf("[%d] Unexpected error in CreateLeague: %s", i, err.Error())
		}
	}

	results, err := db.GetAllLeagues()
	if err != nil {
		t.Fatalf("unexpected error in GetAllLeagues: %s", err.Error())
	}
	if len(results) != 4 {
		t.Fatalf("expected 4 leagues, got %d", len(results))
	}

}
