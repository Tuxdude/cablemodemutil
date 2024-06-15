#!/usr/bin/env bash

set -E -e -o pipefail

ver="${1:?}"
makesystem_dir="${2:?}"
upgrade=0
if [[ "${3}" == "--upgrade" ]]; then
    upgrade=1
fi

makesystem_id="${makesystem_dir:?}/.id"

if [[ "${upgrade:?}" != "1" ]] && [[ -f "${makesystem_id:?}" ]] && [[ "makesystem" == "$(cat ${makesystem_id:?})"  ]]; then
    echo "Makesystem installation already detected, no changes made"
    exit 0
fi

if [[ ${makesystem_dir:?} != "./.makesystem" ]]; then
    echo "Cannot install/upgrade the makesystem to a directory other than \"./.makesystem\", \"${makesystem_dir:?}\" was specified instead"
    exit 1
fi

if [[ "${upgrade:?}" == "1" ]]; then
    echo "Setting up (upgrade) makesystem@v${ver:?} ==> \"${makesystem_dir:?}\""
else
    echo "Setting up (install) makesystem@v${ver:?} ==> \"${makesystem_dir:?}\""
fi

rm -rf "${makesystem_dir:?}"
git clone --quiet --depth 1 --branch v${ver:?} https://github.com/Tuxdude/makesystem.git ${makesystem_dir:?} >/dev/null 2>&1
rm -rf ${makesystem_dir:?}/.git

${makesystem_dir:?}/scripts/post-install.sh

echo "${ver:?}" > .bootstrap/VERSION
exit 0
