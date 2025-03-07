#!/bin/bash

#ccheck if number of arguments  provided is correct
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path_to_example_go_file>"
    exit 1
fi

# path to the example file
example_file=$1

# check if file exits
if [ ! -f "$example_file" ]; then
    echo "file \`$example_file\` not found"
    exit 1
fi

# make output directory if it doesnt exist
output_dir="./examples_out"
mkdir -p "$output_dir"

# get the name of the example file without extension .go
base_name=$(basename "$example_file" .go)

# dir name of the example file
dir_name=$(basename $(dirname "$example_file"))

# make dir for the compiled file
compiled_dir="$output_dir/$dir_name"
mkdir -p "$compiled_dir"

# compile
go build -o "$compiled_dir/$base_name" "$example_file"

# check if compilation was successful
if [ $? -ne 0 ]; then
    echo "failed to compile $example_file"
    exit 1
fi

# run
"$compiled_dir/$base_name"
