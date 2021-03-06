package config

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config interface {
	Save(name string, token string) error
	Exists(name string) bool
	Remotes() []Remote
	SetActive(name string) error
	Active() *Remote
}

type FileConfig struct {
	Path  string
	store Store
}

type Store struct {
	Active  string   `json:"active"`
	Remotes []Remote `json:"remotes"`
}

type Remote struct {
	Name       string `json:"name"`
	Token      string `json:"token"`
	Endpoint   string `json:"endpoint"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
}

func (c *FileConfig) Save(name string, token string) error {
	r := Remote{Name: name, Token: token}
	if err := r.DecodeToken(); err != nil {
		return err
	}

	c.store.Remotes = append(c.Remotes(), r)
	return c.saveAll()
}

func (c *FileConfig) Exists(name string) bool {
	for _, a := range c.Remotes() {
		if a.Name == name {
			return true
		}
	}
	return false
}

func (c *FileConfig) SetActive(name string) error {
	if !c.Exists(name) {
		return fmt.Errorf("remote '%s' does not exist", name)
	}
	c.store.Active = name
	c.saveAll()
	return nil
}

func (c *FileConfig) Active() *Remote {
	activeName := c.store.Active
	if activeName == "" {
		return nil
	}

	for _, r := range c.Remotes() {
		if r.Name == activeName {
			return &r
		}
	}

	return nil
}

func (c *FileConfig) Load() error {
	f, err := os.Open(c.Path)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil
		}

		return err
	}

	d := json.NewDecoder(f)
	if err := d.Decode(&c.store); err != nil {
		return err
	}

	return nil
}

func (c *FileConfig) Remotes() []Remote {
	return c.store.Remotes
}

func (c *FileConfig) saveAll() error {
	b, err := json.MarshalIndent(c.store, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(c.Path, b, 0600)
}

func (r *Remote) DecodeToken() error {
	if r.Token == "" {
		return errors.New("Missing token")
	}
	bs, err := base64.StdEncoding.DecodeString(r.Token)
	if err != nil {
		return err
	}

	data := strings.Split(string(bs), "|")
	r.Endpoint = data[0]
	r.Username = data[1]
	r.Password = data[2]
	r.PrivateKey = data[3]

	return nil
}
