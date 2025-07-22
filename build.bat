@echo OFF

del /f /q "backend/bin/main.exe" > NUL 2>&1
mkdir backend/bin > NUL 2>&1
pushd backend 
go build -o ./bin/main.exe main/main.go
popd
echo Backend build completed.

rmdir /s /q "frontend/build" > NUL 2>&1
pushd frontend
bun run build
popd
echo Frontend build completed.

echo Build completed. Run ".\start.bat" to run the program.