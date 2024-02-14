GOCMD=go1.22.0
GOFMT=gofmt

build: linux-amd64 linux-arm64 windows-arm64 windows-amd64

linux-amd64:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w" -o bin/airgradient-exporter-linux-amd64 github.com/zerklabs/airgradient-exporter/cmd/airgradient-exporter

linux-arm64:
	GOARCH=arm64 GOOS=linux CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w" -o bin/airgradient-exporter-linux-arm64 github.com/zerklabs/airgradient-exporter/cmd/airgradient-exporter

windows-arm64:
	GOARCH=arm64 GOOS=windows CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w" -o bin/airgradient-exporter-windows-arm64.exe github.com/zerklabs/airgradient-exporter/cmd/airgradient-exporter

windows-amd64:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 $(GOCMD) build -ldflags="-s -w" -o bin/airgradient-exporter-windows-amd64.exe github.com/zerklabs/airgradient-exporter/cmd/airgradient-exporter
