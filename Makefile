server_name=staticserver
default: build

build: 
	@rm -f ${server_name}
	@go build -o ${server_name} -gcflags "-N -l"
