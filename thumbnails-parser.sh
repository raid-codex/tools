#!/bin/bash

champions=$(cat $GOPATH/src/github.com/raid-codex/data/docs/champions/current/index.json)
for thumbnail in $(cat /Users/geoffrey/Documents/raid-codex/champions-thumbnails)
do
    filename=$(echo $thumbnail | cut -d'/' -f5)
    wget $thumbnail -O /Users/geoffrey/Documents/raid-codex/thumbnails/$filename
    champion_name=$(echo $filename | cut -d'.' -f1 | tr '[:upper:]' '[:lower:]' )
    slug=$(echo $champions | jq --arg NAME "$champion_name" -r '.[] | select(.slug == $NAME) | .slug')

    if [[ "$slug" = "" ]]
    then
        echo "slug not found for $champion_name"
    else
        mv /Users/geoffrey/Documents/raid-codex/thumbnails/$filename /Users/geoffrey/Documents/raid-codex/thumbnails/image-champion-small-${slug}.jpg
    fi
done