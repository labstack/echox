# SPDX-License-Identifier: LicenseRef-LabStack-Proprietary
# SPDX-FileCopyrightText: © 2026 LabStack LLC
FROM golang:1.27.3 AS build
WORKDIR /app
ENV GOPRIVATE=echo.labstack.com GONOPROXY=echo.labstack.com GONOSUMDB=echo.labstack.com CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=secret,id=netrc,target=/root/.netrc,required=true go mod download
COPY . .
RUN --network=none go build -mod=readonly -o /out/server ./cmd/server

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /out/server /server
USER 65532:65532
ENTRYPOINT ["/server"]
