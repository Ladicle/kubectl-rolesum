#!/bin/bash -ex

# cleanup output directory
rm -r _output/ || :

# build binary
make build-linux
make build-darwin

# compose binaries
for dir in $(find _output/ -type d -name "kubectl-bindrole_*"); do
  tar -zcvf $(basename $dir).tar.gz $dir
done

# check sha256
sha256sum *.tar.gz
