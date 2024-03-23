cd frontend/admin

yarn
yarn build
remove-item -r -force -ea 0 ../../static/html/admin/_app
move-item -force ./build/* ../../static/html/admin

remove-item -r -force -ea 0 ./build

cd ../..

go build -ldflags "-w -s" -trimpath -o bin/ ./cmd
