#!/usr/bin/env bash
set -x

# End to end test
nc -n -D 0.0.0.0 3333 <<EOF
GET 1\r
GET $RANDOM\r
GET -$RANDOM\r
GET $RANDOM\r
$RANDOM\r
wgat?
0
EOF
