FROM ubuntu:20.04 as downloader

ARG VERSION=v0.0.3-rc2024081902
ARG PLATFORM=linux_amd64

RUN apt-get update && apt-get install -y wget
RUN wget https://github.com/gruntwork-io/terragrunt-engine-go/releases/download/${VERSION}/terragrunt-iac-engine-server_${PLATFORM} -O /tmp/terragrunt-engine-server
RUN chmod +x /tmp/terragrunt-engine-server

FROM gruntwork/terragrunt:0.2.0

ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu

COPY --from=downloader /tmp/terragrunt-engine-server /terragrunt-engine-server

RUN mkdir /app
WORKDIR /app
ENTRYPOINT ["/terragrunt-engine-server"]

