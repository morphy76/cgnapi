package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/morphy76/cgnapi/internal/configuration"
)

var profile string
var config bool
var authServer string
var refreshToken string

var add bool
var remove bool
var list bool
var initToken bool

var help bool

func init() {
	flag.StringVar(&profile, "p", "", "profile name")
	flag.BoolVar(&config, "config", false, "config file")
	flag.StringVar(&authServer, "auth-url", "", "auth server url")
	flag.StringVar(&refreshToken, "refresh-token", "", "refresh token")
	flag.BoolVar(&help, "help", false, "Show help")

	flag.BoolVar(&add, "add", false, "add profile")
	flag.BoolVar(&remove, "remove", false, "remove profile")
	flag.BoolVar(&list, "list", false, "list profiles")
	flag.BoolVar(&initToken, "init", false, "init the refresh token for a profile")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "\033[1;34mUsage:\033[0m\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if profile == "" && !list {
		fmt.Println("\033[1;31mprofile is required\033[0m")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	if config || list {
		err := configuration.Main(profile, add, remove, list, initToken, authServer, refreshToken)
		if err != nil {
			fmt.Println("\033[1;31m" + err.Error() + "\033[0m")
			os.Exit(1)
		}
	} else {
		runCmd()
	}
}

func runCmd() {
	panic("unimplemented2")
}
