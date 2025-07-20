docker_compose() {
    docker compose \
        -f Docker/docker-compose.yaml \
        --project-directory . \
        $@
}

svr() {
    go tool templ generate \
        --watch \
        --proxy='http://localhost:3000' \
        --cmd='go run cmd/app/main.go'
}

web() {
    npx parcel watch --no-cache
}

migrate_postgres() {
    migrate \
        -source "file://internal/repos/postgres/Migrations" \
        -database $APP_POSTGRES_URL \
        $@
}


$@
