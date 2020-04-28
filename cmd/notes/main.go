package main

import (
	"fmt"
	"os"

	"github.com/prdpx7/notes-cli"
)

func showUsage() {
	helpMessage := `notes - A simple note-taking app
Usage: notes [OPTIONS]
Options:
	write -  write in markdown file created with today's date stamp
	read - browse all notes in terminal file browser
	sync - upload all markdown files to your github  as Private Gist
Example:
notes write
notes read
notes sync
`
	fmt.Println(helpMessage)
}

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}
	mode := os.Args[1]
	editor, exists := os.LookupEnv("EDITOR")
	if exists == false {
		fmt.Println("EDITOR not set in your enviornment!")
		fmt.Println("edit your env(~/.bashrc etc) and write export EDITOR='vim'")
	}
	if mode == "write" {
		cmd := notes.GetEditorCommand(editor, mode)
		notes.RunEditor(cmd)
	} else if mode == "read" {
		editor := notes.GetWorkingTextEditorWithFileBrowsingSupport()
		cmd := notes.GetEditorCommand(editor, mode)
		notes.RunEditor(cmd)
	} else if mode == "sync" {
		notes.DoCloudSync()
	} else {
		showUsage()
	}
}
