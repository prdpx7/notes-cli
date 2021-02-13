# notes-cli
> A simple CLI app to take notes on markdown files

## Installation
```
git clone https://github.com/prdpx7/notes-cli.git

cd notes-cli/

go build

chmod +x ./notes-cli

sudo cp ./notes-cli /usr/local/bin/notes
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
