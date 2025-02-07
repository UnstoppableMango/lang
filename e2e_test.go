package lang_test

import (
	"context"
	"net/http"
	"os/exec"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
	"github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1/unmangov1alpha1connect"
)

var _ = Describe("LangE2eSuite", func() {
	var (
		ses    *gexec.Session
		client unmangov1alpha1connect.ParserServiceClient
	)

	BeforeEach(func() {
		var err error
		cmd := exec.Command("bin/host")
		ses, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		client = unmangov1alpha1connect.NewParserServiceClient(
			http.DefaultClient,
			"",
			connect.WithGRPC(),
		)
	})

	It("should work", func(ctx context.Context) {
		res, err := client.Parse(ctx, connect.NewRequest(&unmangov1alpha1.ParseRequest{
			Text: `"test"`,
		}))

		Expect(err).NotTo(HaveOccurred())
		Expect(res.Msg.File).NotTo(BeNil())
	})

	AfterEach(func() {
		Eventually(ses.Terminate().Exited).Should(gexec.Exit(0))
	})
})
