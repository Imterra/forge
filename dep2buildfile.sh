#!/bin/bash

usage() {
  printf "Usage: %s DEPFILE\n" $0 >&2
}

if [[ $# -lt 1 ]]; then
  usage
  exit 1
fi

depfile=$1

cat $1 | paste -d',' -s - | sed 's/\\, //g' | tr ',' '\n' | sed 's/.o: \([^ ]*\) /:\1:/g' | \
  awk 'BEGIN{FS=":"} {res = $3; gsub(/ /, " , ", res); printf "%s:\n  type: lib_c\n  sources: [ %s ]\n  resources: [ %s ]\n\n",$1,$2,res}'
