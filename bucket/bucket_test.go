package bucket

import (
	"testing"
	"time"

	"agerotate"
)

type testObject struct {
	age     time.Duration
	deleted bool
}

func (t *testObject) Age() time.Duration {
	return t.age
}

func (t *testObject) Delete() error {
	t.deleted = true
	return nil
}

func (t *testObject) ID() string {
	return t.age.String()
}

func TestCleanup(t *testing.T) {
	var irrelevantDuration time.Duration
	deleted := true

	for _, tc := range []struct {
		id       string
		interval time.Duration
		objects  []time.Duration
		expected []bool
	}{
		{
			id:       "All keepers",
			interval: 0 * time.Second,
			objects:  []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second},
			expected: []bool{!deleted, !deleted, !deleted, !deleted},
		},
		{
			id:       "Keep first",
			interval: 5 * time.Second,
			objects:  []time.Duration{1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second},
			expected: []bool{!deleted, deleted, deleted, deleted},
		},
		{
			id:       "30s interval, 30s age splay",
			interval: 30 * time.Second,
			objects:  []time.Duration{0 * time.Second, 30 * time.Second, 60 * time.Second, 90 * time.Second},
			expected: []bool{!deleted, !deleted, !deleted, !deleted},
		},
		{
			id:       "29s interval, 30s age splay",
			interval: 29 * time.Second,
			objects:  []time.Duration{0 * time.Second, 30 * time.Second, 60 * time.Second, 90 * time.Second},
			expected: []bool{!deleted, !deleted, !deleted, !deleted},
		},
		{
			id:       "31s interval, 30s age splay",
			interval: 31 * time.Second,
			objects:  []time.Duration{0 * time.Second, 30 * time.Second, 60 * time.Second, 90 * time.Second},
			expected: []bool{!deleted, deleted, !deleted, deleted},
		},
		{
			id:       "60s interval, 30s age splay",
			interval: 60 * time.Second,
			objects:  []time.Duration{0 * time.Second, 30 * time.Second, 60 * time.Second, 90 * time.Second},
			expected: []bool{!deleted, deleted, !deleted, deleted},
		},
		{
			id:       "45s interval, 30s age splay",
			interval: 45 * time.Second,
			objects:  []time.Duration{0 * time.Second, 30 * time.Second, 60 * time.Second, 90 * time.Second},
			expected: []bool{!deleted, deleted, !deleted, deleted},
		},
		{
			id:       "90s interval, 30s age splay, shuffled",
			interval: 90 * time.Second,
			objects:  []time.Duration{30 * time.Second, 60 * time.Second, 0 * time.Second, 90 * time.Second},
			expected: []bool{deleted, deleted, !deleted, !deleted},
		},
	} {
		t.Logf("Testing case %q", tc.id)

		b := newBucket(agerotate.Range{irrelevantDuration, tc.interval})
		objs := make([]*testObject, len(tc.objects))
		for i := range tc.objects {
			objs[i] = &testObject{age: tc.objects[i]}
			b.Add(objs[i])
		}

		b.Cleanup()

		if len(tc.expected) != len(b.objects) {
			t.Fatalf("Expected %v results, got %v", len(tc.expected), len(b.objects))
		}
		for i := range tc.expected {
			if tc.expected[i] != objs[i].deleted {
				got := make([]bool, len(objs))
				for j := range objs {
					got[j] = objs[j].deleted
				}
				t.Fatalf("Expected deleted: %v, got %v", tc.expected, got)
			}
		}
	}
}
