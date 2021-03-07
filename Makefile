RUN = docker run --rm --user $$(id -u):$$(id -g)
PROTOC = $(RUN) -v "$$PWD:$$PWD" -w "$$PWD" namely/protoc
PROTOLOCK = $(RUN) -v $$PWD:/protolock -w /protolock nilslice/protolock

EXE = ./ws-worker
SRC_FILES = $(wildcard *.go) \
            $(wildcard */*.go)
PB_FILES = $(patsubst proto/%.proto,proto/%.pb.go,$(wildcard proto/*.proto))

all: $(PB_FILES) lint $(EXE)

$(EXE): $(SRC_FILES)
	go build -o $(EXE)

proto/%.pb.go: proto.lock proto/%.proto
	$(PROTOLOCK) commit
	$(PROTOC) -I=./proto --go_out=plugins=grpc:proto proto/$*.proto

proto.lock:
	$(PROTOLOCK) init

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: run
run: all
	$(EXE) run
