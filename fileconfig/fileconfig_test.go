package fileconfig

import (
	"strings"
	"testing"
	"time"

	"agerotate"
)

const (
	fullInput = `
# Some standalone comment, followed by a blank line

pathGLOB:/path/to/whatever
range:3600:0        # Keep everything from the last hour
range:21600:120	    # Keep one every two hours younger than six hours
# range:21655:123  This one is ignored
RANGE:604800:43200  # Keep one every 12 hours for the last week
`
)

func TestPareLineErrors(t *testing.T) {
	for _, tc := range []struct {
		id          string
		line        string
		expectedErr string
	}{
		{
			id:          "empty line okay",
			line:        "    # gibberish range:   ",
			expectedErr: "",
		},
		{
			id:          "Not enough fields",
			line:        "blurg",
			expectedErr: "Line 0: Missing values",
		},
		{
			id:          "Invalid prefix",
			line:        "blurg:",
			expectedErr: "Line 0: Invalid prefix \"blurg\"",
		},
		{
			id:          "Missing path",
			line:        "PATHGLOB:",
			expectedErr: "Line 0: Must specify path",
		},
		{
			id:          "Multiple paths",
			line:        "PaTHGLOB:a:b",
			expectedErr: "Line 0: Using multiple path values is invalid",
		},
		{
			id:          "Incomplete range",
			line:        "rANgE:0",
			expectedErr: "Line 0: Range lines must have two values",
		},
		{
			id:          "Overspecified range",
			line:        "rANgE:0:b:z",
			expectedErr: "Line 0: Range lines must have two values",
		},
	} {
		t.Logf("Testing case %q", tc.id)
		p := parser{line: tc.line, fieldSep: ":"}
		err := p.parseLine()
		if err == nil {
			if tc.expectedErr != "" {
				t.Fatalf("Expected error %q, got nil", tc.expectedErr)
			}
		} else {
			if err.Error() != tc.expectedErr {
				t.Fatalf("Expected error %q, got %q", tc.expectedErr, err)
			}
		}
	}
}

func TestAddRangeErrors(t *testing.T) {
	for _, tc := range []struct {
		id            string
		line          string
		expectedErr   string
		expectedRange struct {
			age      int
			interval int
		}
	}{
		{
			id:          "Non-int age",
			line:        "RaNGe:1.4:15",
			expectedErr: "Line 0: Invalid age integer: strconv.ParseInt: parsing \"1.4\": invalid syntax",
		},
		{
			id:          "Non-int interval",
			line:        "RaNGe:15:0.6",
			expectedErr: "Line 0: Invalid interval integer: strconv.ParseInt: parsing \"0.6\": invalid syntax",
		},
		{
			id:          "Negative age",
			line:        "RaNGe:-5:0",
			expectedErr: "Line 0: Age values must be positive, got -5",
		},
		{
			id:          "Negative interval",
			line:        "RaNGe:5:-7",
			expectedErr: "Line 0: Interval values must be positive, got -7",
		},
	} {
		t.Logf("Testing case %q", tc.id)
		p := parser{line: tc.line, fieldSep: ":"}
		err := p.parseLine()
		if tc.expectedErr == "" {
			if err != nil {
				t.Fatalf("Expected no error, got %q", err)
			}
			if len(p.ranges) != 1 {
				t.Fatalf("Expected 1 parsed range, got %d", len(p.ranges))
			}
			expAge := time.Duration(tc.expectedRange.age)
			expInterval := time.Duration(tc.expectedRange.interval)
			if p.ranges[0].Age != expAge || p.ranges[0].Interval != expInterval {
				t.Fatalf("Expected age/interval %v, got range %v", tc.expectedRange, p.ranges[0])
			}
		} else {
			if err == nil {
				t.Fatalf("Expected err %q, got nil", tc.expectedErr)
			}
			if tc.expectedErr != err.Error() {
				t.Fatalf("Expected error %q, got %q", tc.expectedErr, err)
			}
		}
	}
}

func TestRangeLarger(t *testing.T) {
	line := "RAnGE:60:15   # Extra fluff"
	p := parser{
		line:     line,
		fieldSep: ":",
		ranges: []agerotate.Range{
			agerotate.Range{Age: 60 * time.Second},
		},
	}
	expected := "Line 0: Age value must be larger than previous age value"
	err := p.parseLine()
	if err == nil {
		t.Fatalf("Expected error %q, got nil", expected)
	}
	if expected != err.Error() {
		t.Fatalf("Expected error %q, got %q", expected, err)
	}
}

func TestFull(t *testing.T) {
	expectedPath := "/path/to/whatever"
	expectedRanges := []agerotate.Range{
		agerotate.Range{
			Age:      3600 * time.Second,
			Interval: 0 * time.Second,
		},
		agerotate.Range{
			Age:      21600 * time.Second,
			Interval: 120 * time.Second,
		},
		agerotate.Range{
			Age:      604800 * time.Second,
			Interval: 43200 * time.Second,
		},
	}

	in := strings.NewReader(fullInput)
	files, ranges, err := Parse(in, ":")
	if err != nil {
		t.Fatalf("Got unexpected error %q", err)
	}
	if string(files) != expectedPath {
		t.Fatalf("Expected files path %q, got %q", expectedPath, files)
	}
	if len(ranges) != len(expectedRanges) {
		t.Fatalf("Expected %d ranges, got %d", len(expectedRanges), len(ranges))
	}
	for i := range expectedRanges {
		if expectedRanges[i] != ranges[i] {
			t.Fatalf("Expected range %d to be %q, got %q", i, expectedRanges[i], ranges[i])
		}
	}
}
