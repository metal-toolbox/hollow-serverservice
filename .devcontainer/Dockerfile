# Used to install CRDB into the devcontainer
FROM cockroachdb/cockroach:latest-v22.2 as CRDB

FROM mcr.microsoft.com/vscode/devcontainers/go:0-1.19-buster

ENV NATS_CLI_VERSION=0.0.35

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        bash-completion \
        uuid-runtime \
        postgresql-client \
        python3-dev \
        python3-pip \
        python3-setuptools \
    && pip3 install semgrep


# Set up crdb
RUN mkdir /usr/local/lib/cockroach
COPY --from=CRDB /cockroach/cockroach /usr/local/bin
COPY --from=CRDB /usr/local/lib/cockroach/libgeos.so /usr/local/lib/cockroach/
COPY --from=CRDB /usr/local/lib/cockroach/libgeos_c.so /usr/local/lib/cockroach/

# Install NATS Tooling
RUN curl -o /tmp/install.sh  https://raw.githubusercontent.com/nats-io/nsc/main/install.sh \
     && chmod +x /tmp/install.sh \
     && su vscode -c "/tmp/install.sh -s /usr/local/bin" \
     && rm -f /tmp/install.sh

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
RUN go install github.com/volatiletech/sqlboiler/v4@latest \
    && go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest \
    && go install github.com/glerchundi/sqlboiler-crdb/v4@latest \
    && go install github.com/nats-io/natscli/nats@v${NATS_CLI_VERSION} \
    && go install github.com/nats-io/nkeys/nk@latest

RUN mkdir /home/vscode/.ssh && \
    chown vscode:vscode /home/vscode/.ssh && \
    chmod 0700 /home/vscode/.ssh

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

WORKDIR /workspace
