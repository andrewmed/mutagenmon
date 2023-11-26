#/bin/sh

mkdir -p MutagenMon.app/Contents/MacOS
GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 go build -buildmode=pie  go.andmed.org/mutagenmon/cmd/mutagenmon && mv mutagenmon MutagenMon.app/Contents/MacOS/