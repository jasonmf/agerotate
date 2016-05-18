// bucket captures objects below a certain age, decides which ones should be deleted, and deletes them.
package bucket

import (
	"sort"
	"time"

	"agerotate"
)

// bucket is a container for Object(s) and is intended to hold those objects younger than the Age of the Range but older than younger buckets.
type bucket struct {
	agerotate.Range
	objects []agerotate.Object
}

func newBucket(r agerotate.Range) *bucket {
	return &bucket{
		r,
		[]agerotate.Object{},
	}
}

func (b *bucket) Add(o agerotate.Object) {
	b.objects = append(b.objects, o)
}

func (b bucket) Age() time.Duration {
	return b.Range.Age
}

// Cleanup sorts the objects in the bucket by Age then deletes objects according to the Interval. The first object in the bucket is always retained. For each object thereafter, if the age of the object is less than the age of the last retained object plus Interval, the newer object is deleted. If the next object is older than the age of the last retained object plus Interval, the newer object is retained and processing continues.
func (b *bucket) Cleanup() error {
	if len(b.objects) < 2 {
		return nil
	}

	sort.Sort(agerotate.ObjectsByAge{b.objects})
	baseAge := b.objects[0].Age()
	for _, o := range b.objects[1:] {
		oAge := o.Age()
		if oAge-baseAge < b.Range.Interval {
			err := o.Delete()
			if err != nil {
				return err
			}
		} else {
			baseAge = oAge
		}
	}
	return nil
}
