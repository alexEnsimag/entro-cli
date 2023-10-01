package cmd

import (
	"alex/entro-cli/pkg/entro"
	"alex/entro-cli/pkg/report"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "entro-cli",
	Short: "entro-cli - a simple CLI to query entro",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Requests a report",
	Run: func(cmd *cobra.Command, args []string) {
		createReport()
	},
}

var getStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get the status of a report",
	Run: func(cmd *cobra.Command, args []string) {
		reportID := args[0]
		getReportStatus(report.ID(reportID))
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a report locally",
	Run: func(cmd *cobra.Command, args []string) {
		reportID := args[0]
		downloadReport(report.ID(reportID))
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

func getAWSCredentials() (accessKeyID, secretAccessKey, sessionToken, region string) {
	// FIXME (alex): verify env variables are not empty and have a valid format
	return os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		os.Getenv("AWS_SESSION_TOKEN"),
		os.Getenv("AWS_REGION")
}

func createReport() {
	accessKeyID, secretAccessKey, sessionToken, region := getAWSCredentials()
	client := entro.NewClient(accessKeyID, secretAccessKey, sessionToken, region)
	id, err := client.CreateReport()
	if err != nil {
		fmt.Println("Failed to create report:", err.Error())
		return
	}
	fmt.Println("Successfully requested report:", id)
}

func getReportStatus(id report.ID) {
	accessKeyID, secretAccessKey, sessionToken, region := getAWSCredentials()
	client := entro.NewClient(accessKeyID, secretAccessKey, sessionToken, region)
	status, err := client.GetReportStatus(id)
	if err != nil {
		fmt.Println("Failed to create report:", err.Error())
		return
	}
	fmt.Println("Report status is:", status)
}

func downloadReport(id report.ID) {
	accessKeyID, secretAccessKey, sessionToken, region := getAWSCredentials()
	client := entro.NewClient(accessKeyID, secretAccessKey, sessionToken, region)
	data, err := client.GetReport(id)
	if err != nil {
		fmt.Println("Failed to download report:", err.Error())
		return
	}

	filePath := "./" + string(id)
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err.Error())
		return
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Failed to write report:", err.Error())
		return
	}
	fmt.Println("Report was successfully saved:", filePath)
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(getStatusCmd)
	rootCmd.AddCommand(downloadCmd)
}
