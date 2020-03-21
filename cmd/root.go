package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	gsw "github.com/doodlesbykumbi/github-secrets-writer/pkg"
)

var owner string
var repo string
var fromFile []string
var fromLiteral []string
var secrets []Secret

func fatal(err error) {
	_, _ = os.Stderr.Write([]byte(fmt.Sprintf("ERROR: %s\n", err)))
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:                   "github-secrets-writer --owner=owner --repo=repo [--from-literal=secretName1=secretValue1] [--from-file=secretName2=/path/to/secretValue2]",
	DisableFlagsInUseLine: true,
	Short:                 "Create or update multiple Github secrets sourced from literal values or files.",
	Long: `Create or update multiple Github secrets sourced from literal values or files.

Key/value pairs representing a secret name and the source of the secret value are provided via the flags --from-file and --from-literal. Depending on the key/value pairs specified a single invocation may carry out zero or more writes to the Github secrets of the repository.

NOTE: An OAuth token **must** be provided via the 'GITHUB_TOKEN' environment variable, this is used to authenticate to the Github API. Access tokens require 'repo' scope for private repos and 'public_repo' scope for public repos. GitHub Apps must have the 'secrets' permission to use the API. Authenticated users must have collaborator access to a repository to create, update, or read secrets.
`,
	Example: `# Write a single secret from a literal value
github-secrets-writer --owner=owner --repo=repo --from-literal=secretName1=secretValue1

# Write a single secret from a file
github-secrets-writer --owner=owner --repo=repo --from-file=secretName1=/path/to/secretValue1
  
# Write multiple secrets, one from a literal value and one from a file
github-secrets-writer --owner=owner --repo=repo --from-literal=secretName1=secretValue1 --from-file=secretName2=/path/to/secretValue2`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				fatal(err)
			}
		}()
		secretWriter := gsw.NewSecretWriter(viper.GetString("token"))
		var hasFailures bool

		fmt.Printf("Write results:\n\n")
		for _, secret := range secrets {
			result, wErr := secretWriter.Write(owner, repo, secret.Name, secret.Value)
			if wErr != nil {
				hasFailures = true
				result = wErr.Error()
			}

			fmt.Printf("%s: %s\n", secret.Name, result)
		}

		if hasFailures {
			err = fmt.Errorf("encountered some failures, see above")
			return
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("token") {
			return fmt.Errorf("envvar not set: GITHUB_TOKEN")
		}

		var err error
		secrets, err = CollectSecrets(fromLiteral, fromFile)
		if err != nil {
			return err
		}

		if len(secrets) == 0 {
			return fmt.Errorf("no secret name-source pairs provided, you must specify at least one of --from-literal or --from-file")
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&owner, "owner", "", "owner of the repository e.g. an organisation or user (required)")
	rootCmd.Flags().StringVar(&repo, "repo", "", "name of the repository (required)")
	rootCmd.Flags().StringArrayVar(&fromFile, "from-file", []string{}, "specify secret name and literal value pairs e.g. secretname=somevalue (zero or more)")
	rootCmd.Flags().StringArrayVar(&fromLiteral, "from-literal", []string{}, "specify secret name and source file pairs e.g. secretname=somefile (zero or more)")

	_ = rootCmd.MarkFlagRequired("owner")
	_ = rootCmd.MarkFlagRequired("repo")
}

// initConfig reads in ENV variables.
func initConfig() {
	_ = viper.BindEnv("token", "GITHUB_TOKEN")
}
