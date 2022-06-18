#!/bin/bash

b='7.4.9'

tag=$(git describe --tags `git rev-list --tags="7.4.9.*" --max-count=1`)

echo $tag

array=(${tag//./ }) 

echo ${array[3]}