/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package bucket

import (
	"agerotate"
)

// Cleanup sets up and invokes actual object cleanup.
func Cleanup(sortedRanges []agerotate.Range, objects agerotate.Objects) error {
	buckets := makeBuckets(sortedRanges)
	overflow, err := readObjects(objects, buckets)
	if err != nil {
		return err
	}

	if err = cleanupBuckets(buckets); err != nil {
		return err
	}

	if err = cleanupOverflow(overflow); err != nil {
		return err
	}

	return nil
}

func makeBuckets(sortedRanges []agerotate.Range) []*bucket {
	buckets := make([]*bucket, len(sortedRanges))
	for i := range sortedRanges {
		buckets[i] = newBucket(sortedRanges[i])
	}
	return buckets
}

// readObjects populates buckets by finding the bucket with the smallest age that's larger than the age of the object. If no buckets are larger than the object it's placed in an overflow list and will be deleted.
func readObjects(objects agerotate.Objects, buckets []*bucket) ([]agerotate.Object, error) {
	overflow := []agerotate.Object{}
	oList, err := objects.List()
	if err != nil {
		return nil, err
	}

	for _, o := range oList {
		found := false
		for _, b := range buckets {
			if o.Age() < b.Age() {
				b.Add(o)
				found = true
				break
			}
		}
		if !found {
			overflow = append(overflow, o)
		}
	}
	return overflow, nil
}

func cleanupBuckets(buckets []*bucket) error {
	for _, b := range buckets {
		if err := b.Cleanup(); err != nil {
			return err
		}
	}
	return nil
}

func cleanupOverflow(overflow []agerotate.Object) error {
	for _, o := range overflow {
		if err := o.Delete(); err != nil {
			return err
		}
	}
	return nil
}
