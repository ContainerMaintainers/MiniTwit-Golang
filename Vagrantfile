# -*- mode: ruby -*-
# vi: set ft=ruby :

$ip_file = "db_ip.txt"

Vagrant.configure("2") do |config|
  config.vm.box = 'digital_ocean'
  config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
  config.ssh.private_key_path = '~/.ssh/id_rsa'
  config.vm.synced_folder ".", "/vagrant", type: "rsync"

  config.vm.define "dbserver", primary: true do |server|
    server.vm.provider :digital_ocean do |provider|
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      provider.image = 'ubuntu-18-04-x64'
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'
      provider.privatenetworking = true
    end

    server.vm.hostname = "dbserver"

    server.trigger.after :up do |trigger|
      trigger.info =  "Writing dbserver's IP to file..."
      trigger.ruby do |env,machine|
        remote_ip = machine.instance_variable_get(:@communicator).instance_variable_get(:@connection_ssh_info)[:host]
        File.write($ip_file, remote_ip)
      end
    end

    server.vm.provision "shell", inline: <<-SHELL
      echo "Installing postgresql"
      sudo apt-get -y install wget ca-certificates
      wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
      sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" >> /etc/apt/sources.list.d/pgdg.list'
      sudo apt-get -y update
      sudo apt-get -y install postgresql-15 postgresql-contrib-15
      echo "Starting postgresql"
      service postgresql start
      echo "Checking status of postgresql"
      service postgresql status
      echo "Updating config files"
      echo "listen_addresses='*'" >> /etc/postgresql/15/main/postgresql.conf
      echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/15/main/pg_hba.conf
      sudo -u postgres psql -c "CREATE DATABASE minitwitdb;"
      sudo -u postgres psql -c "CREATE USER admin WITH ENCRYPTED PASSWORD 'admin';"
      sudo -u postgres psql -c "ALTER DATABASE minitwitdb OWNER TO admin;"
      sudo -u postgres psql -c "\\l"
      sudo service postgresql restart
      sudo -u postgres psql -c "show listen_addresses;"
   SHELL
  end

  config.vm.define "webserver", primary: false do |server|

    server.vm.provider :digital_ocean do |provider|
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      provider.image = 'ubuntu-18-04-x64'
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'
      provider.privatenetworking = true
    end

    server.vm.hostname = "webserver"

    server.trigger.before :up do |trigger|
      trigger.info =  "Waiting to create server until dbserver's IP is available."
      trigger.ruby do |env,machine|
        ip_file = "db_ip.txt"
        while !File.file?($ip_file) do
          sleep(1)
        end
        db_ip = File.read($ip_file).strip()
        puts "Now, I have it..."
        puts db_ip
      end
    end

    server.trigger.after :provision do |trigger|
      trigger.ruby do |env,machine|
        File.delete($ip_file) if File.exists? $ip_file
      end
    end

    server.vm.provision "shell", inline: 'echo "export DOCKER_USERNAME=' + "'" + ENV["DOCKER_USERNAME"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DOCKER_PASSWORD=' + "'" + ENV["DOCKER_PASSWORD"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_PORT=' + "'" + ENV["DB_PORT"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_USER=' + "'" + ENV["DB_USER"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_PASSWORD=' + "'" + ENV["DB_PASSWORD"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_NAME=' + "'" + ENV["DB_NAME"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export PORT=' + "'" + ENV["PORT"] + "'" + '" >> ~/.bash_profile'

    server.vm.provision "shell", inline: <<-SHELL
      export DB_IP=`cat /vagrant/db_ip.txt`
      echo $DB_IP

      echo "Installing docker..."
      sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
      curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
      sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
      apt-cache policy docker-ce
      sudo apt install -y docker-ce
      sudo systemctl status docker
      sudo usermod -aG docker ${USER}

      echo "Verifying that docker works"
      docker run --rm hello-world
      docker rmi hello-world

      echo "Adding environment variables to bash profile"
      echo ". $HOME/.bashrc" >> $HOME/.bash_profile

      echo "Adding DB_IP environment variable to bash profile"
      echo "export DB_IP='$DB_IP'" >> $HOME/.bash_profile

      cp -r /vagrant/* $HOME

      source $HOME/.bash_profile

      echo "Assigning permission to run deploy.sh"
      chmod +x deploy.sh

      echo "Deploying docker image..."
      ./deploy.sh

      echo "================================================================="
      echo "=                            DONE                               ="
      echo "================================================================="
      echo "Navigate in your browser to:"
      THIS_IP=`hostname -I | cut -d" " -f1`
      echo "http://${THIS_IP}:8080"
    SHELL
  end
  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    sudo apt-get update
  SHELL
end