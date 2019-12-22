//          FILE:  input.go
// 
//         USAGE:  input.go
// 
//   DESCRIPTION:  $description
// 
//       OPTIONS:  ---
//  REQUIREMENTS:  ---
//          BUGS:  ---
//         NOTES:  ---
//        AUTHOR:  John van Zantvoort (jvzantvoort), john@vanzantvoort.org
//       COMPANY:  JDC
//       CREATED:  21-Dec-2019
//
// Copyright (C) 2019 JDC
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
//
package input

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	MEMFILE string = "/proc/meminfo"
	FSFILE string = "/proc/filesystems"
	MOUNTSFILE string = "/proc/mounts"
)

// Return the short hostname
func ShortHostname() string {
	name, err := os.Hostname()
	cols := strings.Split(name, ".")
	if err != nil {
		panic(err)
	}
	return cols[0]
}

// GetLines gets all the lines of a file and returns them as a slice of
// strings.
func GetLines(infile string)(retv []string) {
	file, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		retv = append(retv, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return

}

// MemTotalKb returns the total amount of memory found in `/proc/meminfo` in
// kilobytes.
func MemTotalKb()(retint float64) {
	re := regexp.MustCompile(`(MemTotal)\:[ \t]*([0-9]*)[ \t]*kB`)
	for _, line := range GetLines(MEMFILE) {
		if re.MatchString(line) {
			match := re.FindStringSubmatch(line)
			retstr, _ := strconv.Atoi(match[2])
			retint = float64(retstr)
		}
	}
	return
}

// Return a list of used filesystems
func Filesystems()(retv []string) {
	re := regexp.MustCompile("nodev")
	for _, line := range GetLines(FSFILE) {
		if re.MatchString(line) {
			continue
		}
		line = strings.Replace(line, "\t", "", 10)
		retv = append(retv, line)
	}
	return
}

// Lookup and return mountpoints for a provided filesystem
func ProcMounts(fs string)(retv []string) {
	for _, line := range GetLines(MOUNTSFILE) {
		cols := strings.Split(line, " ")
		fstype := cols[2]
		if fs == fstype {
			retv = append(retv, cols[1])
		}
	}
	return
}


// vim: noexpandtab filetype=go
