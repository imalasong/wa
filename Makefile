# 版权 @2019 凹语言 作者。保留所有权利。

hello:
	go run main.go run hello.wa

prime:
	go run main.go run _examples/prime

build-wasm:
	GOARCH=wasm GOOS=js go build -o wa.out.wasm ./main_wasm.go

clean:
