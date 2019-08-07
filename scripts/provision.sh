#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))
source ${SCRIPT_DIR}/project.sh || export PROJECT=engine

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/${PROJECT}/bin" >> /home/develop/.zshrc
echo "alias cdp=\"cd /home/develop/${PROJECT}\"" >> /home/develop/.zshrc

chown -R develop. /home/develop/${PROJECT}

cd /home/develop/${PROJECT}
make -f /home/develop/${PROJECT}/Makefile database-setup

sudo -u develop /home/develop/${PROJECT}/scripts/provision_user.sh
