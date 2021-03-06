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

func checkDir(dir string) error {
	//TODO: fix err return
	if _, err := os.Stat(dir); err == nil {
		// path/to/whatever exists
		fmt.Println(dir)
		fmt.Println("Dir exists")
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		os.MkdirAll(dir, os.ModePerm)
		fmt.Println("Dir does not exist")
	}

	return nil
}

func runConfigure(githubConfig string) (err error) {
	fmt.Println("Running config")
	home, err := os.UserHomeDir()
	inqDirectory := filepath.Join(home, "inq")

	if err != nil {
		return err
	}

	//create sub dir if it doesn't exist
	dirCheckErr := checkDir(inqDirectory)

	if dirCheckErr != nil {
		return dirCheckErr
	}

	fmt.Println("Creating config.json ...")

	configJSON := config{
		GitHub: githubConfig,
	}

	jsonFile, err := json.MarshalIndent(configJSON, "", "    ")

	if err != nil {
		fmt.Println("Error when writing to config file")
	}

	fmt.Println(inqDirectory + "/config.json")

	fileWriteErr := ioutil.WriteFile(inqDirectory+"/config.json", jsonFile, 0644)

	if fileWriteErr != nil {
		fmt.Println(fileWriteErr)
		return fileWriteErr
	}

	cmd := exec.Command("git", "clone", githubConfig)
	cmd.Dir = home
	_, cmdErr := cmd.Output()

	if cmdErr != nil {
		fmt.Println(cmdErr)
		return cmdErr
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
		dirPath := filepath.Join(inqDirectory, topicType)
		//create sub dir if it doesn't exist
		dirCheckErr := checkDir(dirPath)

		if dirCheckErr != nil {
			return dirCheckErr
		}

		notePath = filepath.Join(dirPath, formattedDate)
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

	app.Commands = []cli.Command{
		{
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   "Save a note",
			Action: func(c *cli.Context) error {
				saveLocalErr := saveLocal(topicType)
				if saveLocalErr != nil {
					return saveLocalErr
				}

				stageChangesErr := stageChanges()

				if stageChangesErr != nil {
					return stageChangesErr
				}

				commitChangesErr := commitChanges()

				if commitChanges != nil {
					return commitChangesErr
				}

				return nil
			},
		},
		{
			Name:    "push",
			Aliases: []string{"p"},
			Usage:   "Push pending changes to GitHub",
			Action: func(c *cli.Context) error {
				pushToGitHubErr := pushToGitHub()
				if pushToGitHub != nil {
					return pushToGitHubErr
				}

				fmt.Println("Pushing succesfull")
				return nil
			},
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Configure github repository",
			Action: func(c *cli.Context) error {
				githubConfig := c.Args().Get(0)
				err := runConfigure(githubConfig)
				if err != nil {
					return err
				}

				fmt.Println("inq configured successfuly")
				return nil
			},
		},
	}

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

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
