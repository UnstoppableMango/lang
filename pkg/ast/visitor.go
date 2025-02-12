package ast

import (
	"context"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
)

type server struct {
	unmangov1alpha1.UnimplementedVisitorServiceServer
	visitor Visitor
}

// Visit implements unmangov1alpha1.VisitorServiceServer.
func (s *server) Visit(ctx context.Context, req *unmangov1alpha1.VisitRequest) (*unmangov1alpha1.VisitResponse, error) {
	return &unmangov1alpha1.VisitResponse{
		Continue: s.visitor.Visit(req.Node),
	}, nil
}

func NewServer(visitor Visitor) unmangov1alpha1.VisitorServiceServer {
	return &server{visitor: visitor}
}
