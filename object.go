/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

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
