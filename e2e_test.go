package lang_test

import (
	"bytes"
	"context"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"github.com/unmango/go/vcs/git"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	unmangov1alpha1 "github.com/unstoppablemango/lang/pkg/io/unmango/v1alpha1"
)

var _ = Describe("LangE2eSuite", func() {
	Describe("Lang host", func() {
		var (
			ses     *gexec.Session
			client  unmangov1alpha1.ParserServiceClient
			hostlog *bytes.Buffer
		)

		BeforeEach(func(ctx context.Context) {
			dir, err := git.Root(ctx)
			Expect(err).NotTo(HaveOccurred())
			socket := filepath.Join(dir, "lang-host.sock")

			By("Starting the language host")
			hostlog = &bytes.Buffer{}
			cmd := exec.Command("bin/lang-host", socket)
			ses, err = gexec.Start(cmd, hostlog, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for the langugage host to start")
			Eventually(ses.Out).Should(gbytes.Say("Application started"))

			By("Creating a parser client")
			conn, err := grpc.NewClient("unix:"+socket,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			Expect(err).NotTo(HaveOccurred())
			client = unmangov1alpha1.NewParserServiceClient(conn)
		})

		It("should work", func(ctx context.Context) {
			res, err := client.Parse(ctx, &unmangov1alpha1.ParseRequest{
				Text: `"test"`,
			})

			Expect(err).NotTo(HaveOccurred(), hostlog.String())
			Expect(res.File.Nodes).To(HaveLen(1))
			Expect(res.File.Nodes[0].GetString_().Value).To(Equal("test"))
		})

		AfterEach(func() {
			By("Terminating the host")
			Eventually(ses.Interrupt()).Should(gexec.Exit(0))
		})
	})
})
