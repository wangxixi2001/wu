# certificate
Set 'GOPATH' to '/root/go' and pull the project：
```
cd $GOPATH/src && git clone https://github.com/wangxixi2001/wu.git
```
Add in '/etc/hosts'：
```
127.0.0.1  orderer.example.com
127.0.0.1  peer0.org1.example.com
127.0.0.1  peer1.org1.example.com
```
Add Dependency：
```
cd wu && go mod tidy
```
Run Project：
```
./clean_docker.sh
```
Accessing at '127.0.0.1:9000'
