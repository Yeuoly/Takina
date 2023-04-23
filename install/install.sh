go build --tags netgo cmd/client_cli/main.go
rm -f src/api/client_cli
mv main src/api/client_cli
go build --tags netgo cmd/client_daemon/main.go
rm -f src/api/client_daemon
mv main src/api/client_daemon