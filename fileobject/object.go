/* Copyright (c) 2016 Jason Mansfield


Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

// fileobject implements rotation for filesystem objects using os.Stat and os.Remove.
package fileobject

import (
	"os"
	"path/filepath"
	"time"

	"agerotate"
)

// File captures a file path and it's mtime, providing methods for the Object interface. The mtime is cached to avoid hammering the filesystem during sorting.
type File struct {
	path string
	age  time.Duration
}

func newFile(path string) (File, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}
	return File{
		path: path,
		age:  time.Now().Sub(fi.ModTime()),
	}, nil
}

// ID returns the path for the file object.
func (f File) ID() string {
	return f.path
}

// Age returns the age of the object as a time.Duration.
func (f File) Age() time.Duration {
	return f.age
}

// Delete attempts to remove the file object. No error is returned if it already doesn't exist.
func (f File) Delete() error {
	err := os.Remove(f.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Files wraps a string (assumed to be a path glob) to provide Objects operations on it.
type Files string

// ID returns the path glob for the object.
func (f Files) ID() string {
	return string(f)
}

// List returns the File items matching the glob.
func (f Files) List() ([]agerotate.Object, error) {
	paths, err := filepath.Glob(string(f))
	if err != nil {
		return nil, err
	}
	fObjs := []agerotate.Object{}
	for _, path := range paths {
		nf, err := newFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}
		fObjs = append(fObjs, nf)
	}
	return fObjs, nil
}
