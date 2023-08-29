/*
Copyright © 2023 Tommy Breslein <github.com/tbreslein>
*/
package main

import (
	"fmt"
	"os"

	"github.com/tbreslein/frankenrepo/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		println("everything worked out")
	}
}
