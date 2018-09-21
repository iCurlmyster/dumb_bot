
PLUGIN=print.so

all: plugins
	go build

plugins: $(PLUGIN)

print.so:
	go build -buildmode=plugin -o plugins/print/print.so plugins/print/print.go
