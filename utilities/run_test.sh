#!/bin/bash

# check if number of arguments  provided is correct
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <path_to_test_go_file>"
    exit 1
fi

# path to the test file
test_file=$1

# check file exists
if [ ! -f "$test_file" ]; then
    echo "file \`$test_file\`not found"
    exit 1
fi

# make output dir if not exist
output_dir="./tests_out"
mkdir -p "$output_dir"

# name of the test file without the extension
base_name=$(basename "$test_file" .go)

# dir of the test file
test_dir=$(dirname "$test_file")

# make dir for the compiled test file
compiled_dir="$output_dir/$base_name"
mkdir -p "$compiled_dir"

# compile test
go test -c -v -o "$compiled_dir/$base_name" "$test_dir"

# check if compilation was succesful
if [ $? -ne 0 ]; then
    echo "failed to compile \`$test_file\`"
    exit 1
fi

# run
"$compiled_dir/$base_name"
