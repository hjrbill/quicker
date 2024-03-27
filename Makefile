REVISION := $(shell git rev-parse --short HEAD)
VERSION = $(REVISION)

# 安装依赖
.PHONY: install-dep
install-dep:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go mod tidy

# 整理依赖
.PHONY: tidy-dep
tidy-dep:
	go mod tidy

# 编译 proto 文件
.PHONY: proto
proto:
	protoc --gogofaster_out=./ --proto_path=./ api/document.proto
	protoc -I='D:/hjr learing software/MyProject/search_engine/api' --gogofaster_out=plugins=grpc:./ --proto_path=./ api/index.proto
	protoc -I='D:/hjr learing software/MyProject/search_engine/api' --gogofaster_out=./ --proto_path=./ api/term_query.proto