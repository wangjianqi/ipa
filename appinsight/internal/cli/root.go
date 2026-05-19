package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wangjianqi/AppInsight/internal/ipa"
	"github.com/wangjianqi/AppInsight/internal/ipatool"
	"github.com/wangjianqi/AppInsight/internal/output"
	"github.com/wangjianqi/AppInsight/internal/report"
	"github.com/wangjianqi/AppInsight/internal/system"
)

var Version = "0.1.0"

var jsonFlag bool

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "appinsight",
		Short: "AppInsight CLI - iOS App analysis tool for developers",
		Long:  "AppInsight CLI is a developer-oriented iOS App analysis tool. It can search and download App Store IPAs, analyze visible information including Info.plist, permissions, frameworks, resources, and generate structured reports.",
	}

	root.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output in JSON format")

	root.AddCommand(newDoctorCommand())
	root.AddCommand(newSearchCommand())
	root.AddCommand(newFetchIOSCommand())
	root.AddCommand(newAnalyzeIPACommand())
	root.AddCommand(newReportCommand())

	return root
}

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and dependencies",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := system.RunDoctor(Version)
			return output.PrintJSON(output.NewOK("doctor", result))
		},
	}
}

func newSearchCommand() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search App Store apps via ipatool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			keyword := args[0]

			if !system.CheckTool("ipatool").Available {
				errMsg := system.MissingToolError("ipatool")
				return output.PrintJSON(output.NewError("search", errMsg))
			}

			resp, err := ipatool.Search(keyword, limit)
			if err != nil {
				return output.PrintJSON(output.NewError("search", err.Error()))
			}

			return output.PrintJSON(output.NewOK("search", resp))
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "max number of results")

	return cmd
}

func newFetchIOSCommand() *cobra.Command {
	var bundleID string
	var outputDir string
	var purchase bool

	cmd := &cobra.Command{
		Use:   "fetch-ios",
		Short: "Download IPA from App Store via ipatool",
		RunE: func(cmd *cobra.Command, args []string) error {
			if bundleID == "" {
				return output.PrintJSON(output.NewError("fetch-ios", "--bundle-id is required"))
			}

			if !system.CheckTool("ipatool").Available {
				errMsg := system.MissingToolError("ipatool")
				return output.PrintJSON(output.NewError("fetch-ios", errMsg))
			}

			resp, err := ipatool.Fetch(bundleID, outputDir, purchase)
			if err != nil {
				return output.PrintJSON(output.NewError("fetch-ios", err.Error()))
			}

			return output.PrintJSON(output.NewOK("fetch-ios", resp))
		},
	}

	cmd.Flags().StringVar(&bundleID, "bundle-id", "", "bundle identifier of the app")
	cmd.Flags().StringVar(&outputDir, "output", "./downloads", "output directory for downloaded IPA")
	cmd.Flags().BoolVar(&purchase, "purchase", false, "purchase the app if needed")

	return cmd
}

func newAnalyzeIPACommand() *cobra.Command {
	var outputFile string

	cmd := &cobra.Command{
		Use:   "analyze-ipa <ipa-path>",
		Short: "Analyze an IPA file for visible information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ipaPath := args[0]

			if !system.CheckTool("plutil").Available {
				errMsg := system.MissingToolError("plutil")
				return output.PrintJSON(output.NewError("analyze-ipa", errMsg))
			}

			result, err := ipa.Analyze(ipaPath)
			if err != nil {
				return output.PrintJSON(output.NewError("analyze-ipa", err.Error()))
			}

			if outputFile != "" {
				if err := output.WriteDataToFile(result, outputFile); err != nil {
					return output.PrintJSON(output.NewError("analyze-ipa", fmt.Sprintf("failed to write output file: %v", err)))
				}
			}

			return output.PrintJSON(output.NewOK("analyze-ipa", result))
		},
	}

	cmd.Flags().StringVar(&outputFile, "output", "", "write analysis result to file")

	return cmd
}

func newReportCommand() *cobra.Command {
	var format string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "report <analysis-json>",
		Short: "Generate a Markdown report from analysis JSON",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			analysisPath := args[0]

			analysis, err := report.LoadAnalysis(analysisPath)
			if err != nil {
				return output.PrintJSON(output.NewError("report", err.Error()))
			}

			switch format {
			case "markdown":
				md := report.GenerateMarkdown(analysis)
				if outputFile != "" {
					if err := report.WriteReport(md, outputFile); err != nil {
						return output.PrintJSON(output.NewError("report", fmt.Sprintf("failed to write report: %v", err)))
					}
				}
				fmt.Print(md)
				return nil
			case "html":
				h := report.GenerateHTML(analysis)
				if outputFile != "" {
					if err := report.WriteReport(h, outputFile); err != nil {
						return output.PrintJSON(output.NewError("report", fmt.Sprintf("failed to write report: %v", err)))
					}
				}
				fmt.Print(h)
				return nil
			default:
				return output.PrintJSON(output.NewError("report", fmt.Sprintf("unsupported format: %s (supported: markdown, html)", format)))
			}
		},
	}

	cmd.Flags().StringVar(&format, "format", "markdown", "output format (markdown, html)")
	cmd.Flags().StringVar(&outputFile, "output", "", "write report to file")

	return cmd
}
