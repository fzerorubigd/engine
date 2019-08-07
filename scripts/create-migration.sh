#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))
source ${SCRIPT_DIR}/project.sh

cd ${SCRIPT_DIR}/..

if [[ "$#" -ne 1 ]]; then
    echo "Name for this migration : "
    read NAME
else
    NAME=${1}
fi


DATE=`date +%Y%m%d%H%M%S`
FILE="${DATE}_${NAME}.sql"

cat >>${FILE} <<-EOGO
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

EOGO

echo ${FILE}
