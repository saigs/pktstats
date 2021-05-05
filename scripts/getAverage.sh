#!/bin/bash

#
# average size run
#

ROOT=$HOME/opensource/pktstats
cd $ROOT
echo "Making program.."
make all

# variables
pkts="1000000"   # number of packets
mode="3"         # 0-mutex, 1-atomic, 2-map, 3-concurrent
runxtimes="100"  # run these many times

i="0"
file="/tmp/__avgfile$$__"
echo "Running tests.."
while [ $i -lt $runxtimes ]
do
    $ROOT/build/server -mode $mode -packets $pkts 2>&1 | grep 'RunDatapath .* took' | awk '{ print $(NF-1) }' >> $file
	i=$[$i+1]
done
echo "Average: "
awk '{ s += $1 } END { print s/NR, "us" }' $file
unlink $file
