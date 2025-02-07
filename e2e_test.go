package lang_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os/exec"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"golang.org/x/net/http2"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
	"github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1/unmangov1alpha1connect"
)

var _ = Describe("LangE2eSuite", func() {
	var (
		ses     *gexec.Session
		client  unmangov1alpha1connect.ParserServiceClient
		hostlog *bytes.Buffer
	)

	BeforeEach(func() {
		var err error

		By("Starting the language host")
		hostlog = &bytes.Buffer{}
		cmd := exec.Command("bin/lang-host")
		ses, err = gexec.Start(cmd, hostlog, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		By("Waiting for the langugage host to start")
		Eventually(ses.Out).Should(gbytes.Say("Now listening on: http://localhost:5000"))

		By("Creating a parser client")
		client = unmangov1alpha1connect.NewParserServiceClient(
			newInsecureClient(),
			"http://localhost:5000",
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
		By("Terminating the host")
		Eventually(ses.Interrupt()).Should(gexec.Exit(0))
	})
})

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
				// If you're also using this client for non-h2c traffic, you may want
				// to delegate to tls.Dial if the network isn't TCP or the addr isn't
				// in an allowlist.
				return net.Dial(network, addr)
			},
			// Don't forget timeouts!
		},
	}
}
