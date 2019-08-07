#!/bin/bash -x
set -euo pipefail
SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))
source ${SCRIPT_DIR}/project.sh

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/${PROJECT}/scripts:/home/develop/${PROJECT}/bin" >> /home/develop/.zshrc
