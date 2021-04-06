package main

import (
	"github.com/spf13/cobra"

	"github.com/tinywell/baas/cmd/version"
)

var (
	rootCmd = &cobra.Command{
		Use:   "baas",
		Short: "baas is a manage plantform for hyperledger/fabric",
	}
)

func init() {
	rootCmd.AddCommand(version.VersionCMD)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
