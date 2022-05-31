# This dockerfile builds the build harness for this project. It has everything it needs to build and test this repo.

FROM rockylinux:8

# Make all shells run in a safer way. Ref: https://vaneyckt.io/posts/safer_bash_scripts_with_set_euxo_pipefail/
SHELL [ "/bin/bash", "-euxo", "pipefail", "-c" ]

# hadolint ignore=DL3041
RUN dnf install -y --refresh \
  bind-utils \
  bzip2 \
  bzip2-devel \
  findutils \
  gcc \
  gcc-c++ \
  gettext \
  git \
  jq \
  libffi-devel \
  libxslt-devel \
  make \
  ncurses-devel \
  openssl-devel \
  perl-Digest-SHA \
  readline-devel \
  sqlite-devel \
  unzip \
  wget \
  which \
  xz \
  && dnf clean all \
  && rm -rf /var/cache/yum/

# Install asdf. Get versions from https://github.com/asdf-vm/asdf/releases
ARG ASDF_VERSION="0.10.1"
ENV ASDF_VERSION=${ASDF_VERSION}
# hadolint ignore=SC2016
RUN git clone --branch "v${ASDF_VERSION}" --depth 1 https://github.com/asdf-vm/asdf.git "${HOME}/.asdf" \
  && echo -e '\nsource $HOME/.asdf/asdf.sh' >> "${HOME}/.bashrc" \
  && echo -e '\nsource $HOME/.asdf/asdf.sh' >> "${HOME}/.profile" \
  && source "${HOME}/.asdf/asdf.sh"
ENV PATH="/root/.asdf/shims:/root/.asdf/bin:${PATH}"
# ENV PATH="/home/buildharness/.asdf/shims:/home/buildharness/.asdf/bin:${PATH}"

# Install golang. Get versions using 'asdf list all golang'
ARG GOLANG_VERSION="1.18.2"
ENV GOLANG_VERSION=${GOLANG_VERSION}
RUN asdf plugin add golang \
  && asdf install golang "${GOLANG_VERSION}"

# Install python. Get versions using 'asdf list all python'
ARG PYTHON_VERSION="3.10.4"
ENV PYTHON_VERSION=${PYTHON_VERSION}
RUN asdf plugin add python \
  && asdf install python "${PYTHON_VERSION}"

# Install hadolint. Get versions using 'asdf list all hadolint'
ARG HADOLINT_VERSION="2.10.0"
ENV HADOLINT_VERSION=${HADOLINT_VERSION}
RUN asdf plugin add hadolint \
  && asdf install hadolint "${HADOLINT_VERSION}"

# Install pre-commit. Get versions using 'asdf list all pre-commit'
ARG PRE_COMMIT_VERSION="2.19.0"
ENV PRE_COMMIT_VERSION=${PRE_COMMIT_VERSION}
RUN asdf plugin add pre-commit \
  && asdf install pre-commit "${PRE_COMMIT_VERSION}"

# Install Terraform. Get versions using 'asdf list all terraform'
ARG TERRAFORM_VERSION="1.2.0"
ENV TERRAFORM_VERSION=${TERRAFORM_VERSION}
RUN asdf plugin add terraform \
  && asdf install terraform "${TERRAFORM_VERSION}"

# Install tflint. Get versions using 'asdf list all tflint'
ARG TFLINT_VERSION="0.28.1"
ENV TFLINT_VERSION=${TFLINT_VERSION}
RUN asdf plugin add tflint \
  && asdf install tflint "${TFLINT_VERSION}"

# Install tfsec. Get versions using 'asdf list all tfsec'
ARG TFSEC_VERSION="0.39.37"
ENV TFSEC_VERSION=${TFSEC_VERSION}
RUN asdf plugin add tfsec \
  && asdf install tfsec "${TFSEC_VERSION}"

# Support tools installed as root when running as any other user
ENV ASDF_DATA_DIR="/root/.asdf"

CMD ["/bin/bash"]
