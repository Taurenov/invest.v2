FROM rust:1.82-bookworm AS build
WORKDIR /app
COPY fin-math ./fin-math
COPY fin-grpc ./fin-grpc
WORKDIR /app/fin-grpc
RUN cargo build --release --bin fin-engine

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/fin-grpc/target/release/fin-engine /fin-engine
EXPOSE 50052
ENV ENGINE_HTTP_ADDR=0.0.0.0:50052
CMD ["/fin-engine"]
