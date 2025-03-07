#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path_to_test_go_file>"
    exit 1
fi

# Get the path to the test Go file
test_file=$1

# Check if the file exists
if [ ! -f "$test_file" ]; then
    echo "File not found: $test_file"
    exit 1
fi

# Create the output directory if it doesn't exist
output_dir="./tests_out"
mkdir -p "$output_dir"

# Get the base name of the test file (without the extension)
base_name=$(basename "$test_file" .go)

# Get the directory of the test file
test_dir=$(dirname "$test_file")

# Create the directory for the compiled test file
compiled_dir="$output_dir/$base_name"
mkdir -p "$compiled_dir"

# Compile the Go test file
go test -c -v -o "$compiled_dir/$base_name" "$test_dir"

# Check if the compilation was successful
if [ $? -ne 0 ]; then
    echo "Failed to compile $test_file"
    exit 1
fi

# Run the compiled test file
"$compiled_dir/$base_name"
