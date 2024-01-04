package main

import (
	"fmt"
	"os"

	"github.com/gabriel-antonelli/web-redable-pdf/internal/create_html"
	"github.com/gabriel-antonelli/web-redable-pdf/internal/create_pdf"
	"github.com/spf13/cobra"
)

var (
	outputPDF      string
	createSeparate bool
	rootCmd        = &cobra.Command{
		Use:   "web-readable-pdf [flags] [web-page1] [web-page2] [web-page3]",
		Short: "Generate a web-readable PDF from multiple web pages",
		Long: `web-readable-pdf is a tool to generate a web-readable PDF from multiple web pages.
Specify one or more web pages and an output PDF file.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tempDir := createhtml.GenerateRedableHtml(args)
			if createSeparate {
				createpdf.CreateSeparatedPDFFiles(tempDir, outputPDF)
			} else {
				createpdf.CreateCombinedPDF(tempDir, outputPDF)
			}
		},
	}
)

func main() {
	rootCmd.Flags().StringVarP(&outputPDF, "output", "o", "ARTICLENAME", "Specify the output PDF file name")
	rootCmd.Flags().BoolVarP(&createSeparate, "separate", "s", false, "Create separate PDFs for each page")
	rootCmd.Flags().BoolP("help", "h", false, "Help for web-readable-pdf")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
