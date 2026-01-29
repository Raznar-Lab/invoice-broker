#!/bin/sh
# Load and export every line in .env
export $(grep -v '^#' .env | xargs)

go run . -d
sleep infinity