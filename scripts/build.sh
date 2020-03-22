#!/usr/bin/env bash

# abort on nonzero exitstatus
set -o errexit
# abort on unbound variable
set -o nounset
# don't hide errors within pipes
set -o pipefail

#set -x

main() {
    assert_variable_specified GOPATH
    assert_variable_specified GOROOT

    mkdir -p target

    ${GOROOT}/bin/go build -v -o ./target/stats-to-scale cmd/statsToScale.go
    echo "Binary built to ./target/stats-to-scale"
}

function assert_variable_specified() {
  local variable_name=$1
  if [[ ${!variable_name:-} == "" ]]; then
     echo "${variable_name} environment variable is required"
     exit 4
  fi
}

main "${@}"
