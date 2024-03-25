@echo off

@REM gocov / gocov-html

go get github.com/axw/gocov/gocov
go get github.com/matm/gocov-html

go build -mod=readonly -o ./bin/ github.com/axw/gocov/gocov
go build -mod=readonly -o ./bin/ github.com/matm/gocov-html

go mod tidy

%GOBIN%\gocov test ./... | %GOBIN%\gocov-html > coverage.html
del %GOBIN%\gocov.exe /Q /S
del %GOBIN%\gocov-html.exe /Q /S


@REM @REM codecov

@REM go test ./... -coverprofile=coverage.out
@REM curl.exe --progress-bar -Lo codecov.exe https://uploader.codecov.io/latest/windows/codecov.exe 

@REM @REM  codecov.exe -t ${CODECOV_TOKEN}
@REM codecov.exe -t %1
@REM del codecov.exe /Q /S
