FROM golang:1.25.5-bookworm AS builder

# Copy the entire repository into the build context
WORKDIR /build
COPY . .

# build engine
WORKDIR examples/client-server
RUN make

FROM gruntwork/terragrunt:2.1.7

ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu

COPY --from=builder /build/examples/client-server/terragrunt-engine-server /opt/terragrunt-engine-server

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/opt/terragrunt-engine-server"]

