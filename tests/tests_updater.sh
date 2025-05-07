#!/bin/bash

cd tests || exit 1

for dir in */*/ ; do
    folder_name="${dir%/}"

    if [ -d "$folder_name" ]; then
        echo "Updating test for $folder_name"
        ../docsncode "$folder_name/project" "$folder_name/expected_result"
    fi
done
