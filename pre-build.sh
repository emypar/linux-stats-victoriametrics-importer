#! /bin/bash --noprofile

# Invoked before go build, with the current directory already set to the root of
# the project.

set -e
mkdir -p buildinfo
echo '
// This file was automatically generated by '${0##*/}' at build time. 
// Note that the entire directory is excluded from git control.

package buildinfo

var GitInfo = "'$(git describe --all --dirty=-dirty --long --abbrev=16)'"
var Version = "'$semver'"
' > buildinfo/buildinfo.go

