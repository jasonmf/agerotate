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
