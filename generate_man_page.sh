#!/bin/bash -x
set -e

INSTALL_PANDOC="There is no \"pandoc\" executable in your environment, check your PATH variable or pandoc installation."

echo "Starting to build man page from README.md"

if [ $(which pandoc) ]
then
    pandoc -s -t man README.md -o pazuzu.1
else
   echo $INSTALL_PANDOC
fi
