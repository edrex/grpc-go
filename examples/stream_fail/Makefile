foo/foo.pb.go: foo/foo.proto
	protoc --plugin=$(which protoc-gen-go) foo/foo.proto --go_out=plugins=grpc:foo
