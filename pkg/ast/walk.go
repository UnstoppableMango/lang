package ast

import (
	"context"
	"net"
	"os"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
	"google.golang.org/grpc"
)

type Visitor interface {
	Visit(*unmangov1alpha1.Node) bool
}

func Walk(ctx context.Context, v Visitor, node *unmangov1alpha1.Node) error {
	sock, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", sock.Name())
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	unmangov1alpha1.RegisterVisitorServiceServer(server, NewServer(v))
	go server.Serve(lis)
	defer server.GracefulStop()

	conn, err := grpc.NewClient(sock.Name())
	if err != nil {
		return err
	}

	client := unmangov1alpha1.NewAstServiceClient(conn)
	_, err = client.Walk(ctx, &unmangov1alpha1.WalkRequest{
		Uri:  sock.Name(),
		Node: node,
	})

	return err
}
