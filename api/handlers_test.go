package main

import (
	"net/http"
	"strings"
	"testing"

	"assecor.assessment.test/internal/data"
)

func TestPing(t *testing.T) {
	app := newTestApp(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/healthcheck")

	if code != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, code)
	}

	var input struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	readJSON(t, body, &input)
	if input.Status != "available" {
		t.Errorf("want Status %s; got %s", "available", input.Status)
	}
	if input.Version != "1.0.0" {
		t.Errorf("want Version %s; got %s", "1.0.0", input.Version)
	}
}

func TestShowPerson(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	p := &data.Person{
		Name:     "Max",
		Lastname: "Mustermann",
		Zipcode:  "45555",
		City:     "Musterstadt",
		Color:    int(data.Blue),
	}
	err := app.models.Persons.Insert(p)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody *data.Person
	}{
		{"Valid ID", "/persons/1", http.StatusOK, p},
		{"Non-existent ID", "/persons/2", http.StatusNotFound, nil},
		{"Negative ID", "/persons/-1", http.StatusUnprocessableEntity, nil},
		{"Decimal ID", "/persons/1.23", http.StatusUnprocessableEntity, nil},
		{"String ID", "/persons/foo", http.StatusUnprocessableEntity, nil},
		{"Trailing slash", "/persons/1/", http.StatusNotFound, nil},
	}

	for _, tt := range tests {
		// rebind tt into this lexical scope to avoid concurrency bug from running
		// sub-tests
		// tt := tt

		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Fatalf("want %d; got %d", tt.wantCode, code)
			}
			if tt.wantBody != nil {
				var input struct {
					ID       int64  `json:"id"`
					Name     string `json:"name"`
					Lastname string `json:"lastname"`
					Zipcode  string `json:"zipcode"`
					City     string `json:"city"`
					Color    string `json:"color"`
				}
				readJSON(t, body, &input)

				if input.ID != tt.wantBody.ID {
					t.Errorf("want ID %d; got %d", tt.wantBody.ID, input.ID)
				}
				if input.Name != tt.wantBody.Name {
					t.Errorf("want Name %s; got %s", tt.wantBody.Name, input.Name)
				}
				if input.Lastname != tt.wantBody.Lastname {
					t.Errorf("want Lastname %s; got %s", tt.wantBody.Lastname, input.Lastname)
				}
				if input.Zipcode != tt.wantBody.Zipcode {
					t.Errorf("want Zipcode %s; got %s", tt.wantBody.Zipcode, input.Zipcode)
				}
				if input.City != tt.wantBody.City {
					t.Errorf("want City %s; got %s", tt.wantBody.City, input.City)
				}
				if input.Color != data.Color(tt.wantBody.Color).String() {
					t.Errorf("want Color %s; got %s", data.Color(tt.wantBody.Color).String(), input.Color)
				}
			}
		})
	}
}

func TestShowPersons(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	t.Run("No persons", func(t *testing.T) {
		code, _, body := ts.get(t, "/persons")

		if code != http.StatusOK {
			t.Fatalf("want %d; got %d", http.StatusOK, code)
		}
		js := strings.TrimSpace(string(body))
		if js != "[]" {
			t.Errorf("want %s; got %s", "[]", js)
		}
	})

	testPersons := []*data.Person{
		{
			Name:     "Hans",
			Lastname: "Müller",
			Zipcode:  "67742",
			City:     "Lauterecken",
			Color:    int(data.Blue),
		},
		{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    int(data.Green),
		},
		{
			Name:     "Johnny",
			Lastname: "Johnson",
			Zipcode:  "88888",
			City:     "made up",
			Color:    int(data.Purple),
		},
		{
			Name:     "Milly",
			Lastname: "Millenium",
			Zipcode:  "77777",
			City:     "made up too",
			Color:    int(data.Red),
		},
	}
	for _, p := range testPersons {
		err := app.models.Persons.Insert(p)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []*data.Person
	}{
		{"Valid", "/persons", http.StatusOK, testPersons},
		{"Trailing slash", "/persons/", http.StatusUnprocessableEntity, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Fatalf("want %d; got %d", tt.wantCode, code)
			}
			if tt.wantBody != nil {
				type person struct {
					ID       int64  `json:"id"`
					Name     string `json:"name"`
					Lastname string `json:"lastname"`
					Zipcode  string `json:"zipcode"`
					City     string `json:"city"`
					Color    string `json:"color"`
				}
				var input []person
				readJSON(t, body, &input)

				for len(input) == 0 {
					t.Fatal("unexpected empty array")
				}
				for i, p := range input {
					if p.ID != tt.wantBody[i].ID {
						t.Errorf("Item[%d] want ID %d; got %d", i, tt.wantBody[i].ID, p.ID)
					}
					if p.Name != tt.wantBody[i].Name {
						t.Errorf("Item[%d] want Name %s; got %s", i, tt.wantBody[i].Name, p.Name)
					}
					if p.Lastname != tt.wantBody[i].Lastname {
						t.Errorf("Item[%d] want Lastname %s; got %s", i, tt.wantBody[i].Lastname, p.Lastname)
					}
					if p.Zipcode != tt.wantBody[i].Zipcode {
						t.Errorf("Item[%d] want Zipcode %s; got %s", i, tt.wantBody[i].Zipcode, p.Zipcode)
					}
					if p.City != tt.wantBody[i].City {
						t.Errorf("Item[%d] want City %s; got %s", i, tt.wantBody[i].City, p.City)
					}
					if p.Color != data.Color(tt.wantBody[i].Color).String() {
						t.Errorf("Item[%d] want Color %s; got %s", i, data.Color(tt.wantBody[i].Color).String(), p.Color)
					}
				}
			}
		})
	}
}

