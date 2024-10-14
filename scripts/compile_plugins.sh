#!/bin/bash

# Get the current working directory and resolve its absolute path
rootPath=$(pwd)

# Find all .go files in the 'pl' directory
files=$(find "$rootPath/pl" -name "*.go")

# Loop through each .go file found
for file in $files; do
  # Extract the filename without the extension
  fileName=$(basename "$file")
  fileNameWithoutExt="${fileName%.go}"

  # Define the input and output files
  outputFile="pl/${fileNameWithoutExt}.so"
  inputFile="pl/${fileNameWithoutExt}.go"

  # Build the Go plugin
  if [[ "${ENV}" != "PRODUCTION" ]]; then
    go build -buildmode=plugin -o "$outputFile" "$inputFile"
  else
    go build -buildmode=plugin --trimpath -o "$outputFile" "$inputFile"
  fi

  # Check if the build was successful
  if [ $? -eq 0 ]; then
    echo "Built module: $fileNameWithoutExt"
  else
    echo "Failed to build module: $fileNameWithoutExt"
    exit 1
  fi
done
