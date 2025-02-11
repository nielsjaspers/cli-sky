package datahandler_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/nielsjaspers/cli-sky/bluesky"
	"github.com/nielsjaspers/cli-sky/internal/datahandler"
	"github.com/stretchr/testify/assert"
)

var handle string = "mock.handle"
var wrongHandle string = "wrong.handle"

var mockAuthResponse bluesky.BlueskyAuthResponse = bluesky.BlueskyAuthResponse{
    AccessJwt: "db70a36a-6548-4580-acd1-2b5141dcd3cb",
    RefreshJwt: "53c7cb54-927f-4547-baa6-57fd2a002294",
    Handle: handle,
    Did: "add9392c-4fa9-40bc-ae90-a2a1b4d8c604",
}

// Test for reading the Auth Response data from a file. A mock auth reponse is used and created on setup(). Tests if the given handle is correct, and then if the rest of the response is correct. If everything goes right, cleanup() is called to remove the mock data
func TestReadAuthResponseFromFile(t *testing.T) {
    assert := assert.New(t)

    setup()
    resp, err := datahandler.ReadAuthResponseFromFile(handle)
    if err != nil {
        log.Fatalf("Error reading auth response from file: %v", err)
    }

    // check if handle is read correct
    handleIsCorrect := assert.Equal(handle, resp.Handle, "Handle should equal \"%v\"", handle)
    // log.Printf("\nhandle Have: %v\t handle Expected: %v", resp.Handle, handle)

    // check if rest of file is read correct
    fileContentIsCorrect := assert.Equal(&mockAuthResponse, resp, "file content should equal \"%v\"", mockAuthResponse)
    // log.Printf("\ncontent Have: %v\t content Expected: %v", resp, &mockAuthResponse)
    
    // log.Printf("handleIsCorrect: %v\tfileContentIsCorrect: %v", handleIsCorrect, fileContentIsCorrect)
    if !handleIsCorrect || !fileContentIsCorrect {
        log.Fatalln("Something went terribly wrong...")
    }

    cleanup()
}

// setup for TestReadAuthResponseFromFile. Writes a mock auth reponse to a temporary file.
func setup() {
    err := datahandler.WriteAuthResponseToFile(&mockAuthResponse)
    if err != nil {
        log.Fatalf("Error writing auth response to file: %v", err)
    }
}

// cleanup after TestReadAuthResponseFromFile was succesful. Removed the temporary file from the api_responses folder.
func cleanup() {
    err := os.Remove(fmt.Sprintf("api_responses/auth_response_%s.json", mockAuthResponse.Handle))
    if err != nil {
        log.Fatalf("Error removing mock response file: %v", err)
    }
}
