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
	"log/syslog"
	"math"
	"os"
	"os/user"
	"path"
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

	defer func() {
		if panicname := recover(); panicname != nil {
			log.Println("Found error", panicname)
			os.Exit(1)
		}
	}()

	var result bool
	// Setup log
	log.SetFlags(0)
	logwriter, e := syslog.New(syslog.LOG_INFO, "system_info")
	if e == nil {
		log.SetOutput(logwriter)
	}

	// Get user info
	usrobj, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homedir := usrobj.HomeDir
	username := usrobj.Username
	uid := os.Getuid()
	gid := os.Getgid()

	configuration := getconfiguration(homedir)
	hostname := input.ShortHostname()
	mem := input.MemTotalKb()

	ifc := iface.NewIFace(hostname)

	// Hostname header
	ifc.SetTitle(hostname)
	ifc.PrintHeader1()

	// NUMA
	ifc.SetTitle("NUMA")

	ifc.IncIndent()
	ifc.PrintParamInt("Processors", runtime.NumCPU())
	ifc.PrintParamInt("Memory (GiB)", int(math.Round(mem/1048576.0)))
	ifc.DecIndent()

	// DC
	ifc.SetTitle("Datacenter")
	ifc.PrintHeader2()

	ifc.IncIndent()
	ifc.PrintParamStr("Name", configuration.Datacenter.Name)
	ifc.PrintParamStr("Client", configuration.Datacenter.Client)
	ifc.PrintParamStr("Location", configuration.Datacenter.Location)
	ifc.DecIndent()

	// CU
	ifc.SetTitle("Current User")
	ifc.PrintHeader2()

	ifc.IncIndent()
	ifc.PrintParamStr("username", username)
	ifc.PrintParamStr("homedir", homedir)
	ifc.PrintParamInt("uid", uid)
	ifc.PrintParamInt("gid", gid)
	ifc.DecIndent()

	// Filesystem
	ifc.SetTitle("Filesystems")
	ifc.PrintHeader2()

	ifc.IncIndent()
	for _, fs := range input.Filesystems() {
		for _, mountPoing := range input.ProcMounts(fs) {
			ifc.PrintParamStr(mountPoing, fs)
		}
	}
	ifc.DecIndent()

	// Config
	ifc.SetTitle("Config switches")
	ifc.PrintHeader2()

	ifc.IncIndent()
	ifc.PrintParamTest("Security", "enabled", "disabled", configuration.Server.Secure)
	ifc.DecIndent()

	// Tests
	ifc.SetTitle("Tests")
	ifc.PrintHeader2()

	ifc.IncIndent()
	bashrc := path.Join(homedir, ".bashrc")
	sshdir := path.Join(homedir, ".ssh")
	sshauth := path.Join(sshdir, "authorized_keys")

	result = apptests.TargetParameters(bashrc, 0644, uid, gid)
	ifc.PrintParamTest("bashrc has correct permissions", "yes", "no", result)

	result = apptests.TargetParameters(sshdir, 0700, uid, gid)
	ifc.PrintParamTest("sshdir has correct permissions", "yes", "no", result)

	result = apptests.TargetParameters(sshauth, 0600, uid, gid)
	ifc.PrintParamTest("authorized_keys has correct permissions", "yes", "no", result)

	ifc.DecIndent()

	// ifc.PrintParamTest("  SELinux", "enabled", "disabled", selinux.GetEnabled())

	fmt.Printf("\n\n")
}

// vim: noexpandtab filetype=go
