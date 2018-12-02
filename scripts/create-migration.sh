#!/usr/bin/env bash
set -euo pipefail

if [[ "$#" -ne 1 ]]; then
    echo "Name for this migration : "
    read NAME
else
    NAME=${1}
fi

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))
cd ${SCRIPT_DIR}/..

DATE=`date +%Y%m%d%H%M%S`
FILE="${DATE}_${NAME}.sql"

cat >>$FILE <<-EOGO
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

EOGO

echo ${FILE}