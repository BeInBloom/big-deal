ARG RUST_IMAGE=rust:1.87-slim-bookworm
ARG RUNTIME_IMAGE=debian:bookworm-slim

FROM ${RUST_IMAGE} AS build

WORKDIR /src

RUN apt-get update \
    && apt-get install -y --no-install-recommends protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

COPY Cargo.toml Cargo.lock ./
COPY platform ./platform
COPY services ./services
COPY shared ./shared

RUN cargo build --release -p payment

FROM ${RUNTIME_IMAGE}

COPY --from=build /src/target/release/payment /usr/local/bin/payment

EXPOSE 50051

ENTRYPOINT ["payment"]
