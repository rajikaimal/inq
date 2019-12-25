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

var homeDir, err = os.UserHomeDir()
var inqDir = filepath.Join(homeDir, "inq-notes")

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

func saveLocal(topicType string) error {
	var notePath string

	dt := time.Now()
	formattedDate := dt.Format("01-02-2006")
	home, homeDirErr := os.UserHomeDir()

	if homeDirErr != nil {
		return homeDirErr
	}

	inqDirectory := home + "/inq-notes"

	if topicType == "root" {
		notePath = filepath.Join(inqDirectory, formattedDate)
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
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func stageChanges() error {
	fmt.Println("Staging changes")
	cmd := exec.Command("git", "add", "--all")
	cmd.Dir = inqDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func commitChanges() error {
	fmt.Println("Commiting changes")
	cmd := exec.Command("git", "commit")
	cmd.Dir = inqDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func pushToGitHub() error {
	fmt.Println("Pushing to GitHub")
	cmd := exec.Command("git", "push", "origin", "master")
	cmd.Dir = inqDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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

	app.Commands = []cli.Command{
		{
			Name:    "save",
			Aliases: []string{"a"},
			Usage:   "Saves a note",
			Action: func(c *cli.Context) error {
				saveLocalErr := saveLocal(topicType)
				if saveLocalErr == nil {
					stageChangesErr := stageChanges()
					if stageChangesErr == nil {
						commitChangesErr := commitChanges()
						if commitChangesErr == nil {
							return nil
						} else {
							return commitChangesErr
						}
					} else {
						return stageChangesErr
					}
				} else {
					return saveLocalErr
				}
			},
		},
		{
			Name:    "push",
			Aliases: []string{"a"},
			Usage:   "Push pending changes to GitHub",
			Action: func(c *cli.Context) error {
				pushToGitHubErr := pushToGitHub()
				if pushToGitHubErr == nil {
					fmt.Println("Pushing succesfull")
					return nil
				} else {
					return pushToGitHubErr
				}
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		firstArg := c.Args().Get(0)
		if firstArg == "config" {
			err := runConfigure(githubConfig)
			if err == nil {
				fmt.Println("inq configured successfuly")
			}
		}

		return nil
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
