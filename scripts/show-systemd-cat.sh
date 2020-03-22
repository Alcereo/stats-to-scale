#!/usr/bin/env bash

# abort on nonzero exitstatus
set -o errexit
# abort on unbound variable
set -o nounset
# don't hide errors within pipes
set -o pipefail

#set -x

main() {
    assert_variable_specified D_HOST_IP
    assert_variable_specified D_USER
    assert_variable_specified D_PRIVATE_KEY

    ssh \
    -i ${D_PRIVATE_KEY} \
     ${D_USER}@${D_HOST_IP} \
     sudo systemctl cat stats-to-scale
}

function assert_variable_specified() {
  local variable_name=$1
  if [[ ${!variable_name:-} == "" ]]; then
     echo "${variable_name} environment variable is required"
     exit 4
  fi
}

main "${@}"
