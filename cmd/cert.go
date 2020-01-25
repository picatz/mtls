package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var certCommand = &cobra.Command{
	Use:   "cert",
	Short: "mTLS SSH cert commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not yet implemented!")
	},
}
