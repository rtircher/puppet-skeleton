require 'net/ssh'
require 'net/scp'

APP_SERVER_HOST = '...'
APP_SERVER_USER = '...'
APP_SERVER_SUDO_PWD = '...'
KEY_PATH = 'tool/private-key.pem'

namespace :provision do
  desc "Provision Staging Server"
  task :staging do
    Rake::Task['provision:common'].invoke('staging-instance')
  end

  desc "Provision Production Server"
  task :production do
    Rake::Task['provision:common'].invoke('staging-instance')
  end

  task :common, :server_name do |t, args|
    server_name = args[:server_name]
    puts "Ensuring system is up to date"
    sudo('yum -y update')

    puts "Installing Puppet"
    install_puppet_standalone()

    puts "Uploading Puppet manifest"
    upload_puppet_manifest()

    puts "Applying Puppet manifest"
    apply_puppet_manifest(server_name)

    create_connection_script()
  end
end

def install_puppet_standalone
  # TODO try to install puppet using 'yum install -y puppet' first
  sudo('rpm -ivh http://yum.puppetlabs.com/el/6/products/i386/puppetlabs-release-6-6.noarch.rpm')
  sudo('yum install -y puppet')
end

def upload_puppet_manifest
  system('mkdir -p tools/tmp')
  system('tar czf tools/tmp/puppet-manifests.tgz puppet/*')
  run('rm puppet-manifests.tgz')
  run('rm -rf puppet')
  scp('tools/tmp/puppet-manifests.tgz', '.')
  run('tar xzf puppet-manifests.tgz')
end

def apply_puppet_manifest(manifest, dry_run=false)
  noop_flag = dry_run ? "--noop " : ""
  puppet_root = "/home/#{APP_SERVER_USER}/puppet"
  command = "puppet apply #{noop_flag}--modulepath='#{puppet_root}/modules:#{puppet_root}/vendor/modules' #{puppet_root}/manifests/#{manifest}.pp"
  sudo(command)
end

def run(command)
  Net::SSH.start(APP_SERVER_HOST, APP_SERVER_USER, :keys => KEY_PATH) do |ssh|
    ssh.open_channel do |channel|
      channel.request_pty do |c, success|
        if success
          c.exec(command)
        end
      end
      channel.on_data do |channel, data|
        puts data
      end
    end
  end
end

def sudo(command)
  run("echo #{APP_SERVER_SUDO_PWD} | sudo -S #{command}")
end

def scp(source, destination)
  Net::SCP.start(APP_SERVER_HOST, APP_SERVER_USER, :keys => KEY_PATH) do |scp|
    scp.upload!(source, destination)
  end
end

def create_connection_script
  command = "ssh -i #{KEY_PATH} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o LogLevel=quiet #{APP_SERVER_USER}@#{APP_SERVER_HOST}\n"
  server_name = 'app_server'
  file_name = "connect_#{server_name}"

  File.open(file_name, 'w') do |f|
    f.write("#!/bin/bash\n")
    f.write(command)
    f.chmod(0744)
  end

  puts "Connection script #{file_name} created"
  puts "Run ./#{file_name} to connect to the instance"
end
