@echo OFF

setlocal EnableDelayedExpansion

for /f "usebackq delims=" %%a in (".env") do (
    if not "%%a"=="" if not "%%a:~0,1"=="#" (
        for /f "tokens=1,2 delims==" %%b in ("%%a") do (
            if "%%b"=="FRONTEND_PORT" set "FRONTEND_PORT=%%c"
        )
    )
)

REM "cmd /k" for 2 terminals, "/B cmd /c" for background 

REM To taskill PID:
REM "netstat -ano | findstr :8080"
REM "taskkill /PID ${PID} /F"

if not exist backend/bin (
    mkdir backend/bin > /dev/null 2>&1

    pushd backend 
    go build -o ./bin/main.exe main/main.go
    popd
)

if not exist frontend/build (
    pushd frontend 
    bun run build
    popd

    echo "Init completed. Run this script one more time after program ends."
)

start "Backend" cmd /k "cd backend && bin\main.exe"
start "Frontend" cmd /k "cd frontend && bun run preview --host --open --port !FRONTEND_PORT!"

endlocal