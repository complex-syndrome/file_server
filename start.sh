(
    cd backend || exit
    go run main/main.go
) &

(
    cd frontend || exit
    bun run preview --host --open --port "$FRONTEND_PORT"
) &