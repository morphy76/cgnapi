package command

import (
	"fmt"
)

func Main(
	profile string,
	renew, get, decoded, exp bool,
) error {
	if renew {
		err := RenewToken(profile)
		if err != nil {
			return err
		}
		fmt.Println("\033[1;32mToken renewed successfully!\033[0m")
	} else if get {
		err := GetToken(profile, decoded)
		if err != nil {
			return err
		}
	} else if exp {
		err := GetTokenExp(profile)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no command action specified")
	}
	return nil
}
