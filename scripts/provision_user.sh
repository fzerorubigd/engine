#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/go/src/github.com/fzerorubigd/balloon/scripts:/home/develop/go/src/github.com/fzerorubigd/balloon/bin" >> /home/develop/.zshrc

#make all
