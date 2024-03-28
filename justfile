build *args:
    go build -o build/frankenrepo {{args}} .

run *args:
    go run . {{args}}

test *args:
    go test {{args}} ./...

debug *args:
    dlv debug -- {{args}}

test-c_rust:
    cd ./test/integration/c_rust/exe/cproj; make clean
    cd ./test/integration/c_rust/lib/rustproj; cargo clean
    go run . -C ./test/integration/c_rust test
