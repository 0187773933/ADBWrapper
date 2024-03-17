#!/bin/bash

# Function to check if a string is an integer
function is_int() { return $(test "$@" -eq "$@" > /dev/null 2>&1); }

# Clear all loaded SSH keys and re-add your key
ssh-add -D
ssh-add -k /Users/morpheous/.ssh/githubWinStitch

# Initialize the git repository and configure user
git init
git config --global --unset user.name
git config --global --unset user.email
git config user.name "0187773933"
git config user.email "collincerbus@student.olympic.edu"

# Fetch all tags, sort them, and get the highest tag number
git fetch --tags
LastTag=$(git tag | sort -V | tail -n 1)

# Extract the numeric part of the tag, check if it's an integer, and increment it
if [[ "$LastTag" =~ v1\.0\.([0-9]+)$ ]]; then
    LastTagNumber=${BASH_REMATCH[1]}
    if is_int "$LastTagNumber"; then
        NextCommitNumber=$((LastTagNumber + 1))
    else
        echo "Last tag number is not an integer. Resetting to 1."
        NextCommitNumber=1
    fi
else
    echo "No valid last tag found. Starting from 1."
    NextCommitNumber=1
fi

# Add, commit, and tag changes
git add .
git tag -l | xargs git tag -d

if [ -n "$1" ]; then
    git commit -m "$1"
    git tag "v1.0.$1"
else
    git commit -m "$NextCommitNumber"
    git tag "v1.0.$NextCommitNumber"
fi

# Add remote repository and push changes
git remote add origin git@github.com:0187773933/ADBWrapper.git
git push origin --tags
git push origin master