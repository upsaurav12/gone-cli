package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/upsaurav12/bootstrap/api"
)

var aiCmd = cobra.Command{
	Use:   "ai",
	Short: "user enters the prompt to the ai",
	Long:  "user enters the prompt to the ai",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.CallingApiToGroq(userPrompt)
		if err != nil {
			fmt.Println("error: ", err)
		}

		fmt.Println("result: ", result)
	},
}

var userPrompt string

func init() {

	rootCmd.AddCommand(&aiCmd)

	aiCmd.Flags().StringVar(&userPrompt, "prompt", " ", "user enter the prompt")
}
