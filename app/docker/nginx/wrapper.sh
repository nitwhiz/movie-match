#!/bin/sh

envsubst < /usr/share/nginx/html/env.template.json > /usr/share/nginx/html/env.json

# start nginx
exec nginx -g "daemon off;"
