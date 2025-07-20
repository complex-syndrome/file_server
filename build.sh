#!/bin/sh

while IFS='=' read -r key value; do
    case "$key" in
        ''|\#*) continue ;;
    esac

    if [ "$key" = "FRONTEND_PORT" ]; then
        FRONTEND_PORT="$value"
        export FRONTEND_PORT
        break
    fi
done < .env


mkdir -p backend/bin
pushd backend
go build -o ./bin/main.exe main/main.go
popd
echo Backend build completed.

rm -rf frontend/build
pushd frontend 
bun run build
popd
echo Frontend build completed.

echo Build completed. Run "./start.sh" to run the program.