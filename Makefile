all: get_api
	bash -c "export GOPATH=`pwd` && go install -v ./..."

get_api:
	bash -c "export GOPATH=`pwd` && go get github.com/qiniu/api"
