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
start "Backend" cmd /k "cd backend && bin\main.exe"
start "Frontend" cmd /k "cd frontend && bun run preview --host --open --port !FRONTEND_PORT!"

REM To taskill PID that are using PORT:
REM "netstat -ano | findstr :${PORT}"
REM "taskkill /PID ${PID} /F"

REM Direct api calls are not affected

endlocal