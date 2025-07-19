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


if [ ! -d backend/bin ]; then
    mkdir -p backend/bin

    pushd backend
    go build -o ./bin/main.exe main/main.go
    popd
fi

if [ ! -d frontend/build ]; then (
    pushd frontend 
    bun run build
    popd
)

(
    cd backend || exit
    go run main/main.go
) &

(
    cd frontend || exit
    bun run preview --host --open --port "$FRONTEND_PORT"
) &