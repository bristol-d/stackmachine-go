@rem build the wasm module with tinygo
@setlocal
@cd wasm
@set GOROOT=C:\Go
@echo Building wasm ...
@del go-stackmachine.wasm
@rem tinygo build -o go-stackmachine.wasm -target wasm ./wasm-main.go
@set GOOS=js
@set GOARCH=wasm
go build -o go-stackmachine.wasm
@if exist go-stackmachine.wasm (
    copy go-stackmachine.wasm ..\web
)
