package mock

import (
	"sort"

	"assecor.assessment.test/internal/data"
)

type MockPersonModel struct {
	seqID int64
	db    map[int64]*data.Person
}

func NewTestModels() data.Models {
	return data.Models{
		Persons: &MockPersonModel{
			seqID: 0,
			db:    make(map[int64]*data.Person)},
	}
}

func (m *MockPersonModel) Insert(person *data.Person) error {
	m.seqID++
	person.ID = m.seqID
	m.db[m.seqID] = &data.Person{
		ID:       m.seqID,
		Name:     person.Name,
		Lastname: person.Lastname,
		Zipcode:  person.Zipcode,
		City:     person.City,
		Color:    person.Color,
	}
	return nil
}

func (m *MockPersonModel) Get(id int64) (*data.Person, error) {
	p, ok := m.db[id]
	if ok {
		return p, nil
	}
	return nil, data.ErrRecordNotFound
}

func (m *MockPersonModel) GetAll() ([]*data.Person, error) {
	persons := make([]*data.Person, len(m.db))
	i := 0
	for _, p := range m.db {
		persons[i] = p
		i++
	}
	sort.Slice(persons[:], func(i, j int) bool {
		return persons[i].ID < persons[j].ID
	})
	return persons, nil
}

func (m *MockPersonModel) GetAllByColor(color data.Color) ([]*data.Person, error) {
	var persons []*data.Person
	for _, p := range m.db {
		if p.Color == int(color) {
			persons = append(persons, p)
		}
	}
	if len(persons) == 0 {
		return nil, data.ErrRecordNotFound
	}
	sort.Slice(persons[:], func(i, j int) bool {
		return persons[i].ID < persons[j].ID
	})
	return persons, nil
}
