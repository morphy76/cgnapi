package configuration

import (
	"fmt"
	"sort"
	"strings"
)

func Main(
	profile string,
	add, remove, list, initToken bool,
	authServer, refreshToken string,
) error {
	if list {
		profiles, err := ListProfiles()
		if err != nil {
			return err
		}

		var longestName, longestURL int = len("Profile"), len("Auth Server")
		for name, cfg := range profiles {
			longestName = max(len(name), longestName)
			longestURL = max(len(cfg.AuthServer), longestURL)
		}

		rawPattern := fmt.Sprintf("\033[1;32m%%-%ds\033[0m | %%-%ds\n", longestName, longestURL)
		fmt.Printf("\033[1;34m"+rawPattern+"\033[0m", "Profile", "Auth Server")
		fmt.Println(strings.Repeat("\033[1;34m-\033[0m", longestName+longestURL+3))

		sortedKeys := make([]string, 0, len(profiles))
		for name := range profiles {
			sortedKeys = append(sortedKeys, name)
		}
		sort.Strings(sortedKeys)

		for _, name := range sortedKeys {
			cfg := profiles[name]
			fmt.Printf(rawPattern, name, cfg.AuthServer)
		}
	} else if add {
		err := AddProfile(profile, authServer, refreshToken)
		if err != nil {
			return err
		}
		fmt.Println("\033[1;32mProfile added successfully!\033[0m")
	} else if remove {
		err := RemoveProfile(profile)
		if err != nil {
			return err
		}
		fmt.Println("\033[1;32mProfile removed successfully!\033[0m")
	} else if initToken {
		err := InitToken(profile, refreshToken)
		if err != nil {
			return err
		}
		fmt.Println("\033[1;32mToken initialized successfully!\033[0m")
	} else {
		return fmt.Errorf("no action specified")
	}
	return nil
}
