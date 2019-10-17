package gohome

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/h2non/gock.v1"
)

func TestNew(t *testing.T) {
	var (
		testURL   = "http://test.domain"
		testToken = "mySecretToken"
	)

	hc := new(http.Client)

	defer gock.Off()
	gock.InterceptClient(hc)

	gock.New(testURL).
		Get("/api/").
		MatchHeader("Authorization", "^Bearer "+testToken).
		MatchType("json").
		Reply(200).
		JSON(map[string]string{"message": "API running."})

	client, err := New(
		WithClient(hc),
		WithHost(testURL),
		WithPing(),
		WithAuthToken(testToken),
		NoDiscovery(),
	)

	assert.NoError(t, err)
	assert.NotNil(t, client)

	gock.New(testURL).
		Get("/api/").
		Reply(401).
		BodyString("401: Unauthorized")

	_, err = New(
		WithClient(hc),
		WithHost(testURL),
		WithPing(),
		NoDiscovery(),
	)

	assert.Error(t, err)
	assert.Equal(t, "401: Unauthorized", err.Error())

	gock.New(testURL).
		Get("/api/").
		ReplyError(errors.New("Bad Error"))

	_, err = New(
		WithClient(hc),
		WithHost(testURL),
		WithPing(),
		NoDiscovery(),
	)

	assert.Error(t, err)
	assert.Equal(t, "Get http://test.domain/api/: Bad Error", err.Error())
}

func TestDiscovery(t *testing.T) {
	hc := new(http.Client)

	defer gock.Off()
	gock.InterceptClient(hc)

	gock.New("http://hassio.local").
		Get("/api/discovery_info").
		MatchType("json").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		File("testdata/discovery.json")

	client, err := New(
		WithClient(hc),
	)

	assert.NoError(t, err)
	assert.Equal(t, "0.100.2", client.discoveryInfo.Version)
}
