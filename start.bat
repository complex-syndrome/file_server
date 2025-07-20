@echo OFF

setlocal EnableDelayedExpansion

for /f "usebackq delims=" %%a in (".env") do (
    if not "%%a"=="" if not "%%a:~0,1"=="#" (
        for /f "tokens=1,2 delims==" %%b in ("%%a") do (
            if "%%b"=="FRONTEND_PORT" set "FRONTEND_PORT=%%c"
        )
    )
)

REM To taskill PID using the PORT:
REM "netstat -ano | findstr :${PORT}"
REM "taskkill /PID ${PID} /F"

REM "cmd /k" for 2 terminals, "/B cmd /c" for background 
start "Backend" cmd /k "cd backend && bin\main.exe"
start "Frontend" cmd /k "cd frontend && bun run preview --host --open --port !FRONTEND_PORT!"

REM To access settings in web ui, please use 'bun run dev' instead of 'bun run preview' (Will TODO user auth to replace this)
REM Direct api calls are not affected

endlocal