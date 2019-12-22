//          FILE:  main.go
//
//         USAGE:  main.go
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
package main

import (
	"fmt"
	"github.com/jvzantvoort/system_info/apptests"
	"github.com/jvzantvoort/system_info/config"
	"github.com/jvzantvoort/system_info/iface"
	"github.com/jvzantvoort/system_info/input"
	////  "github.com/opencontainers/selinux/go-selinux"
	"github.com/spf13/viper"
	"log"
	"math"
	"path"
	"os"
	"os/user"
	"runtime"
)

func getconfiguration(homedir string) config.Configuration {
	var configuration config.Configuration
	viper.SetConfigName("sysinfo")  // name of config file (without extension)
	viper.AddConfigPath(homedir)    // optionally look for config in the working directory
	viper.AddConfigPath("/etc/jdc") // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	return configuration
}

func main() {

	usrobj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homedir := usrobj.HomeDir

	configuration := getconfiguration(homedir)
	hostname := input.ShortHostname()
	uid := os.Getuid()
	gid := os.Getgid()
	mem := input.MemTotalKb()

	defer func() {
		if panicname := recover(); panicname != nil {
			fmt.Println("Found error", panicname)
			os.Exit(1)
		}
	}()

	ifc := iface.NewIFace(hostname)

	// Hostname header
	ifc.PrintHeader1(hostname)
	ifc.PrintHeader2("NUMA")
	ifc.PrintParamInt("  Processors", runtime.NumCPU())
	ifc.PrintParamInt("  Memory (GiB)", int(math.Round(mem/1048576.0)))

	ifc.PrintHeader2("Datacenter")
	ifc.PrintParamStr("  Name", configuration.Datacenter.Name)
	ifc.PrintParamStr("  Client", configuration.Datacenter.Client)
	ifc.PrintParamStr("  Location", configuration.Datacenter.Location)

	ifc.PrintHeader2("Current User")
	ifc.PrintParamStr("  homedir", homedir)
	ifc.PrintParamInt("  uid", uid)
	ifc.PrintParamInt("  gid", gid)

	ifc.PrintHeader2("Filesystems")
	for _, fs := range input.Filesystems() {
		for _, mountPoing := range input.ProcMounts(fs) {
			ifc.PrintParamStr("  " + mountPoing, fs)
		}
	}
	ifc.PrintHeader2("Config switches")
	ifc.PrintParamTest("  Security", "enabled", "disabled", configuration.Server.Secure)

	ifc.PrintHeader2("Tests")

	var result bool

	bashrc := path.Join(homedir, ".bashrc")
	sshdir := path.Join(homedir, ".ssh")
	sshauth := path.Join(sshdir, "authorized_keys")

	result = apptests.TargetParameters(bashrc, 0644, uid, gid)
	ifc.PrintParamTest("  bashrc has correct permissions", "yes", "no", result)

	result = apptests.TargetParameters(sshdir, 0700, uid, gid)
	ifc.PrintParamTest("  sshdir has correct permissions", "yes", "no", result)

	result = apptests.TargetParameters(sshauth, 0600, uid, gid)
	ifc.PrintParamTest("  authorized_keys has correct permissions", "yes", "no", result)

	// ifc.PrintParamTest("  SELinux", "enabled", "disabled", selinux.GetEnabled())

	fmt.Printf("\n\n")
}

// vim: noexpandtab filetype=go
