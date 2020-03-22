#!/usr/bin/env bash

# abort on nonzero exitstatus
set -o errexit
# abort on unbound variable
set -o nounset
# don't hide errors within pipes
set -o pipefail

#set -x

main() {
    assert_variable_specified D_INVENTORY
    assert_variable_specified D_USER
    assert_variable_specified D_PRIVATE_KEY
    assert_variable_specified D_BINARY_FILE_SOURCE
    assert_variable_specified D_DATABASE_CONNECTION_STRING

    # Disable creating retry files
    export ANSIBLE_RETRY_FILES_ENABLED=0

    echo "Inventory graph: "
    ansible-inventory \
      --inventory=${D_INVENTORY} \
      --yaml \
      --graph \
      --vars

    ansible-playbook -v \
      --inventory=${D_INVENTORY} \
      --user=${D_USER} \
      --ssh-common-args="-o StrictHostKeyChecking=no" \
      --private-key=${D_PRIVATE_KEY} \
      ./scripts/ansible/monitoring-playbook.yml \
      --extra-vars=binary_file_source=${PWD}/${D_BINARY_FILE_SOURCE} \
      --extra-vars=database_connection_string="'${D_DATABASE_CONNECTION_STRING}'"
}

function assert_variable_specified() {
  local variable_name=$1
  if [[ ${!variable_name:-} == "" ]]; then
     echo "${variable_name} environment variable is required"
     exit 4
  fi
}

main "${@}"
