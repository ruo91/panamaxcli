package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigPersistence(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{Path: dir + "/agent"}
	err = c.Save("Test Agent", testToken)
	assert.NoError(t, err)

	// To make sure it *really* loaded...
	c.store = Store{}
	err = c.Load()
	assert.NoError(t, err)
	if assert.Len(t, c.Remotes(), 1) {
		r := c.Remotes()[0]
		assert.Equal(t, "Test Agent", r.Name)
		assert.Equal(t, testToken, r.Token)
		assert.Equal(t, "https://45.55.152.201:3001", r.Endpoint)
	}
}

func TestErroredBadTokenPersistence(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{Path: dir + "/agent"}
	err = c.Save("Test Agent", "BAD")
	assert.Contains(t, err.Error(), "illegal base64 data")
}

func TestSuccessfulNonExistantLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{Path: dir + "/agent"}
	err = c.Load()
	assert.NoError(t, err)
	assert.Empty(t, c.Remotes())
}

func TestErroredBadFormatLoad(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)
	c := FileConfig{Path: dir + "/agent"}
	err = ioutil.WriteFile(dir+"/agent", []byte("BAD"), 0600)
	assert.NoError(t, err)

	err = c.Load()
	assert.Contains(t, err.Error(), "invalid character")
}

func TestConfigExists(t *testing.T) {
	c := FileConfig{store: Store{Remotes: []Remote{{Name: "Test"}}}}
	assert.True(t, c.Exists("Test"))
	assert.False(t, c.Exists("BadName"))
}

func TestConfigRemotes(t *testing.T) {
	c := FileConfig{store: Store{Remotes: []Remote{{Name: "Test"}}}}
	if assert.Len(t, c.Remotes(), 1) {
		assert.Equal(t, "Test", c.Remotes()[0].Name)
	}
}

func TestConfigSetActive(t *testing.T) {
	dir, err := ioutil.TempDir("", "agent-test")
	defer os.RemoveAll(dir)
	assert.NoError(t, err)

	c := FileConfig{
		Path: dir + "/agent",
		store: Store{
			Active:  "Test",
			Remotes: []Remote{{Name: "Test"}, {Name: "Test2"}},
		},
	}
	assert.NoError(t, c.SetActive("Test2"))
	// To make sure it really got persisted...
	c.store = Store{}
	assert.NoError(t, c.Load())

	assert.Equal(t, "Test2", c.Active().Name)
}

func TestErroredNonexistantRemoteConfigSetActive(t *testing.T) {
	c := FileConfig{}
	err := c.SetActive("nonexistant")
	assert.EqualError(t, err, "remote 'nonexistant' does not exist")
}

func TestConfigActive(t *testing.T) {
	agent := Remote{Name: "Test"}
	c := FileConfig{store: Store{Remotes: []Remote{agent}}}
	assert.Nil(t, c.Active())

	assert.NoError(t, c.SetActive("Test"))
	assert.Equal(t, &agent, c.Active())
}

func TestRemoteDecodeToken(t *testing.T) {
	remote := Remote{Token: testToken}
	err := remote.DecodeToken()
	assert.NoError(t, err)
	assert.Equal(t, "https://45.55.152.201:3001", remote.Endpoint)
	assert.Equal(t, "d55f5518-b56b-459a-aaa3-2ef7c9241bb7", remote.Username)
	assert.Equal(t, "MmZhMmMyNWEtZmE4ZS00MGM4LWE3Y2ItYTAzNzhjMDVkYzY5Cg==", remote.Password)
	assert.Contains(t, remote.PrivateKey, "BEGIN CERTIFICATE")
}

func TestErroredMissingTokenRemoteDecodeToken(t *testing.T) {
	remote := Remote{Token: ""}
	err := remote.DecodeToken()
	assert.EqualError(t, err, "Missing token")
}

func TestErroredBadTokenRemoteDecodeToken(t *testing.T) {
	remote := Remote{Token: "BAD"}
	err := remote.DecodeToken()
	assert.Contains(t, err.Error(), "illegal base64 data")
}
