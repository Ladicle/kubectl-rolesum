#!/bin/bash -ex

# cleanup output directory
rm -r _output/ || :

# build binary
make build-linux
make build-darwin
make build-windows

# compose binaries
cd _output/
for dir in $(find . -type d -name "kubectl-rolesum_*"); do
    cp ../LICENSE $dir
    tar -zcvf $dir.tar.gz $dir
done

# check sha256
sha256sum *.tar.gz
