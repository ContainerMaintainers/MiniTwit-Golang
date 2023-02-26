# -*- mode: ruby -*-
# vi: set ft=ruby :

unless File.exist? "/etc/vbox/networks.conf"
  # See https://www.virtualbox.org/manual/ch06.html#network_hostonly
 puts "Adding network configuration for VirtualBox."
 puts "You will need to enter your root password..."
 system("sudo bash vbox-network.sh")
end

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "generic/ubuntu1804"

  config.vm.network "private_network", type: "dhcp"

  config.vm.synced_folder ".", "/vagrant", type: "rsync"


  config.vm.define "dbserver", primary: true do |server|
    server.vm.network "private_network", ip: "192.168.20.2"
    server.vm.provider "virtualbox" do |vb|
      vb.memory = "1024"
    end
    server.vm.hostname = "dbserver"
    server.vm.provision "shell", privileged: false, inline: <<-SHELL
        echo "Installing Postgres"
        sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
        wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
        sudo apt-get update
        sudo apt-get -y install postgresql
        sudo mkdir -p /data/db
        conf_location=`sudo -u postgres psql -c 'SHOW config_file'`
        sudo sed -i '/  bindIp:/ s/127.0.0.1/0.0.0.0/' $conf_location
        postgres -D /data/db
        mongorestore --gzip /vagrant/dump
    SHELL
  end

  config.vm.define "webserver", primary: true do |server|
    server.vm.network "private_network", ip: "192.168.20.3"
    # server.vm.network "forwarded_port", guest: 5000, host: 5000
    server.vm.provider "virtualbox" do |vb|
      vb.memory = "1024"
    end
    server.vm.hostname = "webserver"
    server.vm.provision "shell", privileged: false, inline: <<-SHELL
        export DB_IP="192.168.20.2"
        cp -r /vagrant/* $HOME
        sudo apt update
        sudo apt upgrade
        sudo apt install golang-go
        go mod download
        nohup go build -o . minitwit.go
        IP=$(ifconfig eth2 | awk -F ' *|:' '/inet /{print $3}')
        echo "================================================================="
        echo "=                            DONE                               ="
        echo "================================================================="
        echo "Navigate in your browser to: http://$IP:5000"
    SHELL
  end

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
      sudo apt-get update
    SHELL
end
