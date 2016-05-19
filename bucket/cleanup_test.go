/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package bucket

import (
	"testing"
	"time"

	"github.com/AgentZombie/agerotate"
)

type testBucketObject struct {
	age               time.Duration
	expectedBucketAge time.Duration
}

func (t testBucketObject) Age() time.Duration {
	return t.age
}

func (t testBucketObject) ID() string {
	return t.age.String() + "/" + t.expectedBucketAge.String()
}

func (t testBucketObject) Delete() error {
	return nil
}

type testBucketObjects []agerotate.Object

func (t testBucketObjects) ID() string {
	return "test object collection"
}

func (t testBucketObjects) List() ([]agerotate.Object, error) {
	return t, nil
}

func TestReadObjects(t *testing.T) {
	var irrelevantInterval time.Duration
	for _, tc := range []struct {
		id        string
		rangeAges []time.Duration
		testObjs  []agerotate.Object
	}{
		{
			id:        "All overflow",
			rangeAges: []time.Duration{10 * time.Second},
			testObjs: []agerotate.Object{
				&testBucketObject{100 * time.Second, 0},
				&testBucketObject{200 * time.Second, 0},
				&testBucketObject{300 * time.Second, 0},
				&testBucketObject{400 * time.Second, 0},
			},
		},
		{
			id:        "2 first bucket, 1 second, 2 overflow",
			rangeAges: []time.Duration{10 * time.Second, 20 * time.Second},
			testObjs: []agerotate.Object{
				&testBucketObject{0 * time.Second, 10 * time.Second},
				&testBucketObject{5 * time.Second, 10 * time.Second},
				&testBucketObject{10 * time.Second, 20 * time.Second},
				&testBucketObject{100 * time.Second, 0},
				&testBucketObject{200 * time.Second, 0},
			},
		},
	} {
		t.Logf("Testing case %q", tc.id)
		ranges := make([]agerotate.Range, len(tc.rangeAges))
		for i := range tc.rangeAges {
			ranges[i] = agerotate.Range{Age: tc.rangeAges[i], Interval: irrelevantInterval}
		}

		buckets := makeBuckets(ranges)
		objects := testBucketObjects(tc.testObjs)
		overflow, err := readObjects(objects, buckets)
		if err != nil {
			t.Fatalf("Unexpected err: %q", err)
		}

		var noDuration time.Duration
		for _, o := range overflow {
			if to, ok := o.(*testBucketObject); ok {
				if to.expectedBucketAge != noDuration {
					t.Fatalf("Unexpected object in overflow: %v", t)
				}
			} else {
				t.Fatalf("Expected testBucketObject, got %v", o)
			}
		}

		for _, b := range buckets {
			for _, o := range b.objects {
				if to, ok := o.(*testBucketObject); ok {
					if to.expectedBucketAge != b.Age() {
						t.Fatalf("Expected object %v in bucket %v, found in %v", to.age, to.expectedBucketAge, b.Age())
					}
				} else {
					t.Fatalf("Expected testBucketObject, got %v", o)
				}

			}
		}
	}
}
