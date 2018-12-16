#!/bin/bash -x
set -euo pipefail

echo -e "\nexport ENV=development\n" >> /home/develop/.zshrc
echo -e "\nexport PATH=\${PATH}:/home/develop/balloon/scripts:/home/develop/balloon/bin" >> /home/develop/.zshrc

#make all
