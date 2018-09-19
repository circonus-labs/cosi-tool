# -*- mode: ruby -*-
# vi: set ft=ruby :
# rubocop:disable Metrics/BlockLength
#
# defines VMs for developing/testing cosi-tool
#

Vagrant.configure('2') do |config|
    config.vm.define 'c7', autostart: false do |c7|
        c7.vm.box = 'maier/centos-7.4.1708-x86_64'
        c7.vm.provider 'virtualbox' do |vb|
            vb.name = 'c7_cosi_tool'
            vb.memory = 2048
            vb.cpus = 1
        end
        c7.vm.synced_folder '.', '/home/vagrant/src/cosi-tool', owner: 'vagrant', group: 'vagrant'
        c7.vm.network 'private_network', ip: '192.168.100.240'
        c7.vm.provision 'shell', inline: <<-SHELL
            [[ -z "$(type -p git)" ]] && {
                echo "Installing needed packages (e.g. git, etc.)"
                yum -q -e 0 makecache fast
                yum -q -e 0 install -y git
            }
            if [[ ! -x /usr/local/go/bin/go ]]; then
                go_ver="1.11"
                echo "Installing go${go_ver}"
                go_tgz="go${go_ver}.linux-amd64.tar.gz"
                [[ -f /vagrant/${go_tgz} ]] || {
                    curl -sSL "https://storage.googleapis.com/golang/${go_tgz}" -o /home/vagrant/$go_tgz
                    [[ $? -eq 0 ]] || { echo "Unable to download go tgz"; exit 1; }
                }
                tar -C /usr/local -xzf /home/vagrant/$go_tgz
                [[ $? -eq 0 ]] || { echo "Error unarchiving $go_tgz"; exit 1; }
            fi
            if [[ ! -x /opt/circonus/agent/sbin/circonus-agentd ]]; then
                agent_ver="0.17.0"
                echo "Installing circonus-agent v${agent_ver}"
                agent_tgz="circonus-agent_${agent_ver}_linux_64-bit.tar.gz"
                [[ -f /vagrant/${agent_tgz} ]] || {
                    curl -sSL "https://github.com/circonus-labs/circonus-agent/releases/download/v${agent_ver}/${agent_tgz}" -o /home/vagrant/$agent_tgz
                    [[ $? -eq 0 ]] || { echo "Unable to download agent tgz"; exit 1; }
                }
                mkdir -p /opt/circonus/agent
                [[ $? -eq 0 ]] || { echo "Error creating /opt/circonus/agent directory"; exit 1; }
                tar -C /opt/circonus/agent -xzf /home/vagrant/$agent_tgz
                [[ $? -eq 0 ]] || { echo "Error unarchiving $agent_tgz"; exit 1; }
            fi
            if [[ ! -x /opt/circonus/cosi-server/sbin/cosi-serverd ]]; then
                server_ver="0.3.0"
                echo "Installing LOCAL cosi-server v${server_ver}"
                server_tgz="cosi-server_${server_ver}_linux_64-bit.tar.gz"
                [[ -f /vagrant/${server_tgz} ]] || {
                    curl -sSL "https://github.com/circonus-labs/cosi-server/releases/download/v${server_ver}/${server_tgz}" -o /home/vagrant/$server_tgz
                    [[ $? -eq 0 ]] || { echo "Unable to download cosi-server tgz"; exit 1; }
                }
                mkdir -p /opt/circonus/cosi-server
                [[ $? -eq 0 ]] || { echo "Error creating /opt/circonus/cosi-server directory"; exit 1; }
                tar -C /opt/circonus/cosi-server -xzf /home/vagrant/$server_tgz
                [[ $? -eq 0 ]] || { echo "Error unarchiving $server_tgz"; exit 1; }
            fi
            [[ -f /etc/profile.d/go.sh ]] || echo 'export PATH="$PATH:/usr/local/go/bin"' > /etc/profile.d/go.sh
            [[ $(grep -c GOPATH /home/vagrant/.bashrc) -eq 0 ]] && {
                mkdir ~vagrant/godev
                chown vagrant.vagrant ~vagrant/godev
                echo 'export GOPATH="${HOME}/godev"' >> /home/vagrant/.bashrc
            }
            exit 0
        SHELL
    end
end
