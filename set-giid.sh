#!/bin/bash
set -e

if [[ $# -ne 2 ]]
then
	echo "Usage: $0 champion_name giid"
	exit 1
fi

champion_path="$GOPATH/src/github.com/raid-codex/data/docs/champions/current"

champion_name=$1
giid=$2

champion_slug=$(cat "${champion_path}/index.json" | jq -r --arg CHAMPION_NAME "$champion_name" '.[] | select(.name == $CHAMPION_NAME) | .slug')

if [[ "$champion_slug" = "" ]]
then
	echo "Champion not found"
	exit 1
fi

champion_file="${champion_path}/${champion_slug}.json"
cat $champion_file | jq --arg GIID "$giid" '.giid = $GIID' | sponge $champion_file
cat $champion_file | jq '.thumbnail = ""' | sponge $champion_file
