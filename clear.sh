#!/bin/bash

# Directories to clear
directories=("logs" "dumps" "backups" "database")

# Base path
base_path="./data"

# Loop through each directory
for dir in "${directories[@]}"; do
    target_dir="$base_path/$dir"
    
    # Check if the directory exists
    if [ -d "$target_dir" ]; then
        # Find and delete all files except .keep
        find "$target_dir" -type f ! -name '.keep' -delete
        echo "Cleared $target_dir"
    else
        echo "Directory $target_dir does not exist"
    fi
done