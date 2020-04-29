# notes-cli
> A simple CLI app to take notes on markdown files

## Installation
* You can also download compressed version from [releases](https://github.com/prdpx7/notes-cli/releases)
    ```
    wget https://github.com/prdpx7/notes-cli/releases/download/v0.1/notes-2020.04.29.tar.gz
    tar -xzf notes-2020.04.29.tar.gz
    ```
* Or you can download the binary directly
    ```
    wget https://github.com/prdpx7/notes-cli/releases/download/v0.1/notes
    ```
* And copy it into the /usr/local/bin directory
    ```
    sudo cp notes /usr/local/bin/notes
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
