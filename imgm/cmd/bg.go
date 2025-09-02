package cmd

import (
	"imgg/internal/image"
	"log/slog"

	"github.com/spf13/cobra"
)

var bgCommand = &cobra.Command{
	Use: "bg",
	Run: func(cmd *cobra.Command, args []string) {
		imgP, _ := cmd.Flags().GetString("image")
		schemaP, _ := cmd.Flags().GetString("schema")
		processor, _ := cmd.Flags().GetString("processor")
		outputP, _ := cmd.Flags().GetString("output")

		if err := image.Process(image.Options{ImageP: imgP, SchemaP: schemaP, Processor: processor, OutputP: outputP}); err != nil {
			slog.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(bgCommand)
	bgCommand.Flags().StringP("image", "i", "", "Image to convert")
	bgCommand.Flags().StringP("schema", "s", "", "Path to colorscheme")
	bgCommand.Flags().StringP("processor", "p", "nn", "Processor type")
	bgCommand.Flags().StringP("output", "o", "modified.png", "Output file")
	bgCommand.MarkFlagRequired("image")
	bgCommand.MarkFlagRequired("schema")
}
