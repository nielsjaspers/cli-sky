package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/nielsjaspers/cli-sky/bluesky"
	"github.com/nielsjaspers/cli-sky/internal/datahandler"
	"github.com/spf13/cobra"
)

var refreshCmd = &cobra.Command{
	Use:   "refresh \"@handle.bsky.com\"",
	Short: "Refresh a session",
	Long:  "Allows you to refresh a session.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handle := args[0]

		// Clean the handle (remove "@" if present)
		cleanedHandle := strings.TrimPrefix(handle, "@")

        authData, err := datahandler.ReadAuthResponseFromFile(cleanedHandle)
        if err != nil {
            log.Fatalf("Error reading Auth data: %v", err)
        }

		fmt.Printf("Using handle: %s\n", cleanedHandle)
        bluesky.RefreshSession(authData.RefreshJwt)
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}

