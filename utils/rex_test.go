package utils

import (
	"testing"
)

func TestRex(t *testing.T) {
	str := `
		service ExampleService {
			rpc SayHello ( .pb.HelloRequest ) returns ( .pb.HelloResponse );
			rpc AnotherFunction ( .pb.AnotherRequest ) returns ( .pb.AnotherResponse );
			rpc ThirdFunction ( .pb.ThirdRequest ) returns ( .pb.ThirdResponse );
			rpc ServerReflectionInfo ( stream .grpc.reflection.v1.ServerReflectionRequest ) returns ( stream .grpc.reflection.v1.ServerReflectionResponse );
			rpc ServerReflectionInfo ( .grpc.reflection.v1.ServerReflectionRequest ) returns ( stream .grpc.reflection.v1.ServerReflectionResponse );
			rpc ServerReflectionInfo ( stream .grpc.reflection.v1.ServerReflectionRequest ) returns ( .grpc.reflection.v1.ServerReflectionResponse );
		}
	`

	funcs, err := GetFunction(str)
	if err != nil {
		t.Fatal(err)
	}

	if len(funcs) != 6 {
		t.Fatalf("expected 3 functions, got %d", len(funcs))
	}

	t.Log(funcs)
}
