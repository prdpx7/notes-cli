# notes-cli
> A simple CLI app to take notes on markdown files

## Installation
Install go, then
```
% go install github.com/prdpx7/notes.cli/cmd/notes
```

## Usage
```
notes - A simple note-taking app
Usage: notes [OPTIONS]
Options:
	write -  write in markdown file automatically created with today's date stamp
	read - browse all notes in terminal file browser
	sync - upload all markdown files to your github as Private Gist
Example:
notes write
notes read
notes sync
```
<img src ="./notes_demo.gif">
