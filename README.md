# Phoenix

## Build
    go build -ldflags "-X main.Version=ver1" 

## Build & run with an available update linux
    go build -ldflags "-X main.Version=ver2" -o releases/Phoenix && go build -ldflags "-X main.Version=ver1" && ./Phoenix

## Build & run with an available update windows
    go.exe build -ldflags "-X main.Version=ver2" -o releases/Phoenix && go.exe build -ldflags "-X main.Version=ver1" && Phoenix.exe
