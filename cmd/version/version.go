package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// version 信息
var (
	version   string
	commitID  string
	buildTime string

	// VersionCMD 打印程序版本信息
	VersionCMD = &cobra.Command{
		Use:   "version",
		Short: "baas 版本信息 ",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s:\n Version: %s\n Commit SHA: %s\n Build Time: %s\n Go version: %s\n OS/Arch: %s\n",
				"baas", version, commitID, buildTime, runtime.Version(), fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
		},
	}
)
