package importer

import (
	"testing"

	"github.com/drgarcia1986/street-fair/pkg/fair"
	"github.com/drgarcia1986/street-fair/pkg/logs"
)

type fakeStreetFair struct {
	createdModels []*fair.Model
}

func (f *fakeStreetFair) Create(m *fair.Model) (*fair.Model, error) {
	f.createdModels = append(f.createdModels, m)
	return m, nil
}

func TestReadFile(t *testing.T) {
	lines, err := readFile("./testdata/sample.csv")
	if err != nil {
		t.Fatalf("want <nil>; got %+v", err)
	}

	if actual := len(lines); actual != 3 {
		t.Errorf("want 3; got %d", actual)
	}

	expectedDistrict := "VILA FORMOSA"
	if actual := lines[0][6]; actual != expectedDistrict {
		t.Errorf("want %s; got %s", expectedDistrict, actual)
	}
}

func TestParseFloat(t *testing.T) {
	var testCases = []struct {
		num      string
		expected float64
	}{
		{"1", 1},
		{"20.0", 20.0},
		{"a", 0},
	}

	log, loggerFinalizer, err := logs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer loggerFinalizer()

	imp := New(log, &fakeStreetFair{createdModels: []*fair.Model{}})
	for _, tt := range testCases {
		if actual := imp.parseFloat(tt.num, "foo", "bar"); actual != tt.expected {
			t.Errorf("want %f; got %f", tt.expected, actual)
		}
	}
}

func TestRun(t *testing.T) {
	log, loggerFinalizer, err := logs.New()
	if err != nil {
		t.Fatal(err)
	}
	defer loggerFinalizer()

	fsf := &fakeStreetFair{createdModels: []*fair.Model{}}
	imp := New(log, fsf)
	if err := imp.Run("./testdata/sample.csv"); err != nil {
		t.Fatalf("want <nil>; got %+v", err)
	}

	if actual := len(fsf.createdModels); actual != 3 {
		t.Errorf("want 3; got %d", actual)
	}
	expected := "VILA FORMOSA"
	if actual := fsf.createdModels[0].District; actual != expected {
		t.Errorf("want %s; got %s", expected, actual)
	}
}
