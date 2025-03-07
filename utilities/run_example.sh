#!/bin/bash


# Check if the correct number of arguments is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path_to_example_go_file>"
    exit 1
fi

# Get the path to the example Go file
example_file=$1

# Check if the file exists
if [ ! -f "$example_file" ]; then
    echo "File not found: $example_file"
    exit 1
fi

# Create the output directory if it doesn't exist
output_dir="./examples_out"
mkdir -p "$output_dir"

# Get the base name of the example file (without the extension)
base_name=$(basename "$example_file" .go)

# Get the directory name of the example file
dir_name=$(basename $(dirname "$example_file"))

# Create the directory for the compiled file
compiled_dir="$output_dir/$dir_name"
mkdir -p "$compiled_dir"

# Compile the Go file
go build -o "$compiled_dir/$base_name" "$example_file"

# Check if the compilation was successful
if [ $? -ne 0 ]; then
    echo "Failed to compile $example_file"
    exit 1
fi

# Run the compiled file
"$compiled_dir/$base_name"
