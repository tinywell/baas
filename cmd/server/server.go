package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/tinywell/baas/internal/api"
)

// version 信息
var (
	version   string
	commitID  string
	buildTime string

	// VersionCMD 打印程序版本信息
	ServerCMD = &cobra.Command{
		Use:   "server",
		Short: "baas 服务 ",
		Run: func(cmd *cobra.Command, args []string) {
			server()
		},
	}
)

func server() {
	r := gin.Default()
	api.AddRouter(r)
	if err := r.Run("0.0.0.0:8080"); err != nil {
		panic(err)
	}
}
