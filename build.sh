#Before we can use the script, we have to make it executable with the chmod command:
#chmod +x ./go-executable-build.sh
#then we can use it  ./go-executable-build.sh yourpackage
#!/usr/bin/env bash

NOW=$(date +"%y%m%d")
APP="frantic-cat"
figlet -f standard "Building $APP" 
echo Building Crossplatform 
echo Clearing old builds
rm -rf ./exec/$NOW
echo Building Windows x86_64 
go-winres simply --icon ./res/icons/app.png
env GOOS=windows GOARCH=amd64 go build  -o "./exec/"$NOW"/windows/"$APP".exe" github.com/mt1976/$APP
echo Building MacOs x86_64 
go-winres simply --icon ./res/icons/app.png
env GOOS=darwin GOARCH=amd64 go build -o "./exec/"$NOW"/darwin/intel/"$APP github.com/mt1976/$APP 
echo Building MacOs arm64 
go-winres simply --icon ./res/icons/app.png
env GOOS=darwin GOARCH=arm64 go build -o "./exec/"$NOW"/darwin/apple/"$APP github.com/mt1976/$APP 
echo Done 