package bucket

import (
	"sort"
	"time"

	"agerotate"
)

type bucket struct {
	agerotate.Range
	objects []agerotate.Object
}

func newBucket(r agerotate.Range) bucket {
	return bucket{
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

func (b *bucket) Cleanup(now time.Time) error {
	if len(b.objects) < 2 {
		return nil
	}

	sort.Sort(agerotate.ObjectsByAge{b.objects, now})
	baseAge := b.objects[0].Age(now)
	for _, o := range b.objects[1:] {
		oAge := o.Age(now)
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