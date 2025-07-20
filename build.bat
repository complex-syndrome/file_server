@echo OFF

mkdir backend/bin > /dev/null 2>&1
pushd backend 
go build -o ./bin/main.exe main/main.go
popd

pushd frontend 
bun run build
popd

echo Build completed. Run "./start.bat" to run the program.
