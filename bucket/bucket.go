/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// bucket captures objects below a certain age, decides which ones should be deleted, and deletes them.
package bucket

import (
	"sort"
	"time"

	"github.com/AgentZombie/agerotate"
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
