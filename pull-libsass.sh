#!/bin/bash
if [ "$1" = "" ]
then
  echo "Usage: $0 <libsass release tag>"
  exit
fi

git subtree pull --prefix libsass_src https://github.com/sass/libsass.git $1 --squash
