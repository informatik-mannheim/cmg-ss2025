FROM mcr.microsoft.com/azure-cli:2.70.0

RUN tdnf install -y wget tar awk git

#install kubectl
ARG BUILDARCH
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${BUILDARCH}/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

#install go
ARG BUILDARCH
RUN curl -O https://dl.google.com/go/go1.24.1.linux-${BUILDARCH}.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.1.linux-${BUILDARCH}.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN rm go1.24.1.linux-${BUILDARCH}.tar.gz
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0
RUN go install -v github.com/go-delve/delve/cmd/dlv@v1.24.0


#install terraform
# RUN wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
# RUN gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint
# RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list
# RUN apt update
# RUN apt-get install -y terraform
# RUN touch ~/.bashrc
# RUN terraform -install-autocomplete

#install docker
RUN tdnf install -y moby-cli
