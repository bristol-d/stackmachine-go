@rem build the wasm module
@setlocal
@cd wasm
@echo Building wasm ...
@del go-stackmachine.wasm
@set GOOS=js
@set GOARCH=wasm
go build -o go-stackmachine.wasm
@if exist go-stackmachine.wasm (
    copy go-stackmachine.wasm ..\web
)
