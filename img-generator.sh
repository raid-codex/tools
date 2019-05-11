#!/bin/bash

if [[ $# -ne 1 ]]
then
    echo "Usage: $0 image-file"
    exit 1
fi

image=$1
filename=$(basename -- "$image")
filename="${filename%.*}"

folder=$filename
cropped_image="${folder}/image-champion-.jpg"
title="${folder}/header.jpg"
title_text="${folder}/text"

rm -rf $folder
mkdir $folder

echo "Generating title..."
convert $image -crop 1334x55+0+0 $title > /dev/null
echo "Generating cropped image..."
convert -crop 600x600+225+55 $image $cropped_image > /dev/null
echo "Extracting text from title..."
tesseract $title $title_text 2> /dev/null > /dev/null

champion_slug=$(name=$(cat "${title_text}.txt" | egrep -E  -o '[A-Z][a-z]{3,}(\s([a-z]{2,})?|Lvl)' | sed 's/Lvl//g' | tr '\n' ' ' | sed 's/  / /g' | xargs ); jq --arg NAME "$name" -r '.[] | select(.name==$NAME) | .image_slug' $GOPATH/src/github.com/raid-codex/champions/export/current/index.json)
if [[ "$champion_slug" != "" ]]
then
    echo "tesseract: Champion is ${champion_slug}"
    mv $cropped_image "${folder}/${champion_slug}.jpg"
    exit 0
fi
output_imgclip=$(imgclip -p --lang eng ${title} 2> /dev/null | grep Result -A1 | tail -n 1 | egrep -E  -o '[A-Z][a-z]{3,}\s([a-z]{2,})?' | tr '\n' ' ' | sed 's/  / /g' | xargs)
echo $output_imgclip > "${title_text}_imgclip.txt"
champion_slug=$(jq --arg NAME "$output_imgclip" -r '.[] | select(.name==$NAME) | .image_slug' $GOPATH/src/github.com/raid-codex/champions/export/current/index.json)
if [[ "$champion_slug" != "" ]]
then
    echo "imgclip: Champion is ${champion_slug}"
    mv $cropped_image "${folder}/${champion_slug}.jpg"
    exit 0
fi