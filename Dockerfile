# syntax=docker/dockerfile:1

FROM docker.io/alpine:3.18 AS base

# Build the web ui from source
FROM --platform=$BUILDPLATFORM docker.io/node:18 AS build-node
WORKDIR /exa
ADD internal/web /exa/internal/web
ADD Makefile /exa/
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash
RUN export NVM_DIR="$HOME/.nvm" ; [ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh"; [ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion"
RUN . "$NVM_DIR/nvm.sh" && nvm install 20
RUN . "$NVM_DIR/nvm.sh" && nvm use 20
RUN . "$NVM_DIR/nvm.sh" && make web

# Build the application from source
FROM docker.io/golang:1.24-bookworm AS build-go

ARG TARGETOS TARGETARCH

WORKDIR /exa

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
COPY --link --from=build-node /exa/internal/web/build /exa/internal/web/build
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make app

# Artifact Target
FROM scratch AS artifact

ARG TARGETOS TARGETARCH TARGETVARIANT

COPY --link --from=build-go /exa/build/exatorrent /exatorrent-${TARGETOS}-${TARGETARCH}${TARGETVARIANT}

# Failover if contexts=artifacts=<path> is not set
FROM scratch AS artifacts
# Releaser flat all artifacts
FROM base AS releaser
WORKDIR /out
RUN --mount=from=artifacts,source=.,target=/artifacts <<EOT
    set -e
    cp /artifacts/**/* /out/ 2>/dev/null || cp /artifacts/* /out/
EOT
FROM scratch AS release
COPY --link --from=releaser /out /

# Final stage
# Deploy the application binary into a lean image
FROM base

LABEL maintainer="zakaria-didah"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.name="zakaria-didah/exatorrent"
LABEL org.label-schema.description="self-hostable torrent client"
LABEL org.label-schema.url="https://github.com/zakaria-didah/exatorrent"
LABEL org.label-schema.vcs-url="https://github.com/zakaria-didah/exatorrent"

COPY --link --from=build-go --chown=1000:1000 /exa/build/exatorrent /exatorrent

RUN apk add -U --upgrade --no-cache \
    ca-certificates

USER 1000:1000

WORKDIR /exa

RUN mkdir -p exadir

EXPOSE 5000 42069

VOLUME /exa/exadir

ENTRYPOINT ["/exatorrent"]

CMD ["-dir", "exadir"]
