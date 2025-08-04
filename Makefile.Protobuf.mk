DOCKER_PROTOBUF_VERSION=0.5.0
DOCKER_PROTOBUF=jaegertracing/protobuf:$(DOCKER_PROTOBUF_VERSION)
PROTOC := docker run --rm -u ${shell id -u} -v$(shell pwd):$(shell pwd) -w$(shell pwd) ${DOCKER_PROTOBUF} --proto_path=$(shell pwd)

PROTO_OUT_DIR = ${PWD}/gen



PROTO_INCLUDES := \
	-Iidl/proto/api_v1

# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
define print_caption
  @echo "🏗️ "
  @echo "🏗️ " $1
  @echo "🏗️ "

endef



# Macro to compile Protobuf $(2) into directory $(1). $(3) can provide additional flags.
# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
# Arguments:
#  $(1) - output directory
#  $(2) - path to the .proto file
#  $(3) - additional flags to pass to protoc, e.g. extra -Ixxx
define proto_compile
  $(call print_caption, "Processing $(2) --> $(1)")

  protoc \
    $(PROTO_INCLUDES) \
	--go_out=$(strip $(1)) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(strip $(1)) \
	--go-grpc_opt=paths=source_relative $(2)

endef

# proto:
# protoc -Iidl/proto/api_v1 --go_out=${PROTO_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${PROTO_OUT_DIR} \
# --go-grpc_opt=paths=source_relative idl/proto/api_v1/*.proto


.PHONY: proto
proto: proto-accounting

.PHONY: proto-accounting
proto-accounting:
	$(call proto_compile,${PWD}/gen/proto,idl/proto/api_v1/*.proto)
