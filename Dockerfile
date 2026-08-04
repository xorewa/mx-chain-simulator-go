FROM golang:1.26.2 AS builder


WORKDIR /multiversx
COPY . .

RUN go mod tidy

WORKDIR /multiversx/cmd/chainsimulator

RUN go build -o chainsimulator

RUN mkdir -p /lib_amd64 /lib_arm64

RUN vm_v14_mod="$(go list -m -f '{{if .Replace}}{{.Replace.Path}}@{{.Replace.Version}}{{else}}{{.Path}}@{{.Version}}{{end}}' github.com/multiversx/mx-chain-vm-v1_4-go)" && \
    cp "$(go env GOPATH)/pkg/mod/${vm_v14_mod}/wasmer/libwasmer_linux_amd64.so" /lib_amd64/ && \
    cp "$(go env GOPATH)/pkg/mod/${vm_v14_mod}/wasmer/libwasmer_linux_arm64_shim.so" /lib_arm64/

RUN vm_exec_mod="$(go list -m -f '{{if .Replace}}{{.Replace.Path}}@{{.Replace.Version}}{{else}}{{.Path}}@{{.Version}}{{end}}' github.com/multiversx/mx-chain-vm-go)" && \
    cp "$(go env GOPATH)/pkg/mod/${vm_exec_mod}/wasmer2/libvmexeccapi.so" /lib_amd64/ && \
    cp "$(go env GOPATH)/pkg/mod/${vm_exec_mod}/wasmer2/libvmexeccapi_arm.so" /lib_arm64/


FROM ubuntu:22.04
ARG TARGETARCH
RUN apt-get update && apt-get install -y git curl

COPY --from=builder /multiversx/cmd/chainsimulator /multiversx

EXPOSE 8085

WORKDIR /multiversx

# Copy architecture-specific files
COPY --from=builder "/lib_${TARGETARCH}/*" "/lib/"

CMD ["/bin/bash"]

ENTRYPOINT ["./chainsimulator"]

