package command

import (
	"fmt"
)

func Main(
	profile string,
	renew bool,
) error {
	if renew {
		err := RenewToken(profile)
		if err != nil {
			return err
		}
		fmt.Println("\033[1;32mToken renewed successfully!\033[0m")
	} else {
		return fmt.Errorf("no command action specified")
	}
	return nil
}
