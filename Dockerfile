# This dockerfile builds the build harness for this project. It has everything it needs to build and test this repo.

FROM rockylinux:8

# Make all shells run in a safer way. Ref: https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
SHELL [ "/bin/bash", "-euxo", "pipefail", "-c" ]

WORKDIR /

# Need root to do rooty things.
USER root

# hadolint ignore=DL3041
RUN dnf install -y --refresh \
  findutils \
  git \
  jq \
  make \
  unzip \
  wget \
  which \
  && dnf clean all \
  && rm -rf /var/cache/yum/

# hadolint ignore=DL3059
RUN useradd -ms /bin/bash buildharness

USER buildharness
WORKDIR /home/buildharness

# Install asdf. Get versions from https://github.com/asdf-vm/asdf/releases
ARG ASDF_VERSION="0.10.1"
ENV ASDF_VERSION=${ASDF_VERSION}
# hadolint ignore=SC2016
RUN git clone --branch "v${ASDF_VERSION}" --depth 1 https://github.com/asdf-vm/asdf.git "${HOME}/.asdf" \
  && echo -e '\nsource $HOME/.asdf/asdf.sh' >> "${HOME}/.bashrc" \
  && echo -e '\nsource $HOME/.asdf/asdf.sh' >> "${HOME}/.profile" \
  && source "${HOME}/.asdf/asdf.sh"
ENV PATH="/home/buildharness/.asdf/shims:/home/buildharness/.asdf/bin:${PATH}"

# Install hadolint. Get versions using 'asdf list all hadolint'
ARG HADOLINT_VERSION="2.10.0"
ENV HADOLINT_VERSION=${HADOLINT_VERSION}
RUN asdf plugin add hadolint \
  && asdf install hadolint "${HADOLINT_VERSION}" \
  && asdf global hadolint "${HADOLINT_VERSION}"

# Install pre-commit. Get versions using 'asdf list all pre-commit'
ARG PRE_COMMIT_VERSION="2.19.0"
ENV PRE_COMMIT_VERSION=${PRE_COMMIT_VERSION}
RUN asdf plugin add pre-commit \
  && asdf install pre-commit "${PRE_COMMIT_VERSION}" \
  && asdf global pre-commit "${PRE_COMMIT_VERSION}"

# Install Terraform. Get versions using 'asdf list all terraform'
ARG TERRAFORM_VERSION="1.2.0"
ENV TERRAFORM_VERSION=${TERRAFORM_VERSION}
RUN asdf plugin add terraform \
  && asdf install terraform "${TERRAFORM_VERSION}" \
  && asdf global terraform "${TERRAFORM_VERSION}"

# Install tflint. Get versions using 'asdf list all tflint'
ARG TFLINT_VERSION="0.28.1"
ENV TFLINT_VERSION=${TFLINT_VERSION}
RUN asdf plugin add tflint \
  && asdf install tflint "${TFLINT_VERSION}" \
  && asdf global tflint "${TFLINT_VERSION}"

# Install tfsec. Get versions using 'asdf list all tfsec'
ARG TFSEC_VERSION="0.39.37"
ENV TFSEC_VERSION=${TFSEC_VERSION}
RUN asdf plugin add tfsec \
  && asdf install tfsec "${TFSEC_VERSION}" \
  && asdf global tfsec "${TFSEC_VERSION}"

# Support tools installed as buildharness when running as root user
USER root
ENV ASDF_DATA_DIR="/home/buildharness/.asdf"
RUN cp /home/buildharness/.tool-versions /root/.tool-versions
ENV HELM_PLUGINS="/home/buildharness/.local/share/helm/plugins"
ENV HELM_REGISTRY_CONFIG="/home/buildharness/.config/helm/registry.json"
ENV HELM_REPOSITORY_CACHE="/home/buildharness/.cache/helm/repository"
ENV HELM_REPOSITORY_CONFIG="/home/buildharness/.config/helm/repositories.yaml"
USER buildharness

CMD ["/bin/bash"]
