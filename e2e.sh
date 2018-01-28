#!/usr/bin/env bash
set -x

# End to end test
netcat -v localhost 3333 <<EOF
GET 1\r
GET $RANDOM\r
GET -$RANDOM\r
GET $RANDOM\r
$RANDOM\r
wgat?
0
QUIT
EOF
