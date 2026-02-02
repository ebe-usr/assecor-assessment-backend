package data

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"assecor.assessment.test/internal/validator"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Persons interface {
		Insert(persion *Person) error
		Get(id int64) (*Person, error)
		GetAll() ([]*Person, error)
		GetAllByColor(color Color) ([]*Person, error)
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Persons: &PersonModel{DB: db},
	}
}

func (m Models) LoadFromCsv(v *validator.Validator, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		v.AddError("csv", err.Error())
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	r.FieldsPerRecord = 4
	lineNumber := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		lineNumber++
		if err == nil {
			val := validator.New()
			person := parseRecord(record)
			ValidatePerson(val, &person)
			if val.Valid() {
				err = m.Persons.Insert(&person)
				if err != nil {
					v.AddError("db", err.Error())
				}
			} else {
				v.AddError("csv",
					fmt.Sprintf("line %d: invalid data record", lineNumber))
			}
		} else {
			v.AddError("csv", err.Error())
		}
	}
}

func parseRecord(r []string) Person {
	// columns: Lastname, Name, Zipcode+City, Color
	var p Person

	p.Lastname = strings.TrimSpace(r[0])
	p.Name = strings.TrimSpace(r[1])
	zc := strings.TrimSpace(r[2])
	for i, r := range zc {
		if !unicode.IsDigit(r) {
			if i > 0 {
				p.Zipcode = zc[:i]
				p.City = strings.TrimSpace(zc[i:])
				break
			}
		}
	}
	if len(p.Zipcode) == 0 {
		p.City = zc
	}
	v, _ := strconv.ParseInt(strings.TrimSpace(r[3]), 10, 32)
	p.Color = int(v)
	return p
}