func TestFilterByPersonsFavoriteColor(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	testPersons := []*data.Person{
		{
			Name:     "Hans",
			Lastname: "Müller",
			Zipcode:  "67742",
			City:     "Lauterecken",
			Color:    int(data.Blue),
		},
		{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    int(data.Green),
		},
		{
			Name:     "Johnny",
			Lastname: "Johnson",
			Zipcode:  "88888",
			City:     "made up",
			Color:    int(data.Blue),
		},
		{
			Name:     "Milly",
			Lastname: "Millenium",
			Zipcode:  "77777",
			City:     "made up too",
			Color:    int(data.Green),
		},
		{
			Name:     "Jonas",
			Lastname: "Müller",
			Zipcode:  "32323",
			City:     "Hansstadt",
			Color:    int(data.Yellow),
		},
		{
			Name:     "Tastatur",
			Lastname: "Fujitsu",
			Zipcode:  "42342",
			City:     "Japan",
			Color:    int(data.Blue),
		},
	}
	for _, p := range testPersons {
		err := app.models.Persons.Insert(p)
		if err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []*data.Person
	}{
		{"Valid ColorID(blue)", "/persons/color/1", http.StatusOK,
			[]*data.Person{testPersons[0], testPersons[2], testPersons[5]}},
		{"Valid ColorID(green)", "/persons/color/2", http.StatusOK,
			[]*data.Person{testPersons[1], testPersons[3]}},
		{"Valid ColorID(yellow)", "/persons/color/5", http.StatusOK,
			[]*data.Person{testPersons[4]}},
		{"Non-existent ColorID", "/persons/color/4", http.StatusNotFound, nil},
		{"Negative ColorID", "/persons/color/-1", http.StatusUnprocessableEntity, nil},
		{"Trailing slash", "/persons/color/", http.StatusUnprocessableEntity, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			if code != tt.wantCode {
				t.Fatalf("want %d; got %d", tt.wantCode, code)
			}
			if tt.wantBody != nil {
				type person struct {
					ID       int64  `json:"id"`
					Name     string `json:"name"`
					Lastname string `json:"lastname"`
					Zipcode  string `json:"zipcode"`
					City     string `json:"city"`
					Color    string `json:"color"`
				}
				var input []person
				readJSON(t, body, &input)

				for i, p := range input {
					if p.ID != tt.wantBody[i].ID {
						t.Errorf("Item[%d] want ID %d; got %d", i, tt.wantBody[i].ID, p.ID)
					}
				}
			}
		})
	}
}

func TestCreatePersons(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	type person struct {
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		Zipcode  string `json:"zipcode"`
		City     string `json:"city"`
		Color    int    `json:"color"`
	}
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody person
	}{
		{"Valid person(ID=1)", "/persons", http.StatusCreated, person{
			Name:     "Hans",
			Lastname: "Müller",
			Zipcode:  "67742",
			City:     "Lauterecken",
			Color:    int(data.Blue)},
		},
		{"Valid person(ID=2)", "/persons", http.StatusCreated, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    int(data.Green)},
		},
		{"Empty name", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    int(data.Green)},
		},
		{"Empty lastname", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    int(data.Green)},
		},
		{"Empty zipcode", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "",
			City:     "Stralsund",
			Color:    int(data.Green)},
		},
		{"Zipcode xxx zipcode", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "xxx",
			City:     "Stralsund",
			Color:    int(data.Green)},
		},
		{"Empty city", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "",
			Color:    int(data.Green)},
		},
		{"Negative colorID", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    -1},
		},
		{"Out of range colorID", "/persons", http.StatusUnprocessableEntity, person{
			Name:     "Peter",
			Lastname: "Petersen",
			Zipcode:  "18439",
			City:     "Stralsund",
			Color:    100},
		},
		{"Trailing slash", "/persons/", http.StatusTemporaryRedirect, person{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.post(t, tt.urlPath, writeJSON(t, tt.wantBody))

			if code != tt.wantCode {
				t.Fatalf("want %d; got %d", tt.wantCode, code)
			}
			if code == http.StatusOK {
				var resp data.Person
				readJSON(t, body, &resp)

				if resp.Name != tt.wantBody.Name {
					t.Errorf("ID=%d want Name %s; got %s", resp.ID, tt.wantBody.Name, resp.Name)
				}
				if resp.Lastname != tt.wantBody.Lastname {
					t.Errorf("ID=%d want Lastname %s; got %s", resp.ID, tt.wantBody.Lastname, resp.Lastname)
				}
				if resp.Zipcode != tt.wantBody.Zipcode {
					t.Errorf("ID=%d want Zipcode %s; got %s", resp.ID, tt.wantBody.Zipcode, resp.Zipcode)
				}
				if resp.City != tt.wantBody.City {
					t.Errorf("ID=%d want City %s; got %s", resp.ID, tt.wantBody.City, resp.City)
				}
				if resp.Color != tt.wantBody.Color {
					t.Errorf("ID=%d want Color %d; got %d", resp.ID, tt.wantBody.Color, resp.Color)
				}
			}
		})
	}
}
