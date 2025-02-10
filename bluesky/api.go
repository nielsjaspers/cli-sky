package bluesky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// "github.com/bluesky-social/indigo/api/atproto"
	// "github.com/bluesky-social/indigo/xrpc"
	"github.com/joho/godotenv"
)

func CreateSession(dotenvPath string) (*http.Response, *BlueskyAuthResponse, string) {
	if dotenvPath == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
			return nil, nil, ""
		}
	} else {
		err := godotenv.Load(dotenvPath)
		if err != nil {
			log.Fatalf("Error loading .env file")
			return nil, nil, ""
		}
	}

	identifier := os.Getenv("BLUESKY_HANDLE")
	if identifier == "" {
		log.Fatalln("BLUESKY_HANDLE is not set in the environment.")
		return nil, nil, ""
	}
	password := os.Getenv("BLUESKY_APP_PASSWORD")
	if password == "" {
		log.Fatalln("BLUESKY_APP_PASSWORD is not set in the environment.")
		return nil, nil, ""
	}

	var postURL string = "https://bsky.social/xrpc/com.atproto.server.createSession"

	body := BlueskyAuth{
		Identifier: identifier,
		Password:   password,
	}

	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		log.Fatalf("Error parsing body to bytes: %v", err)
		return nil, nil, ""
	}

	reader := bytes.NewReader(bodyBytes)

	resp, err := http.Post(postURL, "application/json", reader)
	if err != nil {
		log.Fatalf("Error while doing POST request: %v", err)
		return nil, nil, ""
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatalf("Error closing response body: %v", err)
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
		return nil, nil, ""
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", resp.StatusCode)
	}

	var authResponse BlueskyAuthResponse

	err = json.Unmarshal([]byte(string(responseBody)), &authResponse)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
		return nil, nil, ""
	}

	return resp, &authResponse, string(responseBody)
}

func Post(content string, authResponse *BlueskyAuthResponse) {
	now := time.Now().UTC().Format(time.RFC3339)
	now = strings.Replace(now, "+00:00", "Z", 1)

	createRecordURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

    session := map[string]string {
        "accessJwt": authResponse.AccessJwt,
        "did": authResponse.Did,
    }

    post := map[string]interface{}{
        "text": content,
        "createdAt": now,
    }

    requestBody, err := json.Marshal(map[string]interface{}{
		"repo":       session["did"],
		"collection": "app.bsky.feed.post",
		"record":     post,
	})
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

    req, err := http.NewRequest("POST", createRecordURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+session["accessJwt"])
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	responseBody := buf.String()
	fmt.Println("Response Body:", responseBody)
}
