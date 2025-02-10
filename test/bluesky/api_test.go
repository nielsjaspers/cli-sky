package bluesky_test

import (
	"testing"

	"github.com/nielsjaspers/cli-sky/bluesky"
	"github.com/stretchr/testify/assert"
)

func TestCreateSession(t *testing.T) {
    assert := assert.New(t)

	resp401, _ := bluesky.CreateSession("../../.env.test")
    assert.Equal(401, resp401.StatusCode, "The response should be 401")
	resp200, _ := bluesky.CreateSession("../../.env")
    assert.Equal(401, resp200.StatusCode, "The response should be 401")
}
