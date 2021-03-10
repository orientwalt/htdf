#!/bin/sh

i=1
while [ ${i} -le 100000000 ]
do
  make test chaintype=testnet
  i=`expr ${i} + 1`
  sleep 1
done

echo done
