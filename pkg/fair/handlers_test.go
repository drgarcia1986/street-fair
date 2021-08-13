package fair

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

type fakeStreetFair struct {
	createReturn   *Model
	createErr      error
	allReturn      []Model
	allErr         error
	districtFilter string
	deleteErr      error
	updateErr      error
	getReturn      *Model
	getErr         error
}

func (f *fakeStreetFair) Create(model *Model) (*Model, error) {
	return f.createReturn, f.createErr
}

func (f *fakeStreetFair) All(filters map[string]string) ([]Model, error) {
	f.districtFilter = filters["district"]
	return f.allReturn, f.allErr
}

func (f *fakeStreetFair) Delete(registry string) error {
	return f.deleteErr
}

func (f *fakeStreetFair) Update(model *Model) error {
	return f.updateErr
}

func (f *fakeStreetFair) Get(registry string) (*Model, error) {
	return f.getReturn, f.getErr
}

func TestHandlerGet(t *testing.T) {
	var testCases = []struct {
		title          string
		model          *Model
		methodError    error
		expectedStatus int
	}{
		{
			"Everything Ok",
			fakeModel("4041-5"),
			nil,
			http.StatusOK,
		},
		{
			"NotFound",
			nil,
			ErrNotFound,
			http.StatusNotFound,
		},
		{
			"Unknow Error",
			nil,
			errors.New("some error"),
			http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			fsf := &fakeStreetFair{
				getReturn: tt.model,
				getErr:    tt.methodError,
			}
			api := NewHTTPService(fsf)
			req, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.Get)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got %d; want %d", status, tt.expectedStatus)
			}

			if tt.model == nil {
				return
			}

			var model Model
			if err := json.NewDecoder(rr.Body).Decode(&model); err != nil {
				t.Fatal(err)
			}
			if model.Registry != tt.model.Registry {
				t.Errorf("got %+v; want %+v", model, tt.model)
			}
		})
	}
}

func TestHandlerAll(t *testing.T) {
	var testCases = []struct {
		title          string
		models         []Model
		methodError    error
		expectedStatus int
		url            string
		expectedFilter string
	}{
		{
			"Everything Ok - One Record",
			[]Model{*fakeModel("4041-5")},
			nil,
			http.StatusOK,
			"",
			"",
		},
		{
			"Everything Ok - Three Record",
			[]Model{*fakeModel("4041-5"), *fakeModel("4045-2"), *fakeModel("3048-1")},
			nil,
			http.StatusOK,
			"",
			"",
		},
		{
			"Everything Ok - Filter",
			[]Model{*fakeModel("4041-5")},
			nil,
			http.StatusOK,
			"?district=98",
			"98",
		},
		{
			"Everything Ok - Zero Records",
			[]Model{},
			nil,
			http.StatusOK,
			"",
			"",
		},
		{
			"Some Error",
			nil,
			errors.New("Some error"),
			http.StatusInternalServerError,
			"",
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			fsf := &fakeStreetFair{
				allReturn: tt.models,
				allErr:    tt.methodError,
			}
			api := NewHTTPService(fsf)
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.All)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got %d; want %d", status, tt.expectedStatus)
			}

			if tt.models == nil {
				return
			}
			models := make([]Model, 0)
			if err := json.NewDecoder(rr.Body).Decode(&models); err != nil {
				t.Fatal(err)
			}
			if actual := len(models); actual != len(tt.models) {
				t.Errorf("got %d; want %d", actual, len(tt.models))
			}
			if actual := fsf.districtFilter; actual != tt.expectedFilter {
				t.Errorf("got %s; want %s", actual, tt.expectedFilter)
			}
		})
	}
}

func TestHandlerCreate(t *testing.T) {
	var testCases = []struct {
		title            string
		model            *Model
		payload          string
		methodError      error
		expectedStatus   int
		expectedRegistry string
	}{
		{
			"Everything Ok",
			fakeModel("5171-3"),
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			nil,
			http.StatusCreated,
			"5171-3",
		},
		{
			"Some Error",
			nil,
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			errors.New("Some Error"),
			http.StatusInternalServerError,
			"",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			fsf := &fakeStreetFair{
				createReturn: tt.model,
				createErr:    tt.methodError,
			}
			api := NewHTTPService(fsf)
			jsonStr := []byte(tt.payload)
			req, err := http.NewRequest("POST", "", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.Create)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got %d; want %d", status, tt.expectedStatus)
			}
			if tt.expectedRegistry == "" {
				return
			}
			payload := rr.Body.String()
			if !strings.Contains(payload, fmt.Sprintf(`"registry":"%s"`, tt.expectedRegistry)) {
				t.Errorf("payload doesn't have the registry; got %s", payload)
			}

		})
	}
}

func TestHandlerDelete(t *testing.T) {
	var testCases = []struct {
		title          string
		methodError    error
		expectedStatus int
	}{
		{
			"Everything OK",
			nil,
			http.StatusNoContent,
		},
		{
			"Not Found",
			ErrNotFound,
			http.StatusNotFound,
		},
		{
			"Server Error",
			errors.New("some error"),
			http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			api := NewHTTPService(
				&fakeStreetFair{deleteErr: tt.methodError},
			)

			req, err := http.NewRequest("DELETE", "/REGISTRY", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.Delete)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got %d want %d", status, tt.expectedStatus)
			}
		})
	}
}

func TestHandlerUpdate(t *testing.T) {
	var testCases = []struct {
		title          string
		methodError    error
		payload        string
		urlRegistry    string
		expectedStatus int
	}{
		{
			"Everything OK",
			nil,
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			"5171-3",
			http.StatusOK,
		},
		{
			"Not Found",
			ErrNotFound,
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			"5171-3",
			http.StatusNotFound,
		},
		{
			"Bad Request (wrong registry on payload)",
			nil,
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			"0000-1",
			http.StatusBadRequest,
		},
		{
			"Server Error",
			errors.New("some error"),
			`{"longitude":-46450424,"latitude":-23602582,"setcens":"355030833000022","areap":"3550308005274","cod_district":"32","district":"IGUATEMI","cod_sub_city_hall":"30","sub_city_hall":"SAO MATEUS","region_5":"Leste","region_8":"Leste 2","name":"JD.BOA ESPERANCA","registry":"5171-3","address":"RUA IGUPIARA","address_number":"S/N","neighborhood":"JD BOA ESPERANCA","landmark":""}`,
			"5171-3",
			http.StatusInternalServerError,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.title, func(t *testing.T) {
			api := NewHTTPService(
				&fakeStreetFair{updateErr: tt.methodError},
			)

			jsonStr := []byte(tt.payload)
			req, err := http.NewRequest("PUT", "/REGISTRY", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			req = mux.SetURLVars(req, map[string]string{"registry": tt.urlRegistry})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(api.Update)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("got %d want %d (%s)", status, tt.expectedStatus, rr.Body.String())
			}
		})
	}
}
