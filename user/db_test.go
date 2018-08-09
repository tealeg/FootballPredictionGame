package user

import (
	"encoding/json"
	"os"
	"testing"

	bolt "github.com/coreos/bbolt"
)

func setUpAccountDB() (*AccountDB, error) {
	return NewAccountDB("test")
}

func tearDownAccountDB(adb *AccountDB) {
	adb.db.Close()
	os.Remove("test.db")
}

func TestNewAccountDB(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)
	if err != nil {
		t.Fatalf("Error creating test DB: %s", err.Error())
	}
}

func TestCreate(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}
	adb.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("account"))
		key, _ := adb.Key(acc.Name)
		v := b.Get(key)
		acc2 := &Account{}
		err := json.Unmarshal(v, acc2)
		if err != nil {
			t.Fatalf("unexpected error unmarshalling account: %s", err.Error())
			return err
		}
		if acc2.Name != acc.Name {
			t.Errorf("Expected acc2.Name == %q, but got %q", acc.Name, acc2.Name)
			return nil
		}
		return nil
	})
}

func TestCreateDuplicate(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}
	err = adb.Create(acc)
	if err == nil {
		t.Fatalf("expected error in Create, but none occurred")
	}

}

func TestGet(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)

	acc := Account{Name: "tealeg"}
	err = adb.Create(acc)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}
	acc2, err := adb.Get("tealeg")
	if err != nil {
		t.Fatalf("unexpected in Get: %s", err.Error())
	}
	if acc2.Name != acc.Name {
		t.Errorf("Expected acc2.Name == %q, but got %q", acc.Name, acc2.Name)
	}
}

func TestUpdate(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)

	acc1 := Account{Name: "tealeg", Forename: "Geoff"}
	err = adb.Create(acc1)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}

	acc2 := Account{Name: "tealeg", Forename: "Gott"}
	err = adb.Update("tealeg", acc2)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}

	acc3, err := adb.Get("tealeg")
	if err != nil {
		t.Fatalf("unexpected in Get: %s", err.Error())
	}
	if acc3.Forename == acc1.Forename {
		t.Errorf("Expected acc3.Forename != %q, but got %q", acc1.Forename, acc1.Forename)
	}
	if acc3.Forename != acc2.Forename {
		t.Errorf("Expected acc3.Forename == %q, but got %q", acc2.Forename, acc3.Forename)
	}

}

func TestAdminUserExists(t *testing.T) {
	adb, err := setUpAccountDB()
	defer tearDownAccountDB(adb)

	exists, err := adb.AdminUserExists()
	if err != nil {
		t.Fatalf("unexpected error in AdminUserExists: %s", err.Error())
	}
	if exists {
		t.Fatal("expecte AdminUserExists == false, but got true")
	}

	normalAcc := Account{Name: "nic", IsAdmin: false}
	err = adb.Create(normalAcc)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}

	exists, err = adb.AdminUserExists()
	if err != nil {
		t.Fatalf("unexpected error in AdminUserExists: %s", err.Error())
	}
	if exists {
		t.Fatal("expecte AdminUserExists == false, but got true")
	}

	adminAcc := Account{Name: "geoff", IsAdmin: true}
	err = adb.Create(adminAcc)
	if err != nil {
		t.Fatalf("unexpected error in Create: %s", err.Error())
	}

	exists, err = adb.AdminUserExists()
	if err != nil {
		t.Fatalf("unexpected error in AdminUserExists: %s", err.Error())
	}
	if !exists {
		t.Fatal("expecte AdminUserExists == true, but got false")
	}

}
