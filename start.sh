#!/bin/sh
set -e

/app/bin/migrate up
exec /app/bin/rss-parser
