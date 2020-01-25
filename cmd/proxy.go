package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var proxyCommand = &cobra.Command{
	Use:   "proxy",
	Short: "setup mTLS proxy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!")
	},
}
