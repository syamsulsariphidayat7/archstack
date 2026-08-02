package cli

import (
	"fmt"
)

var Version = "dev"

func versionCmd() error {
	fmt.Printf("archstack %s\n", Version)
	return nil
}
