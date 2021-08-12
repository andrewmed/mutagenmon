#/bin/sh
mkdir -p MutagenMon.app/Contents/MacOS
export TERM=dumb
go build go.andmed.org/mutagenmon/cmd/mutagenmon && mv mutagenmon MutagenMon.app/Contents/MacOS/
