#/bin/sh
set -e

security find-identity -v -p codesigning
codesign --force --options=runtime --timestamp -s "Developer ID Application: Andrey Medvedev (BPN9958X73)" --verbose=2 MutagenMon.app
ditto -c -k --keepParent MutagenMon.app MutagenMon.app.zip 
xcrun notarytool submit MutagenMon.app.zip   --keychain-profile "MutagenMon"    --wait
spctl --assess --verbose MutagenMon.app