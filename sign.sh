#/bin/sh

security find-identity -v -p codesigning
codesign -s "Apple Distribution: Andrey Medvedev (BPN9958X73)" --verbose=4 MutagenMon.app
