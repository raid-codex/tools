#!/bin/bash

function process
{
    base_dir=$1
    cut_slash=$2

    for file in $(ls ${base_dir}/*.PNG)
    do
        filename=$(basename -- "$file")
        filename="${filename%.*}"
        if [[ ! -d "${base_dir}/${filename}" || -f "${base_dir}/${filename}/image-champion-.jpg" ]]
        then
            true
        else
            champion_slug=$(ls ${base_dir}/$filename/* | grep image-champion |  cut -d'/' -f$cut_slash | cut -d'-' -f3- | cut -d'.' -f1)
            cp $file screenshots/screenshot-champion-${champion_slug}.png
        fi
    done
}

process "." 3
process "./step2" 4
process "./step2/step3" 5
process "./step2/step3/step4" 6
process "./step2/step3/step4/step5" 7
