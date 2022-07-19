FROM golang:1.16.11-stretch AS builder
RUN apt-get update && \
    apt-get install -y wget jq hwloc ocl-icd-opencl-dev git libhwloc-dev pkg-config make && \
    apt-get install -y cargo
WORKDIR /app/

RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
RUN cargo --help
RUN git clone https://github.com/application-research/barge . && \
     make all
RUN cp ./barge /usr/local/bin

FROM golang:1.16.11-stretch
RUN apt-get update && \
    apt-get install -y hwloc libhwloc-dev ocl-icd-opencl-dev

COPY --from=builder /app/barge /usr/local/bin