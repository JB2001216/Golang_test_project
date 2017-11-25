/*
Copyright 2014 Tamás Gulácsi.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package punchhole

import (
	"errors"
	"os"
)

var errNoPunch = errors.New("punchHole not supported")

// PunchHole punches a hole in f from offset to offset+size, if non-nil.
var PunchHole func(file *os.File, offset, size int64) error
