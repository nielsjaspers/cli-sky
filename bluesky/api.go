package bluesky

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	// "github.com/bluesky-social/indigo/api/atproto"
	// "github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
)

func CreateSession() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	identifier := os.Getenv("BLUESKY_HANDLE")
	if identifier == "" {
		log.Fatalln("BLUESKY_HANDLE is not set in the environment.")
	}
	password := os.Getenv("BLUESKY_APP_PASSWORD")
	if password == "" {
		log.Fatalln("BLUESKY_APP_PASSWORD is not set in the environment.")
	}

	var postURL string = "https://bsky.social/xrpc/com.atproto.server.createSession"

    body := BlueskyAuth{
        Identifier: identifier,
        Password: password,
    }

    bodyBytes, err := json.Marshal(&body)
	if err != nil {
        log.Fatalf("Error parsing body to bytes: %v", err)
	}

    reader := bytes.NewReader(bodyBytes)

    resp, err := http.Post(postURL, "application/json", reader)
	if err != nil {
		log.Fatal(err)
	}

    defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

    // Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", resp.StatusCode)
	}

	log.Println("Response:", string(responseBody))
}

func Post(post *BlueskyPost) {
}
