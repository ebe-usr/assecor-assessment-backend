package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"assecor.assessment.test/internal/validator"
)

type Person struct {
	ID       int64  `json:"id"` // Unique integer ID for the person
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Zipcode  string `json:"zipcode"`
	City     string `json:"city"`
	Color    int    `json:"color"` // Person's favorite color
}

func ValidatePerson(v *validator.Validator, person *Person) {
	v.Check(person.Name != "", "name", "must be provided")
	v.Check(len(person.Name) <= 250, "name", "must not be more than 250 bytes long")
	v.Check(person.Lastname != "", "lastname", "must be provided")
	v.Check(len(person.Lastname) <= 250, "name", "must not be more than 250 bytes long")
	v.Check(person.Zipcode != "", "zipcode", "must be provided")
	v.Check(validator.ZipCodeRX.MatchString(person.Zipcode), "zipcode", "invalid zip code")
	v.Check(person.City != "", "city", "must be provided")
	v.Check(len(person.City) <= 250, "city", "must not be more than 250 bytes long")
	v.Check(person.Color >= 1 && person.Color < int(LastColorIndex), "color", "id out of range")
}

type Color int

const (
	Blue Color = iota + 1
	Green
	Purple
	Red
	Yellow
	Turquoise
	White
	LastColorIndex
)

var colorName = map[Color]string{
	Blue:      "blau",
	Green:     "grün",
	Purple:    "violett",
	Red:       "rot",
	Yellow:    "gelb",
	Turquoise: "türkis",
	White:     "weiß",
}

func (c Color) String() string {
	return colorName[c]
}

type PersonModel struct {
	DB *sql.DB
}

func (m *PersonModel) Insert(p *Person) error {
	query := `
		INSERT INTO persons (name, lastname, zipcode, city, color)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	args := []interface{}{p.Name, p.Lastname, p.Zipcode, p.City, p.Color}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&p.ID)
}

func (m *PersonModel) Get(id int64) (*Person, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id, name, lastname, zipcode, city, color
		FROM persons
		WHERE id = $1`
	var p Person

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Lastname, &p.Zipcode, &p.City, &p.Color)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &p, nil
}

func (m *PersonModel) GetAll() ([]*Person, error) {
	query := `
		SELECT id, name, lastname, zipcode, city, color
		FROM persons
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	persons := []*Person{}
	for rows.Next() {
		var person Person
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Lastname,
			&person.Zipcode,
			&person.City,
			&person.Color,
		)
		if err != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return persons, nil
}

func (m *PersonModel) GetAllByColor(color Color) ([]*Person, error) {
	query := `
		SELECT id, name, lastname, zipcode, city, color
		FROM persons
		WHERE (color = $1)
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, int(color))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Initialize an empty slice to hold the movie data.
	persons := []*Person{}
	for rows.Next() {
		var person Person
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Lastname,
			&person.Zipcode,
			&person.City,
			&person.Color,
		)
		if err != nil {
			return nil, err
		}
		persons = append(persons, &person)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(persons) == 0 {
		return nil, ErrRecordNotFound
	}
	return persons, nil
}
