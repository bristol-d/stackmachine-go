# Stack machine example

This is an example 16-bit stack machine for teaching computer architecture.

The machine has a 16-bit integer data type, 4K each of code and data memory,
a value stack (capacity 256), and a return stack (capacity 256).

Instructions are 16 bits long; instructions that take an argument are a total
of 32 bits long.

See [docs.md](docs.md) for more details, or try the included [3n+1.asm](3n+1.asm) or [bcd.asm](bcd.asm) example programs that you can paste into the running simulator.

## Building

With golang installed, run `build.bat` in the main folder to produce the binary.

On linux, go to wasm/ and run `GOOS=js GOARCH=wasm go build -o go-stackmachine.wasm` then copy the file to web/.

Serve the files from web, for example with the `simplehttpserver` go module (browsers will complain about loading wasm from file:// so you can't just open index.html.)
