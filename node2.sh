#!/bin/sh
set -e

./tidb/bin/tikv-server --pd-endpoints="127.0.0.1:2379" \
                --addr="127.0.0.1:20161" \
                --data-dir=tidb/tikv2 \
                --log-file=tidb/tikv2.log
