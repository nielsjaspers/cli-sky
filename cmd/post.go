package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nielsjaspers/cli-sky/bluesky"
	"github.com/nielsjaspers/cli-sky/internal/datahandler"
	"github.com/spf13/cobra"
)

var handlePost string

var postCmd = &cobra.Command{
	Use:   "post \"message\"",
	Short: "Post a message to Bluesky",
	Long:  "Allows you to post a message to Bluesky with an optional handle.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]

		// Unescape the string to interpret \n as newlines
		message, err := strconv.Unquote(`"` + message + `"`)
		if err != nil {
			log.Fatalf("Error unescaping message: %v", err)
			return
		}

		// Clean the handle (remove "@" if present)
		cleanedHandle := strings.TrimPrefix(handlePost, "@")

		authData, err := datahandler.ReadAuthResponseFromFile(cleanedHandle)
		if err != nil {
			log.Fatalf("Error reading Auth data: %v", err)
		}

		fmt.Printf("Posting message: %s\n", message)
		fmt.Printf("Using handle: %s\n", cleanedHandle)
		bluesky.Post(message, authData)
	},
}

func init() {
	rootCmd.AddCommand(postCmd)
	postCmd.Flags().StringVarP(&handlePost, "handle", "u", "", "Bluesky handle (optional, @handle.bsky.social)")
}

