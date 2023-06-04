#!/bin/sh
set -e

./tidb/bin/tikv-server --pd-endpoints="127.0.0.1:2379" \
                --addr="127.0.0.1:20160" \
                --data-dir=tidb/tikv1 \
                --log-file=tidb/tikv1.log