#cleanup 

rm "consul_1.0.0_linux_amd64.zip"
rm "consul"
wget 'https://releases.hashicorp.com/consul/1.0.0/consul_1.0.0_linux_amd64.zip'
unzip "consul_1.0.0_linux_amd64.zip"
./consul --version
./consul agent -dev &
