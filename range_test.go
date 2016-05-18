/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package agerotate

import (
	"sort"
	"testing"
	"time"
)

func TestRangeString(t *testing.T) {
	for _, tc := range []struct {
		id            string
		age, interval time.Duration
		expected      string
	}{
		{
			id:       "3 days @ 12 hours",
			age:      3 * 24 * time.Hour,
			interval: 12 * time.Hour,
			expected: "For files younger than 72h0m0s, keep one every 12h0m0s",
		},
		{
			id:       "12 hours @ 3 hours",
			age:      12 * time.Hour,
			interval: 3 * time.Hour,
			expected: "For files younger than 12h0m0s, keep one every 3h0m0s",
		},
	} {
		t.Logf("Testing case %q", tc.id)
		r := Range{tc.age, tc.interval}
		if r.String() != tc.expected {
			t.Fatalf("Got %q, expected %q", r, tc.expected)
		}
	}
}

func TestByAgeSort(t *testing.T) {
	for _, tc := range []struct {
		id       string
		ages     []time.Duration
		expected []time.Duration
	}{
		{
			id:       "Already sorted",
			ages:     []time.Duration{3 * time.Hour, 5 * time.Hour, 10 * time.Hour},
			expected: []time.Duration{3 * time.Hour, 5 * time.Hour, 10 * time.Hour},
		},
		{
			id:       "Reversed",
			ages:     []time.Duration{10 * time.Hour, 5 * time.Hour, 3 * time.Hour},
			expected: []time.Duration{3 * time.Hour, 5 * time.Hour, 10 * time.Hour},
		},
		{
			id:       "Shuffled",
			ages:     []time.Duration{5 * time.Hour, 3 * time.Hour, 10 * time.Hour},
			expected: []time.Duration{3 * time.Hour, 5 * time.Hour, 10 * time.Hour},
		},
	} {
		t.Logf("Testing case %q", tc.id)
		ranges := make([]Range, len(tc.ages))
		for i, age := range tc.ages {
			ranges[i] = Range{Age: age}
		}
		sort.Sort(ByAge(ranges))
		got := make([]time.Duration, len(tc.ages))
		for i, r := range ranges {
			got[i] = r.Age
		}
		for i, gotAge := range got {
			if gotAge != tc.expected[i] {
				t.Fatalf("In case %q, got %v, expected %v", tc.id, got, tc.expected)
			}
		}
	}
}
