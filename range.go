package agerotate

import (
	"fmt"
	"time"
)

// Range identifies a set of items for rotation. Age specifies the youngest items that belong to the set. Interval defines the minimum age gap between items to keep.
type Range struct {
	Age      time.Duration
	Interval time.Duration
}

// String profiles a human-readable string for a range.
func (r Range) String() string {
	return fmt.Sprintf("For files older than %s, keep one every %s", r.Age, r.Interval)
}

// ByAge implements sort.Interface to sort Range objects by Age, ascending.
type ByAge []Range

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
