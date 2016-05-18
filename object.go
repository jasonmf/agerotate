// agerotate provides cleanup of timestamped objects with user-defined schedules. Objects are grouped into "buckets" based on their age. Each bucket has an age threshold and an age interval. If a bucket has an interval of two hours, roughly one object will be kept for every two hours. This allows the user to keep many objects that are recent and fewer objects that are older.
package agerotate

import (
	"time"
)

// Object is the interface objects implement to be managed by agerotate.
type Object interface {
	// Age returns the age of the object.
	Age() time.Duration
	// Delete attempts to remove the object.
	Delete() error
	// ID returns an identifier string that is intended to be unique.
	ID() string
}

// Objects is the interface for a container of Object objects.
type Objects interface {
	// ID returns an identifier string that is intended to be unique.
	ID() string
	// List retrieves all of the available Objects.
	List() ([]Object, error)
}

// ObjectsByAge implements sort.Interface to sort Objects by Age, ascending.
type ObjectsByAge struct {
	O []Object
}

func (a ObjectsByAge) Len() int           { return len(a.O) }
func (a ObjectsByAge) Swap(i, j int)      { a.O[i], a.O[j] = a.O[j], a.O[i] }
func (a ObjectsByAge) Less(i, j int) bool { return a.O[i].Age() < a.O[j].Age() }
