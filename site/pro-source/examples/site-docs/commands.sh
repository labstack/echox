#!/usr/bin/env bash
# SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
# SPDX-FileCopyrightText: © 2026 LabStack LLC
# Source this file to select an individual documentation command.

install_pro() {
  # docs:start protect-netrc
  chmod 600 "$HOME/.netrc"
  # docs:end protect-netrc
  # docs:start install
  # docs:start private-pattern
  go env -w GOPRIVATE=echo.labstack.com
  # docs:end private-pattern
  # docs:start get-version
  go get echo.labstack.com/pro@v0.2.0
  # docs:end get-version
  go mod verify
  # docs:end install
}

probe_registry() {
  # docs:start proxy-probe
  curl --silent --show-error --netrc-file "$HOME/.netrc" --output /dev/null --write-out '%{http_code}\n' \
    https://go.labstack.io/echo.labstack.com/pro/@v/v0.2.0.info
  # docs:end proxy-probe
}

mirror_pro() {
  # docs:start mirror
  go env -w GOPROXY=https://mirror.corp.example/go GONOSUMDB=echo.labstack.com GONOPROXY=none
  # docs:end mirror
}

inspect_go_settings() {
  # docs:start settings
  go env GOPROXY GOPRIVATE GONOPROXY GONOSUMDB GOSUMDB GOFLAGS GOAUTH
  go list -m echo.labstack.com/pro
  # docs:end settings
}

module_directory() {
  # docs:start module-directory
  cd "$(go list -m -f '{{.Dir}}' echo.labstack.com/pro)"
  # docs:end module-directory
}

run_getting_started() {
  # docs:start getting-started
  go run ./examples/getting-started -browser -key-file /tmp/echo-pro-getting-started.key
  # docs:end getting-started
}

exercise_getting_started() {
  # docs:start unauthenticated
  curl --include http://127.0.0.1:8080/v1/tenants/acme/invoices
  # docs:end unauthenticated
  # docs:start authenticated
  # The key stays in curl's header file input; shell tracing never expands it.
  {
    printf 'X-API-Key: '
    cat /tmp/echo-pro-getting-started.key
  } |
    curl --include --header @- http://127.0.0.1:8080/v1/tenants/acme/invoices
  # docs:end authenticated
  # docs:start cross-tenant
  {
    printf 'X-API-Key: '
    cat /tmp/echo-pro-getting-started.key
  } |
    curl --include --header @- http://127.0.0.1:8080/v1/tenants/globex/invoices
  # docs:end cross-tenant
  # docs:start create-invoice
  {
    printf 'X-API-Key: '
    cat /tmp/echo-pro-getting-started.key
  } |
    curl --include --header @- --request POST http://127.0.0.1:8080/v1/tenants/acme/invoices
  # docs:end create-invoice
}

run_recipes() {
  # docs:start admin
  go run ./examples/admin
  # docs:end admin
  # docs:start admin-demo
  go run ./examples/admin-demo -addr 127.0.0.1:8080
  # docs:end admin-demo
  # docs:start quota-otel
  go run ./examples/quota-otel
  # docs:end quota-otel
  # docs:start stripe
  go run ./examples/stripe
  # docs:end stripe
  # docs:start stripe-online
  go run ./examples/stripe -online
  # docs:end stripe-online
}

test_packages() {
  # docs:start test-auth
  go test -run Example ./auth/...
  # docs:end test-auth
  # docs:start test-tenant
  go test -run Example ./tenant/...
  # docs:end test-tenant
  # docs:start test-quota
  go test -run Example ./quota/...
  # docs:end test-quota
  # docs:start test-webhook
  go test -run Example ./webhook/...
  # docs:end test-webhook
  # docs:start test-audit
  go test -run Example ./audit/...
  # docs:end test-audit
  # docs:start test-execution
  go test -run Example ./execution/...
  # docs:end test-execution
  # docs:start test-admin
  go test -run Example ./admin/...
  # docs:end test-admin
  # docs:start test-recipes
  go test ./examples/getting-started ./examples/quota-otel ./examples/admin-demo
  # docs:end test-recipes
}

build_application() {
  # docs:start docker-build
  docker build --secret id=netrc,src="$HOME/.netrc" -t invoice-api:v0.2.0 .
  # docs:end docker-build
  # docs:start vendor
  go mod vendor
  go test -mod=vendor ./...
  go build -mod=vendor ./...
  # docs:end vendor
}
