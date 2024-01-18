echo ".:| Cross compiling for windows, linux and mac |:."

# check if dist folder exists
if [ -d "dist" ]; then
  rm -rf dist
fi

mkdir dist

echo "* Compiling for Mac OS (amd64)"
GOOS=darwin GOARCH=amd64 go build -ldflags '-w' -o dist/wsreplay-darwin-amd64
echo -n " - "; file dist/wsreplay-darwin-amd64

echo "* Compiling for Mac OS (arm64)"
GOOS=darwin GOARCH=arm64 go build -ldflags '-w' -o dist/wsreplay-darwin-arm64
echo -n " - "; file dist/wsreplay-darwin-arm64

echo "Compiling for Linux (amd64)"
GOOS=linux GOARCH=amd64 go build -ldflags '-w' -o dist/wsreplay-linux-amd64
echo -n " - "; file dist/wsreplay-linux-amd64

echo "Compiling for Linux (arm64)"
GOOS=linux GOARCH=arm64 go build -ldflags '-w' -o dist/wsreplay-linux-arm64
echo -n " - "; file dist/wsreplay-linux-arm64

echo "Compiling for windows"
GOOS=windows GOARCH=amd64 go build -ldflags '-w' -o dist/wsreplay.exe
echo -n " - "; file dist/wsreplay.exe

echo "Compiling for windows ARM"
GOOS=windows GOARCH=arm64 go build -ldflags '-w' -o dist/wsreplay-arm64.exe
echo -n " - "; file dist/wsreplay-arm64.exe
