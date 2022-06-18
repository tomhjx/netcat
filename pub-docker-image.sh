#!/bin/bash
root=`dirname $0`
tag0=$1

for f in $root/docker/*
do
    if test -d $f
    then
        tag1="${f##*/}"
        tag="$tag0-$tag1"
        echo $f:$tag
        docker build $f -t $tag
        docker push $tag
    fi
done