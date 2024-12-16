#!/bin/bash
# nicenum -- Given a number, shows it in comma-separated form.
# expects DD - decimal point delimiter and TD - a thousands delimiter to be instantiated
# instantate function, if a second arg is specified, echo output to stdout

nicenumber()
{
  # assume that '.' is the decimal separator in the inout value to this script
  # the decimal separator in the output is '.', unless specified by the user with the -d flag.

  # get integers - left side of decimal point
  integer=$(echo $1 | cut -d. -f1)
  # get decimal - right side of decimal point
  decimal=$(echo $1 | cut -d. -f2)

  # check if number has more than the integer part
  if [ "$decimal" != "$1" ]; then
    # it means there is a fractional part, so let us include it
    result="${DD:='.'}$decimal"
  fi

  thousands=$integer

  while [ $thousands -gt 999 ]; do
    remainder=$(($thousands % 1000)) # get three significant figures

    # force leading zeros on the remainder
    while [ ${#remainder} -lt 3 ]; do
      remainder="0$remainder"
    done

    # build the integer side from right to left
    result="${TD:=','}${remainder}${result}"
    thousands=$(($thousands / 1000))
  done

  nicenum="${thousands}${result}"
  if [ ! -z $2 ]; then
    echo $nicenum
  fi
}

DD="."
TD=","

# main script
while getopts "d:t" opt; do
  case $opt in
    d ) DD="$OPTARG"  ;;
    t ) TD="$OPTARG"  ;;
  esac
done
shift $(($OPTIND - 1))

# input validation
if [ $# -eq 0 ]; then
  echo "Usage: $(basename $0) [-d c] [-t c] number"
  echo "        -d specifies the decimal point delimiter"
  echo "        -t specifies the thousands delimiter"
  exit 1
fi

nicenumber $1 1

exit 0