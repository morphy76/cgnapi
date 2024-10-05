package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/morphy76/cgnapi/internal/command"
	"github.com/morphy76/cgnapi/internal/configuration"
)

var profile string
var config bool
var authServer string
var realm string
var clientID string
var refreshToken string

var add bool
var remove bool
var list bool
var initToken bool

var renew bool

var help bool

func init() {
	flag.StringVar(&profile, "p", "", "profile name")
	flag.BoolVar(&config, "config", false, "configure profiles")
	flag.StringVar(&authServer, "auth-url", "", "auth server url, to be used with -config in conjunction with -add")
	flag.StringVar(&realm, "realm", "", "realm, to be used with -config in conjunction with -add")
	flag.StringVar(&clientID, "client-id", "", "client id, to be used with -config in conjunction with -add")
	flag.StringVar(&refreshToken, "refresh-token", "", "refresh token, to be used with -config in conjunction with -add (opt) and -init")
	flag.BoolVar(&help, "help", false, "Show help")

	flag.BoolVar(&add, "add", false, "config action, add a new profile, to be used with -config")
	flag.BoolVar(&remove, "remove", false, "config action, remove an existing profile, to be used with -config")
	flag.BoolVar(&list, "list", false, "config action, list profiles, no need to specify -config")
	flag.BoolVar(&initToken, "init", false, "config action, init the refresh token for an existing profile, to be used with -config")

	flag.BoolVar(&renew, "renew", false, "command action, renew the access token for the given profile")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "\033[1;34mUsage:\033[0m\n")
		flag.PrintDefaults()

		fmt.Fprint(flag.CommandLine.Output(), "\n\033[1;34mExamples:\033[0m\n")
		fmt.Fprint(flag.CommandLine.Output(), "  \033[1;32mConfigure a new profile:\033[0m\n")
		fmt.Fprint(flag.CommandLine.Output(), "    cgnapi -config -add -p <profile_name> -client-id <clientID> -realm <realm> -auth-url <auth_server_url> [-refresh-token <refresh_token>]\n")
		fmt.Fprint(flag.CommandLine.Output(), "  \033[1;32mList profiles:\033[0m\n")
		fmt.Fprint(flag.CommandLine.Output(), "    cgnapi -list\n")
		fmt.Fprint(flag.CommandLine.Output(), "  \033[1;32mInit refresh token for an existing profile:\033[0m\n")
		fmt.Fprint(flag.CommandLine.Output(), "    run -config -init -p <profile_name> -refresh-token <refresh_token>\n")
		fmt.Fprint(flag.CommandLine.Output(), "  \033[1;32mRenew the access token for an existing profile:\033[0m\n")
		fmt.Fprint(flag.CommandLine.Output(), "    run -renew -p <profile_name>\n")
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
		err := configuration.Main(profile, add, remove, list, initToken, realm, clientID, authServer, refreshToken)
		if err != nil {
			fmt.Println("\033[1;31m" + err.Error() + "\033[0m")
			os.Exit(1)
		}
	} else {
		err := command.Main(profile, renew)
		if err != nil {
			fmt.Println("\033[1;31m" + err.Error() + "\033[0m")
			os.Exit(1)
		}
	}
}
