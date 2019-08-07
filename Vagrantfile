# -*- mode: ruby -*-
# vi: set ft=ruby :

ENV['VAGRANT_DEFAULT_PROVIDER'] = 'docker'

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.username = "develop"
  config.ssh.password = "bita123"

  config.vm.provision "shell" do |s|
    s.path   = "scripts/provision.sh"
    s.args   = ["/home/develop/engine"]
  end

  config.vm.network "forwarded_port", guest: 8090,    	host: 8090     # webserver
  config.vm.network "forwarded_port", guest: 8080,    	host: 8080     # goconvey
  config.vm.network "forwarded_port", guest: 8025,    	host: 8025     # mailhog
  config.vm.network "forwarded_port", guest: 5432,      host: 5432     # postgres
  config.vm.network "forwarded_port", guest: 22,        host: 5555     # ssh server
  config.vm.synced_folder ".", "/home/develop/engine", owner: "develop", group: "develop", create: true

  config.vm.provider "docker" do |d|
    d.build_dir = "./image"
    d.has_ssh = true
    d.cmd = ["/bin/bash", "/home/develop/engine/scripts/init.sh"]
  end
end
