package bluesky

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func RefreshSession(refreshJwt string) (*BlueskyAuthResponse, error) {
	url := "https://bsky.social/xrpc/com.atproto.server.refreshSession"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+refreshJwt)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh session failed: %s, %s", resp.Status, string(body))
	}

	// Parse the JSON response
	var refreshResponse BlueskyAuthResponse
	err = json.Unmarshal(body, &refreshResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &refreshResponse, nil
}

func Post(content string, authResponse *BlueskyAuthResponse) {
	now := time.Now().UTC().Format(time.RFC3339)
	now = strings.Replace(now, "+00:00", "Z", 1)

	createRecordURL := "https://bsky.social/xrpc/com.atproto.repo.createRecord"

	session := map[string]string{
		"accessJwt": authResponse.AccessJwt,
		"did":       authResponse.Did,
	}

	facets := parseFacets(content)

	postRecord := PostRecord{
		Text:      content,
		CreatedAt: now,
		Facets:    facets,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"repo":       session["did"],
		"collection": "app.bsky.feed.post",
		"record":     postRecord,
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

func parseFacets(text string) []Facet {
	var facets []Facet
	mentions := parseMentions(text)
	urls := parseURLs(text)
	tags := parseTags(text)

	for _, m := range mentions {
		did, err := resolveHandle(m["handle"])
		if err != nil {
			log.Printf("Error resolving handle %s: %v\n", m["handle"], err)
			continue // Skip this mention if handle resolution fails
		}

		start, err := strconv.Atoi(m["start"])
		if err != nil {
			log.Printf("Error converting start index to integer: %v\n", err)
			continue // Skip this mention if start index conversion fails
		}

		end, err := strconv.Atoi(m["end"])
		if err != nil {
			log.Printf("Error converting end index to integer: %v\n", err)
			continue // Skip this mention if end index conversion fails
		}

		facets = append(facets, Facet{
			Index: FacetIndex{
				ByteStart: start,
				ByteEnd:   end,
			},
			Features: []FacetFeature{{
				Type: "app.bsky.richtext.facet#mention",
				Did:  did,
			}},
		})
	}

	for _, u := range urls {
		start, err := strconv.Atoi(u["start"])
		if err != nil {
			log.Printf("Error converting start index to integer: %v\n", err)
			continue // Skip this URL if start index conversion fails
		}

		end, err := strconv.Atoi(u["end"])
		if err != nil {
			log.Printf("Error converting end index to integer: %v\n", err)
			continue // Skip this URL if end index conversion fails
		}

		facets = append(facets, Facet{
			Index: FacetIndex{
				ByteStart: start,
				ByteEnd:   end,
			},
			Features: []FacetFeature{{
				Type: "app.bsky.richtext.facet#link",
				URI:  u["url"],
			}},
		})
	}

	for _, t := range tags {
		start, err := strconv.Atoi(t["start"])
		if err != nil {
			log.Printf("Error converting start index to integer: %v\n", err)
			continue // Skip this tag if start index conversion fails
		}

		end, err := strconv.Atoi(t["end"])
		if err != nil {
			log.Printf("Error converting end index to integer: %v\n", err)
			continue // Skip this tag if end index conversion fails
		}

		facets = append(facets, Facet{
			Index: FacetIndex{
				ByteStart: start,
				ByteEnd:   end,
			},
			Features: []FacetFeature{{
				Type: "app.bsky.richtext.facet#tag",
				Tag:  t["tag"],
			}},
		})
	}

	return facets
}

func parseMentions(text string) []map[string]string {
	var spans []map[string]string
	mentionRegex := regexp.MustCompile(`(?m)(?:^|\W)(@([a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)`)
	textBytes := []byte(text)

	matches := mentionRegex.FindAllSubmatchIndex(textBytes, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			start := match[2]
			end := match[3]
			handle := string(textBytes[match[2]:match[3]])

			// Remove the leading "@" symbol from the handle
			if len(handle) > 0 && handle[0] == '@' {
				handle = handle[1:]
			}

			spans = append(spans, map[string]string{
				"start":  strconv.Itoa(start),
				"end":    strconv.Itoa(end),
				"handle": handle,
			})
		}
	}
	return spans
}

func parseURLs(text string) []map[string]string {
	var spans []map[string]string
	urlRegex := regexp.MustCompile(`(?m)(?:^|\W)((https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*)))`)
	textBytes := []byte(text)
	textLength := len(textBytes)

	matches := urlRegex.FindAllSubmatchIndex(textBytes, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			start := match[2]
			end := match[3]

			if start >= 0 && end <= textLength && start < end {
				urlStr := string(textBytes[start:end])

				// Validate the URL
				_, err := url.ParseRequestURI(urlStr)
				if err != nil {
					log.Printf("Invalid URL: %s, error: %v\n", urlStr, err)
					continue // Skip invalid URLs
				}

				spans = append(spans, map[string]string{
					"start": strconv.Itoa(start),
					"end":   strconv.Itoa(end),
					"url":   urlStr,
				})
			} else {
				log.Printf("Invalid URL indices: start=%d, end=%d, textLength=%d\n", start, end, textLength)
			}
		}
	}
	return spans
}

func parseTags(text string) []map[string]string {
	var spans []map[string]string
	tagRegex := regexp.MustCompile(`(?m)(#([a-zA-Z0-9_]+))`)
	textBytes := []byte(text)
	textLength := len(textBytes)

	matches := tagRegex.FindAllSubmatchIndex(textBytes, -1)

	for _, match := range matches {
		if len(match) >= 4 {
			start := match[2]
			end := match[3]

			if start >= 0 && end <= textLength && start <= end {
				tag := string(textBytes[start:end])

				spans = append(spans, map[string]string{
					"start": strconv.Itoa(start),
					"end":   strconv.Itoa(end),
					"tag":   tag,
				})
			} else {
				log.Printf("Invalid tag indices: start=%d, end=%d, textLength=%d\n", start, end, textLength)
			}
		}
	}
	return spans
}

func resolveHandle(handle string) (string, error) {
	resolveURL := "https://bsky.social/xrpc/com.atproto.identity.resolveHandle?handle=" + url.QueryEscape(handle)

	resp, err := http.Get(resolveURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("handle resolution failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var resolveResponse map[string]interface{}
	err = json.Unmarshal(body, &resolveResponse)
	if err != nil {
		return "", err
	}

	did, ok := resolveResponse["did"].(string)
	if !ok {
		return "", fmt.Errorf("DID not found in response")
	}

	return did, nil
}

