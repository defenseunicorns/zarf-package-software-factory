Vagrant.configure("2") do |config|
    config.vm.provider "virtualbox" do |vb|
      vb.check_guest_additions = false
      vb.cpus = 10
      vb.memory = 20000
    end
  
    config.vm.boot_timeout = 600
  
    config.vm.disk :disk, size: "100GB", primary: true
    config.vm.box = "boxomatic/ubuntu-20.04"
  
    config.vm.hostname = "zarf-examples"
    config.vm.synced_folder '.', '/vagrant', disabled: true
    config.vm.synced_folder './build/', '/build', SharedFoldersEnableSymlinksCreate: false
  
    config.vm.network "forwarded_port", guest: 80, host: 8080
    config.vm.network "forwarded_port", guest: 443, host: 8443
    config.vm.network "forwarded_port", guest: 9080, host: 9080
    config.vm.network "forwarded_port", guest: 9443, host: 9443
  
    config.ssh.insert_key = false
    config.ssh.extra_args = [ "-t", "cd /build; sudo su" ]
  
    config.vm.provision "shell", inline: <<-SHELL
      # The partition is 100GB but the filesystem isn't yet
      growpart /dev/sda 1 && resize2fs /dev/sda1
  
      # Elasticsearch needs this
      sysctl -w vm.max_map_count=262144
  
      # Create a simulated airgap
      echo "0.0.0.0 registry.opensource.zalan.do ghcr.io registry.hub.docker.com hub.docker.com charts.helm.sh repo1.dso.mil github.com registry.dso.mil registry1.dso.mil docker.io index.docker.io auth.docker.io registry-1.docker.io dseasb33srnrn.cloudfront.net production.cloudflare.docker.com" >> /etc/hosts
    SHELL
  end