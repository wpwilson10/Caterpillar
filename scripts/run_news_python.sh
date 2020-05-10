#!/bin/bash
# command to save file format correctly in vim
# :set fileformat=unix 
# command to set file exectuate permissions
# sudo chmod u+x run.sh

# go to file location
PRJ_FILE="/home/patrick/Documents/Projects/Caterpillar/"
cd $PRJ_FILE

# run data crawler with given flag
sudo python3 ./cmd/main.py $1
