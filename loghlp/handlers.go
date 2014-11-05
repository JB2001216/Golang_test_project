/*
Copyright 2014 Tamás Gulácsi

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

// Package loghlp collects some small log15 handlers
package loghlp

import (
	"flag"
	"os"

	"github.com/tgulacsi/go/loghlp/gloghlp"
	"github.com/tgulacsi/go/loghlp/tsthlp"
)

var (
	fsOrig = flag.CommandLine

	TestHandler = tsthlp.TestHandler
	GLogHandler = gloghlp.GLogHandler
)

func MaskFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func UnmaskFlags() {
	flag.CommandLine = fsOrig
}
