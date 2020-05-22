#!/usr/bin/env bash
set -euo pipefail

source /bd_build/buildconfig

curl https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
echo "deb http://apt.postgresql.org/pub/repos/apt/ `lsb_release -cs`-pgdg main" >> /etc/apt/sources.list.d/pgdg.list

apt update
$minimal_apt_get_install wget sudo git zsh nano libsqlite3-dev autoconf bison build-essential libssl-dev \
                libyaml-dev libreadline6 libreadline6-dev zlib1g zlib1g-dev htop redis-server postgresql postgresql-contrib \
                mercurial ruby-dev realpath pkg-config unzip dnsutils re2c python3-pip openjdk-8-jre jq \
                python3-dev libpq-dev tmux bzr libsodium-dev cmake python3-setuptools iputils-ping iproute2

# Install golang
cd /tmp
curl https://godeb.s3.amazonaws.com/godeb-amd64.tar.gz | gunzip | tar xvf -
/tmp/godeb install 1.14.3

GOBIN=/usr/local/bin GOPATH=/tmp go get -v -u github.com/mailhog/MailHog
pip3 install --no-cache-dir psycopg2-binary pgcli

# Create vagrant user
bash -c "echo root:bita123 | chpasswd"
groupadd develop
useradd -g develop -m -s /usr/bin/zsh develop
bash -c "echo %develop ALL=NOPASSWD:ALL > /etc/sudoers.d/vagrant"
chmod 0440 /etc/sudoers.d/vagrant
bash -c "echo develop:bita123 | chpasswd"
sudo -Hu develop -- wget -O /home/develop/.zshrc http://git.grml.org/f/grml-etc-core/etc/zsh/zshrc
sudo -Hu develop -- wget -O /home/develop/.zshrc.local  http://git.grml.org/f/grml-etc-core/etc/skel/.zshrc

# Go related stuff
bash -c "echo export GOPATH=/home/develop >> /etc/environment"
sudo -Hu develop -- bash -c "echo export GOPATH=/home/develop >> /home/develop/.zshrc.local"
sudo -Hu develop -- bash -c "echo export PATH=$PATH:/usr/local/go/bin:/home/develop/bin >> /home/develop/.zshrc.local"
