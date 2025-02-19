#!/bin/bash

## echon -- ensure echo prints values without a newline

echon()
{
  printf "%s" "$*"
  # echo "$*" | tr -d '\n'
}

## RUNNING THE SCRIPT
if [ $# -eq 0 ]; then
  echo "Enter something to echo out."
fi

echon "$1"

exit 0