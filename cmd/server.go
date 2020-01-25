package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "mTLS SSH server commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!")
	},
}
