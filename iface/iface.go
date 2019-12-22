//          FILE:  iface.go
//
//         USAGE:  iface.go
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
package iface

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

type IFace struct {
	hostname string
	uid      int
}

func NewIFace(hostname string) *IFace {
	retv := IFace{hostname: hostname, uid: os.Getuid()}
	return &retv
}

func times(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

// Header1 wraps `instr` as a H1 header
func (ifc IFace) Header1(instr string) string {
	return "\n" + instr + "\n" + strings.Repeat("=", len(instr))
}

// Header2 wraps `instr` as a H2 header
func (ifc IFace) Header2(instr string) string {
	return "\n" + instr + "\n" + strings.Repeat("-", len(instr))
}

func (ifc IFace) PrintHeader1(instr string) {
	fmt.Println(ifc.Header1(instr))
}

func (ifc IFace) PrintHeader2(instr string) {
	fmt.Println(ifc.Header2(instr))
}

func (ifc IFace) PrintParamInt(paramname string, paramvalue int) {
	yellow := color.New(color.FgYellow).SprintFunc()
	intparam := fmt.Sprintf("%-8d", paramvalue)
	fmt.Printf("%-70s %8s\n", ifc.PadRight(paramname, 70, "."), yellow(intparam))
}

func (ifc IFace) PrintParamStr(paramname string, paramvalue string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%-70s %8s\n", ifc.PadRight(paramname, 70, "."), yellow(paramvalue))
}

func (ifc IFace) PrintParamTest(paramname string, trueval string, falseval string, testval bool) {
	if testval {
		ifc.Success(paramname + " (" + trueval + ")")
	} else {
		ifc.Failure(paramname + " (" + falseval + ")")
	}
}

func (ifc IFace) PadRight(str string, length int, pad string) string {
	return str + " " + times(pad, length-len(str)-1)
}

func (ifc IFace) Success(titlestr string) {
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%-70s [ %-16s ]\n", ifc.PadRight(titlestr, 70, "."), green("SUCCESS"))
}

func (ifc IFace) Failure(titlestr string) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%-70s [ %-16s ]\n", ifc.PadRight(titlestr, 70, "."), red("FAILED"))
}

// vim: noexpandtab filetype=go
