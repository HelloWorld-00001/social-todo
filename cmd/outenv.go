package cmd

import (
	ssu "github.com/coderconquerer/social-todo/cmd/servicesetup"
	"github.com/spf13/cobra"
)

var outEnvCmd = &cobra.Command{
	Use:   "outenv",
	Short: "Output all environment variables to std",
	Run: func(cmd *cobra.Command, args []string) {
		ssu.NewServices().OutEnv()
	},
}
