.PHONY:

run:
	go run user_service/cmd/main.go

build_user_proto: 
	protoc \
	--proto_path=proto \
	--go_out=user_service/userpb --go_opt=paths=source_relative \
	--go-grpc_out=user_service/userpb --go-grpc_opt=paths=source_relative \
	proto/user/v1/user.proto

build_movie_proto:
	protoc \
	--proto_path=proto \
	--go_out=movie_service/moviepb --go_opt=paths=source_relative \
	--go-grpc_out=movie_service/moviepb --go-grpc_opt=paths=source_relative \
	proto/movie/v1/movie.proto