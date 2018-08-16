package competition

import "testing"

func TestMakeSeasonId(t *testing.T) {
	id := makeSeasonId(1, 2018, 2019)
	expected := "0001-2018-2019"
	if id != expected {
		t.Errorf("expected season Id to be %q, but got %q", expected, id)
	}
}

func TestNewSeason(t *testing.T) {
	s := NewSeason(1, 2018, 2019)

	if s == nil {
		t.Fatalf("got nil Season pointer")
	}

	expected := "0001-2018-2019"
	if s.ID != expected {
		t.Errorf("expected s.ID == %q, but got %q", expected, s.ID)
	}

	if s.StartYear != 2018 {
		t.Errorf("expected StartYear == %d, but got %d", 2018, s.StartYear)
	}

	if s.EndYear != 2019 {
		t.Errorf("expected EndYear == %d, but got %d", 2019, s.EndYear)
	}
}

func TestCreateAndGetSeason(t *testing.T) {
	db := setUpDB(t)
	defer tearDownDB(db)

	s := NewSeason(1, 2018, 2019)
	id, err := db.CreateSeason(s)
	if err != nil {
		t.Fatalf("unexpected error in db.CreateSeason: %s", err.Error())
	}

	s2, err := db.GetSeason(id)
	if err != nil {
		t.Fatalf("unexpected error in db.GetSeason")
	}

	if s2 == nil {
		t.Fatalf("GetSeason returned nil Season pointer")
	}

	if s2.ID != s.ID {
		t.Errorf("s2.ID == %q but expected %q", s2.ID, s.ID)
	}
	if s2.StartYear != s.StartYear {
		t.Errorf("s2.StartYear == %d but expected %d", s2.StartYear, s.StartYear)
	}
	if s2.EndYear != s.EndYear {
		t.Errorf("s2.EndYear == %d but expected %d", s2.EndYear, s.EndYear)
	}

}
