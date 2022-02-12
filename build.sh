#/bin/sh

mkdir -p MutagenMon.app/Contents/MacOS
go build -buildmode=pie  go.andmed.org/mutagenmon/cmd/mutagenmon && mv mutagenmon MutagenMon.app/Contents/MacOS/