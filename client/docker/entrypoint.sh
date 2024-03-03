#!/usr/bin/env sh

VARS=$(env | grep VERTEX_ | grep _ADDR=)                              # Get all env variables starting with VERTEX_ and ending with _ADDR
VARS=$(echo "$VARS" | sed -e 's/VERTEX_\(.*\)_ADDR=\(.*\)/\1:"\2"/g') # Transform VARS from VERTEX_XXX_ADDR=yyy to XXX: yyy
VARS=$(echo "$VARS" | tr '[:upper:]' '[:lower:]')                     # All lowercase

echo "$VARS"

# Write env variables to /usr/share/nginx/html/config.js
echo "window.api_urls = {" > /usr/share/nginx/html/config.js
for VAR in $VARS
do
    echo "  $VAR," >> /usr/share/nginx/html/config.js
done
echo "};" >> /usr/share/nginx/html/config.js

# Start nginx
exec nginx -g 'daemon off;'
