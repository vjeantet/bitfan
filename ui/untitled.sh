CC="x86_64-w64-mingw32-gcc" CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.version=TMPSQL -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/ETS.exe .
