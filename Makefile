.PHONY: darwin-amd64 darwin-arm64 freebsd-386 freebsd-amd64 freebsd-arm linux-386 linux-amd64 linux-arm linux-arm64 windows-386 windows-amd64



all:darwin-amd64 darwin-arm64 freebsd-386 freebsd-amd64 freebsd-arm linux-386 linux-amd64 linux-arm linux-arm64 windows-386 windows-amd64

darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/hb_darwin_amd64 -ldflags "-s -w"
darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/hb_darwin_arm64 -ldflags "-s -w"
freebsd-386:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o build/hb_freebsd_386 -ldflags "-s -w"
freebsd-amd64:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o build/hb_freebsd_amd64 -ldflags "-s -w"
freebsd-arm:
	CGO_ENABLED=0 GOOS=freebsd GOARCH=arm go build -o build/hb_freebsd_arm -ldflags "-s -w"
linux-386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o build/hb_linux_386 -ldflags "-s -w"
linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/hb_linux_amd64 -ldflags "-s -w"
linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o build/hb_linux_arm -ldflags "-s -w"
linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/hb_linux_arm64 -ldflags "-s -w"
windows-386:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/hhb_windows_win32.exe -ldflags "-s -w"
windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/hb_windows_win64.exe -ldflags "-s -w"

