FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:latest
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y gnupg software-properties-common wget

#install kubectl
RUN apt-get install google-cloud-cli-app-engine-go kubectl google-cloud-cli

#install go
RUN curl -O https://dl.google.com/go/go1.23.1.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
RUN rm go1.23.1.linux-amd64.tar.gz
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0
RUN go install -v github.com/go-delve/delve/cmd/dlv@latest

#install terraform
RUN wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
RUN gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint
RUN echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | tee /etc/apt/sources.list.d/hashicorp.list
RUN apt update
RUN apt-get install -y terraform
RUN touch ~/.bashrc
RUN terraform -install-autocomplete

#install docker
RUN apt install -y docker.io