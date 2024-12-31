#!/bin/bash

# Validate Date -- validates a date, taking into account leap years

## >>add normdate function here<<
normDate=../../03-normalizing-date-formats/bash/check.sh

exceedsDaysInMonth()
{
  # Given a month name and day number, return 0 if the specified day value is less than
  # or equal to the number of days in the month.
  # return 1 if otherwise.

  case $(echo $1 | tr '[:upper:]' '[:lower:]') in
    jan* ) days=31  ;; feb* ) days=28 ;;
    mar* ) days=31  ;; apr* ) days=30 ;;
    may* ) days=31  ;; jun* ) days=30 ;;
    jul* ) days=31  ;; aug* ) days=31 ;;
    sep* ) days=30  ;; oct* ) days=31 ;;
    nov* ) days=30  ;; dec* ) days=31 ;;
    * ) echo "$0: unknown month name $1" >&2
        exit 1
  esac
  
  # check if days inputted is less than or greater than the number of days for month
  if [ $2 -lt 1 -o $2 -gt $days ]; then
    return 1
  else
    return 0 # indicates a valid day
  fi
}

isLeapYear() 
{
  # return 0, if year is a leap year; 1 if otherwise
  # The formula for checking the whether a year is a leap year is:
  # 1. years not divisible by 4 are not leap years
  # 2. years divisible by 4 and 400 are leap years
  # 3. years divisible by 4, and not by 400 byt divisible by 100 are not leap years
  # 4. all other years divisible by 4 are leap years

  year=$1
  
  if [ "$((year % 4))" -ne 0 ]; then
    return 1 # not a leap year
  elif [ "$((year % 400))" -eq 0 ]; then
    return 0 # a leap year
  elif [ "$((year % 100))" -eq 0 ]; then
    return 1
  else
    return 0
  fi
}

## MAIN SCRIPT
if [ $# -ne 3 ]; then
  echo "Usage: $0 month day year" >&2
  echo "Typical input formats are: August 3 1984 or 8 3 1984" >&2
  exit 1
fi

# normalize date and store return values for validation
newDate="$($normDate "$@")"
if [ $? -eq 1 ]; then
  exit 1 # error reported by normDate already
fi

# split the normalized date format
month="$(echo $newDate | cut -d\  -f1)"
day="$(echo $newDate | cut -d\  -f2)"
year="$(echo $newDate | cut -d\  -f3)"

# check if the date is valid
if ! exceedsDaysInMonth $month "$day"; then
  if [ "$month" = "Feb" -a "$day" -eq "29" ]; then
    if ! isLeapYear $year ; then
      echo "$0: $year is not a leap year, so Feb does not have 29 days." >&2
      exit 1
    fi
  else 
    echo "$0: invalid day value: $month does not have $day days." <&2
    exit 1
  fi
fi

echo "Valid date: $newDate"

exit 0