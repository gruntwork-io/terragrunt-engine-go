FROM golang:1.25.5-bookworm AS builder

# FIXME: Set to stable version of engine
ARG ENGINE_VERSION=feat/updating-schema-using-oneof

# clone base engine repo
WORKDIR src
RUN git clone https://github.com/gruntwork-io/terragrunt-engine-go.git -b ${ENGINE_VERSION} .
# add examples from working directory
RUN rm -rf ./examples/*
COPY . ./examples/client-server

# build engine
WORKDIR examples/client-server
RUN make

FROM gruntwork/terragrunt:2.1.7

ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu

COPY --from=builder /go/src/examples/client-server/terragrunt-engine-server /opt/terragrunt-engine-server

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/opt/terragrunt-engine-server"]

