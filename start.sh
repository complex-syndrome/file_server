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

(
    cd backend || exit
    go run main/main.go
) &

(
    cd frontend || exit
    bun run preview --host --open --port "$FRONTEND_PORT"
) &