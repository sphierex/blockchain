
# =============================================================================
function run {
    go run cmd/apps/node/main.go
}

# =============================================================================
# Depends

function tidy {
    go mod tidy
    go mod vendor
}

function deps-upgrade {
    go get -u -v ./...
    go mod tidy
    go mod vendor

}

# =============================================================================

$0 = $args[0]

if (Get-Command -Name $0 -CommandType Function -ErrorAction SilentlyContinue) {
    & $0
} else {
    Write-Error "Invaild command."
}