#!/usr/bin/env bash

git clone https://github.com/vertex-center/openapi.git vertex-openapi

find vertex-openapi/ -type f -name 'openapi.*.yaml' | while read file; do

    filename=$(basename -- "$file")
    dirname=$(echo "$filename" | sed -e 's/openapi.//' -e 's/.yaml//')

    mkdir -p "api/$dirname/next"
    cp "$file" "api/$dirname/next/openapi.yml"

    echo "Moved $filename to api/$dirname/next/openapi.yml"

done

rm -rf vertex-openapi
