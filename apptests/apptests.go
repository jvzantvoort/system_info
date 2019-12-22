//          FILE:  tests.go
// 
//         USAGE:  tests.go
// 
//   DESCRIPTION:  $description
// 
//       OPTIONS:  ---
//  REQUIREMENTS:  ---
//          BUGS:  ---
//         NOTES:  ---
//        AUTHOR:  John van Zantvoort (jvzantvoort), john@vanzantvoort.org
//       COMPANY:  JDC
//       CREATED:  22-Dec-2019
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
package apptests

import (
    "fmt"
    "os"
    "log"
    "syscall"
    "strconv"
)

func TargetParameters(targetpath string, mode int, uid int, gid int) bool {
	retv := true
	if ! TargetHasPermission(targetpath, mode) {
		retv = false
	}
	if ! TargetHasUID(targetpath, uid) {
		retv = false
	}
	if ! TargetHasGID(targetpath, gid) {
		retv = false
	}
	return retv
}

func TargetHasPermission(targetpath string, mode int) bool {
	targetStat, err := os.Stat(targetpath)
	if err != nil {
		log.Fatal(err)
	}
	target_sys := targetStat.Sys()
	num := int(target_sys.(*syscall.Stat_t).Mode & 0777)
	if num == mode {
		return true
	} else {
	return false
}
}

func TargetHasUID(targetpath string, uid int) bool {
	targetStat, err := os.Stat(targetpath)
	if err != nil {
		log.Fatal(err)
	}
	target_sys := targetStat.Sys()
	target_uid, _ := strconv.Atoi(fmt.Sprint(target_sys.(*syscall.Stat_t).Uid))
	if uid == target_uid {
		return true
	} else {
		return false
	}
}

func TargetHasGID(targetpath string, gid int) bool {
	targetStat, err := os.Stat(targetpath)
	if err != nil {
		log.Fatal(err)
	}
	target_sys := targetStat.Sys()
	target_gid, _ := strconv.Atoi(fmt.Sprint(target_sys.(*syscall.Stat_t).Gid))
	if gid == target_gid {
		return true
	} else {
		return false
	}
}


// vim: noexpandtab filetype=go
