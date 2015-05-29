#!/usr/bin/env bash
set -ex

if [ -z "$version" ]; then
   echo "version is not set"
   exit 1
fi

for platform in "darwin" "dragonfly" "freebsd" "linux" "netbsd" "openbsd" \
			 "plan9" "solaris" "windows"; do
    for arch in "amd64" "386"; do
	make package version=${version} platform=${platform} arch=${arch}
    done
done

for platform in "freebsd" "linux" "netbsd"; do
    make package version=${version} platform=${platform} arch=arm
done
