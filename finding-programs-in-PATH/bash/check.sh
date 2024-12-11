#!/bin/bash
# inpath() -- verifies that a program is valid or can be found in the PATH directory listing

in_path()
{
  # Given a command and a path, try to find command
  # Return 0 if found, 1 if not
  # temporarily modifies the IFS (Internal Field Separator) but restores it afterwards

  cmd=$1
  ourpath=$2
  result=1
  oldIFS=$IFS
  IFS=":"

  for directory in "$ourpath"
  do
    if [ -x $directory/$cmd ]; then
      result=0
    fi
  done

  IFS=$oldIFS
  return $result
}

checkForCmdInPath()
{
  var=$1

  if [ "$var" != "" ]; then
    if [ "${var:0:1}" = "/" ]; then
      if [ ! -x $var ]; then
        return 1
      fi
    elif ! in_path $var "$PATH"; then
      return 2
    fi
  fi
}


# getting user input
if [ $# -ne 1 ]; then
  echo "Usage $0 command" >&2
  exit 1
fi
checkForCmdInPath "$1"
case $? in 
  0 ) echo "$1 found in PATH" ;;
  1 ) echo "$1 not found or not executable" ;;
  2 ) echo "$1 not found in PATH" ;;
esac

exit 0