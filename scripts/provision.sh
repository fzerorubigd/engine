#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/balloon/bin" >> /home/develop/.zshrc
echo "alias cdp=\"cd /home/develop/balloon\"" >> /home/develop/.zshrc

chown develop. /home/develop/balloon

cd /home/develop/balloon
make -f /home/develop/balloon/Makefile database-setup

sudo -u develop /home/develop/balloon/scripts/provision_user.sh