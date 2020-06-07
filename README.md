# Phoenix

In Ancient Greek folklore, a phoenix is a long-lived bird that cyclically regenerates or is otherwise born again. Associated with the sun, a phoenix obtains new life by arising from the ashes of its predecessor. Almost alike this web server can regenerate/update itself into a newer version.  

This was done as a learning exercise to get more fluent in golang, therefore it's not recommended to be used in a productive environment :)

## Install
    go get github.com/sprengr/Phoenix

## Build
    go build -ldflags "-X main.Version=ver1" 

## Build & run with an available update linux
    go build -ldflags "-X main.Version=ver2" -o releases/Phoenix && go build -ldflags "-X main.Version=ver1" && ./Phoenix

## Build & run with an available update windows
    go.exe build -ldflags "-X main.Version=ver2" -o releases/Phoenix.exe && go.exe build -ldflags "-X main.Version=ver1" && Phoenix.exe

## Security considerations
* Only releases using the same name as the running executable and responding to `--version` will be applied. This reduces the risk to release a wrong executable.
* Only regular files will be considered updates.
* `Install` functionality uses `Check`, e.g. you can only install something if `Check` was run. 
* Caution: Everything that's considered an allowed update will be run! To prevent someone to for example upgrade to a "read the manual really fast" `rm -rf /` (😁) version it's recommended to only run the service as an unprivilged user.