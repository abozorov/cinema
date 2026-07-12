.PHONY: services

run_user_service:
	go run cmd/user/main.go

run_movie_service:
	go run cmd/movie/main.go


.PHONY:
build_user_proto: 
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/userpb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/userpb --go-grpc_opt=paths=source_relative \
	user/v1/user.proto

build_movie_proto:
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/moviepb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/moviepb --go-grpc_opt=paths=source_relative \
	movie/v1/movie.proto

build_booking_proto:
	protoc \
	-I=./grpc_api/ \
	--go_out=./grpc_api/generate/bookingpb --go_opt=paths=source_relative \
	--go-grpc_out=./grpc_api/generate/bookingpb --go-grpc_opt=paths=source_relative \
	booking/v1/booking.proto