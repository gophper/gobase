// Copyright  2015  gophper.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.
package gobase

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"reflect"
	"runtime"
	"syscall"
)

var (
	// 日志
	Log *BaseLog
	//配置文件
	Config     *GoConfig
	SigHandler = make(map[string]interface{})

	// 定义命令行参数
	verbose = flag.Bool("v", false, "Verbose output")
	help    = flag.Bool("h", false, "Show this help")
	chroot  = flag.Bool("w", false, "Setup chroot")
	cfgfile = flag.String("c", "", "Config file")
	workdir = flag.String("d", "", "Setup work dir")
	pidfile = flag.String("p", "", "Pid file")
)

func init() {
	flag.Parse()
	if *help {
		Help()
		return
	}

	if *workdir != "" {
		fmt.Println("workdir: ", *workdir, os.Args)
		if err := syscall.Chdir(*workdir); err != nil {
			fmt.Printf("Can't change to work dir [%s]: %s\n", *workdir, err)
			os.Exit(1)
		}

		if *chroot {
			pwd, _ := os.Getwd()
			if err := syscall.Chroot(pwd); err != nil {
				fmt.Printf("Can't Chroot to %s: %s\n", *workdir, err)
				os.Exit(1)
			}
			fmt.Printf("I'll Chroot to %s !\n", pwd)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	Config = LoadConfig("")
	CreatePid()
	Log = defaultLog()

	SigHandler["sighup"] = func() {
		Log.Debug("reload config")
		Config = LoadConfig("")
	}

	if ok, _ := Config.Bool("sys.signal", false); ok {
		go SignalHandle(SigHandler)
	}
}

func LoadConfig(configFile string) (cfg *GoConfig) {
	var err error
	if configFile == "" {
		if *cfgfile != "" {
			configFile = *cfgfile
		} else {
			configFile = "etc/" + path.Base(os.Args[0]) + ".conf"
		}
	}
	if !IsExist(configFile) {
		fmt.Println("config is not set,use the default configuration")
		configFile = "gobase.conf"
		FilePutContent(configFile, "")
	}
	cfg, err = NewConfig(configFile, 5)
	if err != nil {
		fmt.Println("read config file error: ", err)
		os.Exit(1)
	}

	return cfg
}

func defaultLog() *BaseLog {
	logType, _ := Config.String("log.type", "console")
	logFile, _ := Config.String("log.file", "")
	logLevel, _ := Config.Int("log.level", 5)
	logFlag, _ := Config.Int("log.flag", 9)
	opt := &LogOptions{Type: logType, File: logFile, Level: logLevel, Flag: logFlag}
	return NewLog(opt)
}

func SignalHandle(funcs map[string]interface{}) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGHUP)

	for {
		select {
		case s := <-ch:
			switch s {
			default:
			case syscall.SIGHUP:
				if f, ok := funcs["sighup"]; ok {
					if ff := reflect.ValueOf(f); ff.Kind() == reflect.Func {
						ff.Call(nil)
					}
				}
				break
			case syscall.SIGINT:
				if f, ok := funcs["sigint"]; ok {
					if ff := reflect.ValueOf(f); ff.Kind() == reflect.Func {
						ff.Call(nil)
					}
				}
				os.Exit(1)
			case syscall.SIGUSR1:
				if f, ok := funcs["sigusr1"]; ok {
					if ff := reflect.ValueOf(f); ff.Kind() == reflect.Func {
						ff.Call(nil)
					}
				}
			case syscall.SIGUSR2:
				if f, ok := funcs["sigusr2"]; ok {
					if ff := reflect.ValueOf(f); ff.Kind() == reflect.Func {
						ff.Call(nil)
					}
				}
			}
		}
	}
}

// create pid file
func CreatePid() {
	pid := os.Getpid()

	if pid < 1 {
		fmt.Println("Get pid err")
		os.Exit(1)
	}

	var pidf string
	if *pidfile != "" {
		pidf = *pidfile
	} else {
		pidf, _ = Config.String("sys.pid", "")
		if pidf == "" {
			pidf = "/var/run/" + path.Base(os.Args[0]) + ".pid"
		}
	}

	if pidf == "" {
		fmt.Println("pid file not setup")
		os.Exit(1)
	}

	f, err := os.OpenFile(pidf, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open pid file err ", err)
		os.Exit(1)
	}

	f.WriteString(GetIntStr(pid))
	f.Close()
}

func Help() {
	fmt.Printf(
		"\nUseage: %s [ Options ]\n\n"+
			"Options:\n"+
			"  -c Server config file [Default: etc/serverd.conf]\n"+
			"  -d Work dir [Default: publish]\n"+
			"  -h Display this mssage\n"+
			"  -p Pid file [Default: /var/run/serverd.pid]\n"+
			"  -w Enable chroot to work dir [Required: -d ]\n\n"+
			"------------------------------------------------------\n\n",
		os.Args[0])

	os.Exit(0)
}
