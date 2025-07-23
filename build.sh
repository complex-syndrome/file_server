#!/bin/sh

mkdir -p backend/bin
pushd backend && go build -o ./bin/main.exe main/main.go && popd
echo Backend build completed.

rm -rf frontend/build
pushd frontend && bun install && bun run build && popd
echo Frontend build completed.

echo Build completed. Run "./start.sh" to run the program.