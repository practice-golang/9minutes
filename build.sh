cd frontend/admin

yarn
yarn build

rm -rf ../../static/html/admin/_app/*
mv -f ./build/* ../../static/html/admin

rm -rf ./build

cd ../..

go build -ldflags "-w -s" -trimpath -o bin/ ./cmd
