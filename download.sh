#!/bin/sh
set -e

# Download the package.
export version=v5.0.1
wget https://download.pingcap.org/tidb-${version}-linux-amd64.tar.gz
wget http://download.pingcap.org/tidb-${version}-linux-amd64.sha256

# Check the file integrity. If the result is OK, the file is correct.
sha256sum -c tidb-${version}-linux-amd64.sha256

# Extract the package.
tar -xzf tidb-${version}-linux-amd64.tar.gz
mv tidb-${version}-linux-amd64 tidb
