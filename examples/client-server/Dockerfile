FROM golang:1.21.13-bullseye as builder

ARG ENGINE_VERSION=v0.0.3-rc2024081902

# clone base engine repo
WORKDIR src
RUN git clone https://github.com/gruntwork-io/terragrunt-engine-go.git -b ${ENGINE_VERSION} .
# add examples from working directory
RUN rm -rf ./examples/*
COPY . ./examples/client-server

# build engine
WORKDIR examples/client-server
RUN make

FROM gruntwork/terragrunt:0.2.0

ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu

COPY --from=builder /go/src/examples/client-server/terragrunt-engine-server /opt/terragrunt-engine-server

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/opt/terragrunt-engine-server"]

