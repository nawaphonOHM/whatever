#!/bin/bash

set -e

TEMP_JSON_FILE=$(mktemp || exit 1)
trap 'rm --force -- "$TEMP_JSON_FILE"' EXIT

TEMP_ERROR=$(mktemp || exit 1)
echo 0 > "$TEMP_ERROR"

TEMP_LOCK=$(mktemp || exit 1)
exec 1996> "$TEMP_LOCK" || exit 1

function lock() {
  flock --exclusive 1996
}

function unlock() {
  flock --unlock 1996
}

function search() {
  local GO_FILE="$1"

  local OUTPUT='null'
  OUTPUT=$(grep --line-number --with-filename --extended-regexp --regexp="//nolint:.+" "$GO_FILE")

  if [ "$OUTPUT" != 'null' ]; then
    echo 1 > "$TEMP_ERROR"
  fi

  OUTPUT=$(sed --regexp-extended 's/([^:]+):([0-9]+):(.+)/\1|\2|\3/' <<< "$OUTPUT")

  lock

  jq --arg RAW_JSON "$OUTPUT" '
    ($RAW_JSON | split("|")) as [$location, $lineNumber, $code] |
    . += [{ location: $location, lineNumber: ($lineNumber | tonumber), code: $code }]
  ' "$TEMP_JSON_FILE" > temp_json_obj.json || exit 1
  mv temp_json_obj.json "$TEMP_JSON_FILE" || exit 1

  unlock
}

mapfile -t GO_FILES <<< "$(find . -name '*.go')"

lock
printf '[]' > "$TEMP_JSON_FILE"
unlock

for GO_FILE in "${GO_FILES[@]}"; do
  search "$GO_FILE" &
done

wait

HAS_ERROR=$(cat "$TEMP_ERROR")

if [ "$HAS_ERROR" -eq 1 ]; then
  printf "Foud Lint Rule Bypassing: \n" >&2
  jq '.' "$TEMP_JSON_FILE" >&2
  exit 1
else
  printf "Not Found Lint Rule Bypassing\n" >&1
  exit 0
fi
