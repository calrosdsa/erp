
PROTO_OUT_DIR = ${PWD}/gen
PROTO_DIT = ${PWD}/idl/proto/api_v1


include Makefile.Protobuf.mk
.PHONY:
protoc:
	protoc -Iidl/proto/api_v1 --go_out=${PROTO_OUT_DIR} --go_opt=paths=source_relative --go-grpc_out=${PROTO_OUT_DIR} \
	--go-grpc_opt=paths=source_relative idl/proto/api_v1/*.proto
	


mockery:
	mockery



tag=latest
docker:
	docker build -t jmiranda0521/erp:$(tag) .
	docker push jmiranda0521/erp:$(tag)	


models:
	@cd cmd/gen-models && go run main.go 

annotations:
	@cd cmd/gen-annotations && go run main.go

.PHONY: accounting
accounting:
	@cd accounting/cmd && go run main.go

.PHONY: app
app:
	@cd cmd/app && go run main.go


# STARTING CONSUL AGENT
.PHONY: consul
consul:
	consul agent -dev -ui -node member


# generate:
# 	@echo running code generation
# 	@go generate ./...
# 	@echo done

doc:
	docker run --rm -v ${PWD}/documents:/documents asciidoctor/docker-asciidoctor sample.adoc

test:
	go test ./...

test-short:
	go test -short ./...


start:
	@cd cmd/all && go run main.go


PG_URL = postgresql://postgres:12ab34cd56ef@10.0.0.151:5432/erp_dev

backup-test-db:
	pg_dump ${PG_URL} --format plain --data-only --verbose --file "db/data.sql" --table entities \
	--table actions --table currencies --table party_types --table states --table unit_of_measures \
	--table unit_of_measure_translations 
	pg_dump ${PG_URL} --format plain --schema-only --verbose --file "db/schema.sql"
	cd db && copy schema.sql + data.sql + custom-data.sql init.sql


