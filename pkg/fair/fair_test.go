package fair

import (
	"testing"

	"github.com/drgarcia1986/street-fair/pkg/tests"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func fakeModel(registry string) *Model {
	return &Model{
		Longitude:      -46550164,
		Latitude:       -23558733,
		Setcens:        "355030885000091",
		Areap:          "3550308005040",
		CodDistrict:    "87",
		District:       "VILA FORMOSA",
		CodSubCityHall: "26",
		SubCityHall:    "ARICANDUVA-FORMOSA-CARRAO",
		Region5:        "Leste",
		Region8:        "Leste 1",
		Name:           "VILA FORMOSA",
		Registry:       registry,
		Address:        "RUA MARAGOJIPE",
		AddressNumber:  "S/N",
		Neighborhood:   "VL FORMOSA",
		Landmark:       "TV RUA PRETORIA",
	}
}

func testCreate(sf StreetFair, t *testing.T) {
	m := fakeModel("4041-0")
	newModel, err := sf.Create(m)
	if err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if actual := newModel.Registry; actual != m.Registry {
		t.Errorf("got %s; want %s", actual, m.Registry)
	}
}

func testCreateWithNullRegistry(sf StreetFair, t *testing.T) {
	m := fakeModel("4041-0")
	m.Registry = ""

	if _, err := sf.Create(m); err != ErrInvalidStreetFair {
		t.Error("got <nil>; want ErrInvalidStreetFair")
	}

}

func testAll(sf StreetFair, t *testing.T) {
	m1, m2 := fakeModel("4041-0"), fakeModel("4045-2")
	for _, m := range []*Model{m1, m2} {
		if _, err := sf.Create(m); err != nil {
			t.Fatalf("creating models, got %+v; want <nil>", err)
		}
	}

	models, err := sf.All(map[string]string{})
	if err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if actual := len(models); actual != 2 {
		t.Errorf("got %d; want 2", actual)
	}

	if actual := models[0].Registry; actual != m1.Registry && actual != m2.Registry {
		t.Errorf("got %s; want %s or %s", actual, m1.Registry, m2.Registry)
	}
}

func testAllWithFilter(sf StreetFair, t *testing.T) {
	expectedRegistry := "4045-2"
	m1, m2 := fakeModel("4041-5"), fakeModel(expectedRegistry)

	for _, m := range []*Model{m1, m2} {
		if _, err := sf.Create(m); err != nil {
			t.Fatalf("creating models, got %+v; want <nil>", err)
		}
	}

	models, err := sf.All(map[string]string{"registry": expectedRegistry})
	if err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if actual := len(models); actual != 1 {
		t.Errorf("got %d; want 1", actual)
	}

	if actual := models[0].Registry; actual != expectedRegistry {
		t.Errorf("got %s; want %s", actual, expectedRegistry)
	}
}

func testDelete(sf StreetFair, t *testing.T) {
	expectedRegistry := "4041-5"
	if _, err := sf.Create(fakeModel(expectedRegistry)); err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if err := sf.Delete(expectedRegistry); err != nil {
		t.Errorf("got %+v; want <nil>", err)
	}
}

func testDeleteNotFound(sf StreetFair, t *testing.T) {
	if _, err := sf.Create(fakeModel("4041-5")); err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if err := sf.Delete("9999-1"); err != ErrNotFound {
		t.Errorf("got %+v; want ErrNotFound", err)
	}
}

func testGet(sf StreetFair, t *testing.T) {
	expectedRegistry := "4041-5"
	if _, err := sf.Create(fakeModel(expectedRegistry)); err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	m, err := sf.Get(expectedRegistry)
	if err != nil {
		t.Errorf("got %+v; want <nil>", err)
	}
	if actual := m.Registry; actual != expectedRegistry {
		t.Errorf("got %s; want %s", actual, expectedRegistry)
	}
}

func testGetNotFound(sf StreetFair, t *testing.T) {
	if _, err := sf.Create(fakeModel("4041-5")); err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	if _, err := sf.Get("9999-1"); err != ErrNotFound {
		t.Errorf("got %+v; want ErrNotFound", err)
	}
}

func testUpdate(sf StreetFair, t *testing.T) {
	m, err := sf.Create(fakeModel("4041-5"))
	if err != nil {
		t.Fatalf("got %+v; want <nil>", err)
	}

	expectedName := "A New Street Fair"
	m.Name = expectedName

	if err = sf.Update(m); err != nil {
		t.Errorf("got %+v; want <nil>", err)
	}

	nm, err := sf.Get(m.Registry)
	if err != nil {
		t.Fatal(err)
	}
	if actual := nm.Name; actual != expectedName {
		t.Errorf("got %s; want %s", actual, expectedName)
	}
}

func testUpdateNotFound(sf StreetFair, t *testing.T) {
	if err := sf.Update(fakeModel("0000-1")); err != ErrNotFound {
		t.Errorf("got %+v; want ErrNotFound", err)
	}
}

func testSetup(db *gorm.DB) error {
	if r := db.Where("1 = 1").Delete(&Model{}); r.Error != nil {
		return r.Error
	}
	return nil
}

func TestStreetFair(t *testing.T) {
	var unitTests = []struct {
		title string
		test  func(sf StreetFair, t *testing.T)
	}{
		{"Create", testCreate},
		{"CreateWithNullRegistry", testCreateWithNullRegistry},
		{"All", testAll},
		{"AllWithFilter", testAllWithFilter},
		{"Delete", testDelete},
		{"DeleteNotFound", testDeleteNotFound},
		{"Get", testDelete},
		{"GetNotFound", testGetNotFound},
		{"Update", testUpdate},
		{"UpdateNotFound", testUpdateNotFound},
	}

	for _, ut := range unitTests {
		db, err := tests.NewDB()
		if err != nil {
			t.Fatal(err)
		}

		d, err := New(db, logrus.New())
		if err != nil {
			t.Fatal(err)
		}

		if err := testSetup(db); err != nil {
			t.Fatal(err)
		}

		t.Run(ut.title, func(t *testing.T) {
			ut.test(d, t)
		})
	}
}
