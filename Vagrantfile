# -*- mode: ruby -*-
# vi: set ft=ruby :

require "yaml"

Vagrant.require_version ">= 1.8.0"

# install required plugins if necessary
if ARGV[0] == 'up'
    # add required plugins here
    required_plugins = %w( vagrant-rsync-only-changed )
    missing_plugins = []
    required_plugins.each do |plugin|
        missing_plugins.push(plugin) unless Vagrant.has_plugin? plugin
    end

    if ! missing_plugins.empty?
        install_these = missing_plugins.join(' ')
        puts "Found missing plugins: #{install_these}.  Installing..."
        exec "vagrant plugin install #{install_these}; vagrant up"
    end
end

Vagrant.configure(2) do |config|
    config.vm.define "spoton2-hubhotspot" do |config|
        if ENV['PACKAGING_BASE_BOX']
            config.vm.box = "kaorimatz/ubuntu-16.04-amd64"
        else
            config.vm.box = "spoton_07242017_234815"
        end

        config.vm.provider :virtualbox do |vb|
            vb.customize ["modifyvm", :id, "--ioapic", "on"]
            vb.customize ["modifyvm", :id, "--nictype1", "virtio"]

            vb.memory = 2048
            vb.cpus = 2
        end

        config.vm.synced_folder "newtests/data", "/mnt"

        config.vm.synced_folder ".", "/vagrant", type: "rsync",
            rsync__exclude: [
                "*.box",
                "*.swp",
                ".git/",
                "node_modules",
                "spoton.com/node_modules",
                "spoton.com/dist/",
                "hotspot/static/js3/node_modules",
                "hotspot/static/js3/dist/",
                "hub/node_modules",
                "hub/static/js3/dist/",
                "hub/static/css/app.css",
                "newtests/data",
                "newtests/testlogs",
        ].concat(
            `python2.7 #{File.dirname(__FILE__)}/sass.py outfiles`.each_line.map {
                |line| "#{line}".strip
            })


        if not ENV['VAGRANT_NO_TRIGGERS']
            config.vm.network "private_network", ip: "192.168.50.100"
        end

        # Port forwarding for browser sync
        if not ENV['PACKAGING_BASE_BOX'] and not ENV['VAGRANT_NO_TRIGGERS']
            config.vm.network "forwarded_port", guest: 3000, host: 3000, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 3001, host: 3001, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 5000, host: 5000, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 5001, host: 5001, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 8009, host: 8009, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 8000, host: 8000, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 8005, host: 8005, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 27017, host: 27017, protocol: "tcp"
            config.vm.network "forwarded_port", guest: 5432, host: 1234, protocol: "tcp"
        end

        # enable access to mongo from the host for deployment purposes
        if ENV['SPOTON_DEPLOY']
            config.vm.network "forwarded_port", guest: 27017, host: 27017, protocol: "tcp"
        end

        if ENV['PACKAGING_BASE_BOX']
            config.ssh.insert_key = false

            config.vm.provision :shell, :inline => "if ! grep -q $(cat /etc/hostname) /etc/hosts; then echo >> /etc/hosts; echo 127.0.0.1 $(cat /etc/hostname) >> /etc/hosts; fi"

            config.vm.provision "ansible_local" do |ansible|
                ansible.version = "latest"
                ansible.playbook = "vagrant-ansible/galaxy_playbook.yml"
                ansible.sudo = true
                ansible.verbose = "vvvv"
                # 2.3.0.0 came out 12 April and breaks a ton of stuff
                # Lock version for now
                ansible.version = "2.2.2.0"
                ansible.install_mode = "pip"
            end

            config.vm.provision "ansible_local" do |ansible|
                ansible.version = "latest"
                ansible.playbook = "vagrant-ansible/playbook.yml"
                ansible.sudo = true
                ansible.verbose = "vvvv"
                # 2.3.0.0 came out 12 April and breaks a ton of stuff
                # Lock version for now
                ansible.version = "2.2.2.0"
                ansible.install_mode = "pip"
            end
        end

        if Vagrant.has_plugin?("vagrant-gatling-rsync")
            config.gatling.rsync_on_startup = false
        end
    end
end
