/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AgentZombie/agerotate/bucket"
	"github.com/AgentZombie/agerotate/fileobject/config"
)

var (
	ConfigPath = flag.String("config", "", "Path to file rotation config.")
	FieldSep   = flag.String("fieldsep", ":", "Field separator for range lines.")
	ShowFormat = flag.Bool("showfmt", false, "Take no action, just print the config format.")
)

func errorExit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(-1)
}

func showFormat() {
	fmt.Printf(`
Configuration is done with a simple text file having one configuration 
directive per line. For each line, leading and trailing whitespace is ignored. 
Any text followed by # is ignored, including the #. Blank lines and lines 
composed of only whitespace and/or characters prefixed by # are ignored.

There are two configuration directives: PATHGLOB and RANGE. PATHGLOB and RANGE
need not be capitalized. All configuration directives are followed by %s, and
then one or more values separated by %s. The separator can be changed with the
fieldsep command line flag.

PATHGLOB specifies a filesystem glob to select files for rotation. The PATHGLOB
line is required, can appear anywhere in the file, and must appear only once.

RANGE identifies a set of files for rotation by their age. Each RANGE line has
exactly two values: Age and Interval. Files with mtimes younger than Age but
greater than or equal to the Age on the previous RANGE line match this RANGE.
Interval specifies a minimum length of time between files to keep within this
range.

The youngest file in a RANGE is always retained. Beyond that a file is only
retained if it's age is at least Interval greater than the last file retained.

The config file may have as many RANGE lines as are desired. The Age of each
RANGE line must be unique and must be greater than the Age of the previous
RANGE lines, except for the first RANGE line.

An Interval of 0 is valid and causes all files in the range to be retained.
It's often a good idea for the first RANGE to have an Interval of 0 so all
of the most recent files are kept.

Times for Age and Interval are specified using the syntax specified in
https://golang.org/pkg/time/#ParseDuration. Units larger than "h" are not
available because calendar math is frought with peril.

Sample Config:
  # RANGE:Age:Interval
  pathglob:/path/to/files/*.gz
  range:6h:0      # For the first six hours, keep all.
  range:72h:4h	  # For files under 72 hours, keep one per 4 hours.
  range:720h:24h  # For files under 30 days, keep one per day.
  range:4320h:72h # For files under six months, keep one every 3 days.
  # Beyond six months, files are deleted.
`, *FieldSep, *FieldSep)
}

func main() {
	flag.Parse()

	if *ShowFormat {
		showFormat()
		return
	}

	cfg, err := os.Open(*ConfigPath)
	if err != nil {
		errorExit("Error opening config %q: %v\n", *ConfigPath, err)
	}

	files, ranges, err := config.Parse(cfg, *FieldSep)
	if err != nil {
		errorExit("Error parsing config %q: %v\n", *ConfigPath, err)
	}

	err = bucket.Cleanup(ranges, files)
	if err != nil {
		errorExit("Error doing cleanup: %v\n", err)
	}
}
