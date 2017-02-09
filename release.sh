#!/bin/bash
#https://gist.github.com/bclinkinbeard/1331790
versionLabel=$1
productName="bitfan"

rm releases/*.tgz
rm releases/*.zip

go generate .

arch=amd64
os=darwin
product=${productName}_v${versionLabel}_${os}_${arch}
env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/$product .
tar czfv releases/$product.tgz releases/$product
rm releases/$product

arch=amd64
os=linux
product=${productName}_v${versionLabel}_${os}_${arch}
env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/$product .
tar czfv releases/$product.tgz releases/$product
rm releases/$product

#arch=386
#os=linux
#product=${productName}_v${versionLabel}_${os}_${arch}
#env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/$product .
#tar czfv releases/$product.tgz releases/$product
#rm releases/$product

arch=arm
os=linux
product=${productName}_v${versionLabel}_${os}_${arch}
env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/$product .
tar czfv releases/$product.tgz releases/$product
rm releases/$product

arch=amd64
os=windows
product=${productName}_v${versionLabel}_${os}_${arch}
env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/${product}.exe .
zip -r releases/${product}.zip releases/${product}.exe
rm releases/$product.exe

#arch=386
#os=windows
#product=${productName}_v${versionLabel}_${os}_${arch}
#env GOOS=$os GOARCH=$arch go build -ldflags="-s -w -X main.version=${1} -X main.buildstamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'` " -o releases/${product}.exe .
#zip -r releases/${product}.zip releases/${product}.exe
#rm releases/$product.exe

