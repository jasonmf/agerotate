package agerotate

import (
	"time"
)

type Object interface {
	// Age returns the age of the object from the supplied time as a time.Duration. If the supplied Time.IsZero(), the current time should be used.
	Age(time.Time) time.Duration
	// Delete attempts to remove the object.
	Delete() error
	// ID returns an identifier string that is intended to be unique.
	ID() string
}

type Objects interface {
	// ID returns an identifier string that is intended to be unique.
	ID() string
	// List retrieves all of the available Objects.
	List() ([]Object, error)
}

// ObjectsByAge implements sort.Interface to sort Objects by Age, ascending.
type ObjectsByAge struct {
	O   []Object
	Now time.Time
}

func (a ObjectsByAge) Len() int           { return len(a.O) }
func (a ObjectsByAge) Swap(i, j int)      { a.O[i], a.O[j] = a.O[j], a.O[i] }
func (a ObjectsByAge) Less(i, j int) bool { return a.O[i].Age(a.Now) < a.O[j].Age(a.Now) }
