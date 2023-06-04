#!/bin/sh
set -e
./tidb/bin/pd-server --name=pd1 \
                --data-dir=tidb/pd1 \
                --client-urls="http://127.0.0.1:2379" \
                --peer-urls="http://127.0.0.1:2380" \
                --initial-cluster="pd1=http://127.0.0.1:2380" \
                --log-file=tidb/d1.log
