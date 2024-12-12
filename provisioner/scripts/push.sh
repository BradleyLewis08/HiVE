#!/bin/bash

# Check if a commit message was provided
if [ $# -eq 0 ]; then
    echo "Error: Please provide a commit message."
    echo "Usage: $0 <commit_message>"
    exit 1
fi

# Remove the 'server' binary if it exists
if [ -f "server" ]; then
    rm server
    echo "Removed 'server' binary."
else
    echo "'server' binary not found. Skipping removal."
fi

# Add all changes to git
git add .

# Commit changes with the provided message
git commit -m "$1"

# Push changes to the remote repository
git push

echo "Git update completed successfully."