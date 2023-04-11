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
    server.vm.provision "shell", inline: 'echo "export DB_USER=' + "'" + ENV["DB_USER"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_PASSWORD=' + "'" + ENV["DB_PASSWORD"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_NAME=' + "'" + ENV["DB_NAME"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export DB_PORT=' + "'" + ENV["DB_PORT"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export PORT=' + "'" + ENV["PORT"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export SESSION_KEY=' + "'" + ENV["SESSION_KEY"] + "'" + '" >> ~/.bash_profile'
    server.vm.provision "shell", inline: 'echo "export GIN_MODE=' + "'" + ENV["GIN_MODE"] + "'" + '" >> ~/.bash_profile'

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

      echo "Adding DB_HOST environment variable to bash profile"
      echo "export DB_HOST='$DB_IP'" >> $HOME/.bash_profile

      cp -r /vagrant/* $HOME

      source $HOME/.bash_profile

      echo "Assigning permission to run deploy.sh and env_file.sh"
      chmod +x deploy.sh
      chmod +x env_file.sh

      echo "Building application"
      docker build -t $DOCKER_USERNAME/minitwit:latest --build-arg db_user=$DB_USER --build-arg db_host=$DB_HOST --build-arg db_password=$DB_PASSWORD --build-arg db_name=$DB_NAME --build-arg db_port=$DB_PORT --build-arg port=$PORT --build-arg session_key=$SESSION_KEY --build-arg gin_mode=$GIN_MODE .

      echo "Logging into docker"
      docker login --username $DOCKER_USERNAME --pasword $DOCKER_PASSWORD¨

      echo "Installing loki-docker-driver"
      docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions

      echo "Running docker image..."
      docker run --rm -d -p $PORT:$PORT --name minitwit $DOCKER_USERNAME/minitwit:latest

      echo "================================================================="
      echo "=                            DONE                               ="
      echo "================================================================="
      echo "Navigate in your browser to:"
      THIS_IP=`hostname -I | cut -d" " -f1`
      echo "http://${THIS_IP}:8080"
    SHELL
  end
  
  config.vm.define "monitoring", primary: true do |server|
    server.vm.provider :digital_ocean do |provider|
      provider.ssh_key_name = ENV["SSH_KEY_NAME"]
      provider.token = ENV["DIGITAL_OCEAN_TOKEN"]
      provider.image = 'ubuntu-18-04-x64'
      provider.region = 'fra1'
      provider.size = 's-1vcpu-1gb'
      provider.privatenetworking = true
    end

    server.vm.hostname = "monitoring"

    server.vm.provision "shell", inline: <<-SHELL

      cp -r /vagrant/* $HOME

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
    
      echo "Installing prometheus"
      docker pull prom/prometheus

      echo "Installing grafana"
      docker pull grafana/grafana:4.5.2

      echo "Running prometheus"
      docker run -d -p 9090:9090 -v ./prometheus.yml:/etc/prometheus/prometheus.yml --name prometheus prom/prometheus

      echo "Running grafana"
      docker run -d -p 3000:3000 --name grafana grafana/grafana:4.5.2

   SHELL
  end

  config.vm.provision "shell", privileged: false, inline: <<-SHELL
    sudo apt-get update
  SHELL
end
