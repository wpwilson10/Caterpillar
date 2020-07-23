#!/bin/bash
# command to save file format correctly in vim
# :set fileformat=unix 
# command to set file exectuate permissions
# sudo chmod u+x run.sh

#set connection data accordingly
source_host=localhost
source_port=6379
source_db=0
target_host=localhost
target_port=6379
target_db=1
source_password=FtcPsYLE6tQTxNlFqh8SSWVxg4xvJegYtos4dyI2lRV97l+ImqCvo5iqkSRQ1FP+2LE5yMjL12orWBDL

#copy all keys without preserving ttl!
redis-cli -h $source_host -p $source_port -a $source_password -n $source_db keys \* | while read key; do echo "Copying $key"; redis-cli --raw -h $source_host -p $source_port -a $source_password -n $source_db DUMP "$key" | head -c -1|redis-cli -x -h $target_host -p $target_port -a $source_password -n $target_db RESTORE "$key" 0; done
