#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(readlink -f $(dirname ${BASH_SOURCE[0]}))

echo "export GOPATH=/home/develop/go" >> /home/develop/.zshrc
echo "export GOPATH=/home/develop/go" >> /etc/environment
echo "export PATH=$PATH:/usr/local/go/bin:/home/develop/go/src/github.com/fzerorubigd/balloon/bin" >> /home/develop/.zshrc
echo "alias cdp=\"cd /home/develop/go/src/github.com/fzerorubigd/balloon\"" >> /home/develop/.zshrc

chown -R develop. /home/develop/go
chown -R develop. /home/develop/go/src
chown -R develop. /home/develop/go/src/github.com
chown -R develop. /home/develop/go/src/github.com/fzerorubigd
chown -R develop. /home/develop/go/src/github.com/fzerorubigd/balloon

cd /home/develop/go/src/github.com/fzerorubigd/balloon
make -f /home/develop/go/src/github.com/fzerorubigd/balloon/Makefile database-setup

sudo -u develop /home/develop/go/src/github.com/fzerorubigd/balloon/scripts/provision_user.sh