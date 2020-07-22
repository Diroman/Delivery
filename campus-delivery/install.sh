go get github.com/golang/protobuf/protoc-gen-go

PB_REL="https://github.com/protocolbuffers/protobuf/releases"
download_link=$PB_REL/download/v3.11.4/protoc-3.11.4-linux-x86_64.zip
curl -LO download_link
unzip protoc-3.11.4-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"
rm -rf download_link

protoc -I . message.proto --go_out=plugins=grpc:.