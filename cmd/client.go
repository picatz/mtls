package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var clientCommand = &cobra.Command{
	Use:   "client",
	Short: "mTLS SSH client commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!")
	},
}
