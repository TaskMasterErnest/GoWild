#!/bin/bash
# validAlpaNum -- ensures that input consists only of alphanumeric values

validAlpaNum()
{
  # validate argument: return 0 if upper+lower+digit
  # return 1, if otherwise
  
  # remove all unacceptable characters
  validChars="$(echo $1 | sed -e 's/[^[:alnum:]]//g')"

  if [ -z "$1" ]; then
    echo "Invalid input!"
    exit 1
  elif [ "$validChars" == "$1" ]; then
    return 0
  else
    return 1
  fi
}

# interactive script
/bin/echo -n "Enter input: "
read input

# validate input
if ! validAlpaNum "$input"; then
  echo "Enter only letters and numbers." >&2
  exit 1
else
  echo "Input is valid!"
fi

exit 0