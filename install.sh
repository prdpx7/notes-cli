#!/bin/sh

echo "git clone notes-cli"
git clone https://github.com/prdpx7/notes-cli.git

echo "cd notes-cli/"
cd notes-cli/

echo "go build"
go build

echo "chmod +x ./notes-cli"
chmod +x ./notes-cli

echo "Need root privilage to copy notes into /usr/local/bin/"
sudo cp ./notes-cli /usr/local/bin/notes

echo "Done!"
