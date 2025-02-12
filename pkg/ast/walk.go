package ast

import (
	"context"
	"os"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
	"google.golang.org/grpc"
)

type Visitor interface {
	Visit(*unmangov1alpha1.Node) bool
}

func Walk(ctx context.Context, v Visitor, node *unmangov1alpha1.Node) error {
	sock, err := os.CreateTemp("", "")

	server := grpc.NewServer()
	unmangov1alpha1.RegisterVisitorServiceServer(server, NewServer(v))

	client := unmangov1alpha1.NewAstServiceClient(nil)
	_, err := client.Walk(ctx, &unmangov1alpha1.WalkRequest{
		Node: node,
	})
}
