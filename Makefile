TARGETBIN=updateBsp
.PHONY:	all ${TARGETBIN}.exe ${TARGETBIN} protoc

BUILD_ROOT=$(PWD)
all: ${TARGETBIN}.exe  ${TARGETBIN}

${TARGETBIN}:
	@gofmt -l -w ${BUILD_ROOT}/
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.cn && \
	go build -ldflags "-w -s" -o $@ ithings.go
	@chmod 777 $@
	

${TARGETBIN}-ARM64:
	@gofmt -l -w ${BUILD_ROOT}/
	@echo "编译 ARM64 动态版本..."
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.cn && \
	GOARCH=arm64 GOOS="linux" CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -buildmode=pie -ldflags "-s -w -linkmode external -extldflags '-z relro -z now'" -o $@ ithings.go
	@chmod 777 $@
	@echo "编译完成: $@"
	@file $@

install:
	@mkdir -p out
	@chmod 777 ${TARGETBIN}.exe  ${TARGETBIN}
	@cp -a conf ${TARGETBIN}.exe  ${TARGETBIN}  out/
	sync;sync
	@echo "[Done]"

.PHONY: clean  install
clean:
	@rm -rf ${TARGETBIN}.exe  ${TARGETBIN} *.log *.db *.tar.gz
	@echo "[clean Done]"
