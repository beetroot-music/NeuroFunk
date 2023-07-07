package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v3"
)

var opts struct {
	UserAgent string `short:"a" long:"user-agent" description:"Sets the user agent for the client" default:"NeuroFunk"`
	HostName  string `short:"n" long:"host-name" required:"true" description:"The name of the host given to clients"`
	Relay     string `short:"r" long:"relay" description:"The URL of the relay" default:"http://localhost:8080"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if flags.WroteHelp(err) {
			return
		}
		log.Fatalf("Failed to parse args: %+v", err)
		return
	}

	_, err = loadTestData()
	if err != nil {
		log.Fatalf("Failed to load test data: %+v", err)
		return
	}

	err = createSession()
	if err != nil {
		log.Fatalf("Failed to create session: %+v", err)
		return
	}
}

func createSession() error {
	requestBody := struct {
		UserAgent string
		HostName  string
	}{
		UserAgent: opts.UserAgent,
		HostName:  opts.HostName,
	}

	var responseBody struct {
		ClientURL    string
		SessionToken string
	}

	jsonData, err := json.Marshal(&requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s/player/register", opts.Relay), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("Server responded with '%s'", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		return err
	}

	fmt.Printf("Got session token: %s\n", responseBody.SessionToken)
	fmt.Printf("Got client URL:    %s\n", responseBody.ClientURL)

	return nil
}

func loadTestData() (data TestData, err error) {
	yamlData, err := os.ReadFile("testdata.yml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return
	}

	return
}

type Track struct {
	Index       int
	ID          string
	Title       string
	Album       string
	Artists     []string
	Explicit    int
	ExternalURL string
}

type TestData struct {
	Tracks []Track
}

// afk
