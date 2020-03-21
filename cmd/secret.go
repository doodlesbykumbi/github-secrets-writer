package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
)

type Secret struct {
	Name  string
	Value []byte
}

func CollectSecrets(literalSources, fileSources []string) ([]Secret, error) {
	var allSecrets []Secret

	var err error

	secrets, err := secretsFromLiteralSources(literalSources)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"literal sources %v", literalSources))
	}
	allSecrets = append(allSecrets, secrets...)

	secrets, err = secretsFromFileSources(fileSources)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"file sources: %v", fileSources))
	}
	allSecrets = append(allSecrets, secrets...)

	var secretNames = map[string]bool{}
	for _, secret := range allSecrets {
		if secretNames[secret.Name] {
			return nil, fmt.Errorf(
				"multiple sources provided for secret name: %s",
				secret.Name,
			)
		}

		secretNames[secret.Name] = true
	}

	return allSecrets, nil
}

func secretsFromLiteralSources(sources []string) ([]Secret, error) {
	var secrets []Secret
	for _, s := range sources {
		secretName, secretValue, err := parseSecretSource(s)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, Secret{Name: secretName, Value: []byte(secretValue)})
	}
	return secrets, nil
}

func secretsFromFileSources(sources []string) ([]Secret, error) {
	var secrets []Secret
	for _, s := range sources {
		secretName, fPath, err := parseSecretSource(s)
		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadFile(fPath)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, Secret{Name: secretName, Value: content})
	}
	return secrets, nil
}

// parseSecretSource parses the source key=val pair into its component pieces.
// This functionality is distinguished from strings.SplitN(source, "=", 2) since
// it returns an error in the case of empty keys, values, or a missing equals sign.
func parseSecretSource(source string) (secretName, value string, err error) {
	// leading equal is invalid
	if strings.Index(source, "=") == 0 {
		return "", "", fmt.Errorf("invalid secret source %v, expected key=value", source)
	}
	// split after the first equal (so values can have the = character)
	items := strings.SplitN(source, "=", 2)
	if len(items) != 2 {
		return "", "", fmt.Errorf("invalid secret source %v, expected key=value", source)
	}
	return items[0], strings.Trim(items[1], "\"'"), nil
}
