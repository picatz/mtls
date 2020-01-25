package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "mtlssh"}
	rootCmd.AddCommand(serverCommand)
	rootCmd.AddCommand(clientCommand)
	rootCmd.AddCommand(certCommand)
	rootCmd.AddCommand(proxyCommand)
	rootCmd.Execute()
}
