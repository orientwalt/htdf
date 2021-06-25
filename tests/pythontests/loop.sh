#!/bin/sh

# O
# ```shell CentOS7 
# yum update && yum upgrade 
# yum install  install libtool -y
# yum install epel-release -y
# yum install python3-devel -y
# yum install python36 -y
# wget https://ftp.gnu.org/gnu/autoconf/autoconf-latest.tar.gz
# wget https://ftp.gnu.org/gnu/automake/automake-1.16.2.tar.gz
# tar xzf autoconf-latest.tar.gz
# tar xzf automake-1.16.2.tar.gz
# cd autoconf-2.71
# ./configure --prefix=/usr
# make && make install
# cd ..
# cd automake-1.16.2
# ./configure --prefix=/usr
# make && make install
# ```


i=1
while [ ${i} -le 100000000 ]
do
  make test chaintype=regtest
  i=`expr ${i} + 1`
  sleep 180
done

echo done
