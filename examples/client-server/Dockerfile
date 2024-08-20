FROM gruntwork/terragrunt:0.2.0

ARG VERSION=v0.0.3-rc2024081902
ARG PLATFORM=linux_amd64
ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN apt-get update && apt-get install -y wget

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu
RUN mkdir /app
RUN wget https://github.com/gruntwork-io/terragrunt-engine-go/releases/download/${VERSION}/terragrunt-iac-engine-server_${PLATFORM} -O /app/terragrunt-engine-server
RUN chmod +x /app/terragrunt-engine-server
ENTRYPOINT ["/app/terragrunt-engine-server"]