FROM gruntwork/terragrunt:0.2.0

ENV LISTEN_ADDRESS=0.0.0.0:50051

RUN mise use --global -y opentofu@1.7.0
RUN ln -s /root/.local/share/mise/installs/opentofu/1.7.0/bin/tofu /bin/tofu
RUN mkdir /app
ADD terragrunt-engine-server /app/terragrunt-engine-server

ENTRYPOINT ["/app/terragrunt-engine-server"]