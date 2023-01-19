#!/bin/bash

# Store a list of folders in $PATH
PATH_FOLDERS=$(echo $PATH | sed 's/:/\n/g')

# Store the owners of each folder in PATH_FOLDERS
PATH_OWNERS=$(echo "$PATH_FOLDERS" | xargs -I {} stat -c "%U" {})

# Get the first occurence of the current use in PATH_OWNERS and store it's line number
FIRST_OWNED_INDEX=$(echo "$PATH_OWNERS" | grep -n "$USER" | head -n 1 | cut -d: -f1)

# Get the folder at the line number FIRST_OWNED_INDEX in PATH_FOLDERS
FIRST_OWNED_PATH_FOLDER=$(echo "$PATH_FOLDERS" | sed -n "$FIRST_OWNED_INDEX"p)

echo "$FIRST_OWNED_PATH_FOLDER"