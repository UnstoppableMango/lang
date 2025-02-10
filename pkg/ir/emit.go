package ir

import (
	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
	"tinygo.org/x/go-llvm"
)

func EmitString(ctx llvm.Context, s *unmangov1alpha1.Node_String_) llvm.Value {
	return llvm.ConstString(s.String_.Value, true)
}

func EmitNode(ctx llvm.Context, s *unmangov1alpha1.Node) llvm.Value {
	switch n := s.GetValue().(type) {
	case *unmangov1alpha1.Node_String_:
		return EmitString(ctx, n)
	}

	return llvm.Value{}
}

func EmitFile(ctx llvm.Context, f *unmangov1alpha1.File) llvm.Value {
	return llvm.ConstString(f.String(), true)
}
