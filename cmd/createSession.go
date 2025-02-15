package cmd

import (
	"fmt"
	"strings"

	"github.com/nielsjaspers/cli-sky/bluesky"
	"github.com/nielsjaspers/cli-sky/internal/datahandler"
	"github.com/spf13/cobra"
)

var createSessionCmd = &cobra.Command{
	Use:   "create \"@handle.bsky.com\"",
	Short: "Create a new session",
	Long:  "Allows you to create a new session.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handle := args[0]

		// Clean the handle (remove "@" if present)
		_ = strings.TrimPrefix(handle, "@")

		_, authResponse, responseBody := bluesky.CreateSession("")
		fmt.Println(responseBody)

		err := datahandler.WriteAuthResponseToFile(authResponse)
		if err != nil {
		    panic(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(createSessionCmd)
}
