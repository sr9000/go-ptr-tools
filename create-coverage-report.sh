#!/usr/bin/env bash

# Create an HTML report; store next to script.
go tool cover -html="coverage.out" -o "coverage.html"

# Extract total coverage: the decimal number from the last line of the function report.
COVERAGE=$(go tool cover -func="coverage.out" | tail -1 | grep -Eo '[0-9]+\.[0-9]')

echo "coverage: $COVERAGE% of statements"

date "+%s,$COVERAGE" >> "coverage.log"

# Pick a color for the badge.
if awk "BEGIN {exit !($COVERAGE >= 90)}"; then
	COLOR=brightgreen
elif awk "BEGIN {exit !($COVERAGE >= 80)}"; then
	COLOR=green
elif awk "BEGIN {exit !($COVERAGE >= 70)}"; then
	COLOR=yellowgreen
elif awk "BEGIN {exit !($COVERAGE >= 60)}"; then
	COLOR=yellow
elif awk "BEGIN {exit !($COVERAGE >= 50)}"; then
	COLOR=orange
else
	COLOR=red
fi

# Download the badge; store next to script.
curl -s "https://img.shields.io/badge/coverage-$COVERAGE%25-$COLOR" > "coverage.svg"
