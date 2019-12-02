package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

type config struct {
	GitHub string
}

func runConfigure(githubConfig string) (err error) {
	if githubConfig == "default" {
		return nil
	}

	home, err := os.UserHomeDir()
	inqDirectory := filepath.Join(home, "inq")

	if err != nil {
		return err
	}

	if _, err := os.Stat(inqDirectory); err == nil {
		// path/to/whatever exists
		fmt.Println("Config directory exists ...")
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		os.MkdirAll(inqDirectory, os.ModePerm)
	}

	fmt.Println("Creating config.json ...")

	configJSON := config{
		GitHub: githubConfig,
	}

	jsonFile, err := json.MarshalIndent(configJSON, "", "    ")

	if err != nil {
		fmt.Println("Error when writing to config file")
	}

	_ = ioutil.WriteFile(inqDirectory+"config.json", jsonFile, 0644)

	cmd := exec.Command("git", "clone", githubConfig)
	cmd.Dir = home
	_, cmdErr := cmd.Output()

	if cmdErr != nil {
		return nil
	}

	return nil
}

func saveOnGitHub() string {
	return "Saving on GitHub"
}

type configuration interface {
	readConfig() string
}

func (c config) readConfig() string {
	//home, err := os.UserHomeDir()
	//inqDirectory := filepath.Join(home, "inq")

	//jsonFile, jsonFileErr := ioutil.ReadFile(inqDirectory)

	//var jsonConfig Config

	//jsonFile.Unmarshal(jsonFile, &jsonConfig)
	// TODO
	return ""
}

func saveLocal(topicType string) {
	var notePath string

	dt := time.Now()
	formattedDate := dt.Format("01-02-2006")
	home, homeDirErr := os.UserHomeDir()

	if homeDirErr != nil {
		return
	}
	inqDirectory := home + ""

	if topicType == "root" {
		notePath = filepath.Join(notePath)
		notePath = notePath + ".md"
	} else {
		notePath = filepath.Join(inqDirectory, topicType, formattedDate)
		notePath = notePath + ".md"
	}

	fmt.Println("Saving note")
	cmd := exec.Command("vim", notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(err)
}

func pushToGitHubByTopic(topic string) {
	fmt.Println("Pushing to GitHub")
	cmd := exec.Command("git push", filepath.Join(topic))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(err)
}

func pushToGitHub() {
	dt := time.Now()
	formattedDate := dt.Format("01-02-2006")

	fmt.Println("Pushing to GitHub")
	cmd := exec.Command("git push", formattedDate+".md")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fmt.Println(err)
}

func main() {
	var githubConfig string
	var topicType string

	app := cli.NewApp()
	app.Name = "inq"
	app.Usage = ""
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "github",
			Value:       "default",
			Usage:       "Configuration type",
			Destination: &githubConfig,
		},
		cli.StringFlag{
			Name:        "topic",
			Value:       "root",
			Usage:       "Topic type",
			Destination: &topicType,
		},
	}

	app.Action = func(c *cli.Context) error {
		firstArg := c.Args().Get(0)
		if firstArg == "config" {
			err := runConfigure(githubConfig)
			if err == nil {
				fmt.Println("inq configured successfuly")
			}
		} else if firstArg == "save" {
			saveLocal(topicType)
		} else if firstArg == "push" {
			pushToGitHub()
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
