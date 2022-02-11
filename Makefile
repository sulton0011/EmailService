CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

# proto file function
proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}
	ls genproto/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"

# migration up function
migrations-up:
	migrate -source file:./migrations -database 'postgres://mac:sulton0011@localhost:5432/email?sslmode=disable' up

# migration down function
migrations-down:
	migrate -source file:./migrations -database 'postgres://mac:sulton0011@localhost:5432/email?sslmode=disable' down

