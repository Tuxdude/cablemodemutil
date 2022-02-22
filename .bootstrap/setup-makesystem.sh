#!/usr/bin/env bash

set -e -o pipefail

makesystem_dir="${1:?}"
makesystem_id="${makesystem_dir:?}/.id"

if [[ -f "${makesystem_id:?}" ]] && [[ "makesystem" == "$(cat ${makesystem_id:?})"  ]]; then
    echo "Makesystem installation already detected, no changes made"
    exit 0
fi

if [[ ${makesystem_dir:?} != "./.makesystem" ]]; then
    echo "Cannot install the makesystem to a directory other than \"./.makesystem\", \"${makesystem_dir:?}\" was specified instead"
    exit 1
fi

script_parent_dir="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
ver=$(cat ${script_parent_dir:?}/VERSION)

echo "Setting up makesystem@v${ver:?} ==> \"${1:?}\""
rm -rf "${makesystem_dir:?}"
git clone --quiet --depth 1 --branch v${ver:?} https://github.com/Tuxdude/makesystem.git ${makesystem_dir:?} >/dev/null 2>&1
rm -rf ${makesystem_dir:?}/.git
exit 0
