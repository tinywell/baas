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

// @title baas 平台后端 API
// @version 1.0
// @description fabric 区块链管控台 - baas 后端 API

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
