ARG RUST_IMAGE=rust:1.95-slim-bookworm
ARG RUNTIME_IMAGE=debian:bookworm-slim

FROM ${RUST_IMAGE} AS build

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends libprotobuf-dev protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock ./
COPY platform ./platform
COPY services ./services
COPY shared ./shared

RUN cargo build --release -p inventory

FROM ${RUNTIME_IMAGE}

COPY --from=build /src/target/release/inventory /usr/local/bin/inventory

EXPOSE 50052

ENTRYPOINT ["inventory"]
