package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jpetrucciani/poglets/pkg/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	pogletsCmd.AddCommand(clientCmd)

	clientCmd.Flags().StringP("url", "r", "", "server address i.e. ws://127.0.0.1:8000")
	clientCmd.Flags().StringP("upstream", "u", "", "upstream server i.e. http://127.0.0.1:3000")
	clientCmd.Flags().StringP("token", "t", "", "authentication token")
	clientCmd.Flags().StringP("token-from", "f", "", "read the authentication token from a file")
	clientCmd.Flags().Bool("print-token", false, "prints the token in server mode")
	clientCmd.Flags().Bool("strict-forwarding", true, "forward only to the upstream URLs specified")

	clientCmd.Flags().Bool("insecure", false, "allow the client to connect to a server without encryption")
	clientCmd.SilenceUsage = true
	clientCmd.SilenceErrors = true
}

type UpstreamParser interface {
	Parse(input string) (map[string]string, error)
}

type ArgsUpstreamParser struct {
}

func (a *ArgsUpstreamParser) Parse(input string) (map[string]string, error) {
	upstreamMap, err := buildUpstreamMap(input)

	return upstreamMap, err
}

func buildUpstreamMap(args string) (map[string]string, error) {
	items := make(map[string]string)

	entries := strings.Split(args, ",")
	for _, entry := range entries {
		kvp := strings.Split(entry, "=")
		if len(kvp) == 1 {
			items[""] = strings.TrimSpace(kvp[0])
		} else {
			items[strings.TrimSpace(kvp[0])] = strings.TrimSpace(kvp[1])
		}
	}

	for k, v := range items {
		hasScheme := (strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://"))
		if !hasScheme {
			items[k] = fmt.Sprintf("http://%s", v)
		}
	}

	if len(items) > 1 {
		return items, fmt.Errorf("cannot use multiple upstream values")
	}

	return items, nil
}

// clientCmd represents the client sub command.
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the poglets tunnel client.",
	Long:  `Start the poglets tunnel client.`,
	Example: `  # Start a secure tunnel connection over the internet to forward a
  # server running on port 3000
  poglets client \
  --url=wss://192.168.0.101 \
  --upstream=http://127.0.0.1:3000 \
  --token TOKEN
  
  # Start an insecure tunnel connection over your local network
  poglets client \
    --url=ws://192.168.0.101:80 \
    --upstream=http://127.0.0.1:3000 \
    --token TOKEN \
    --insecure

  Note: You can pass the --token argument followed by a token value to both the server and client to prevent unauthorized connections to the tunnel.`,
	RunE:          runClient,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// runClient does the actual work of reading the arguments passed to the client sub command.
func runClient(cmd *cobra.Command, _ []string) error {

	upstream, err := cmd.Flags().GetString("upstream")
	if err != nil {
		return errors.Wrap(err, "failed to get 'upstream' value")
	}

	if len(upstream) == 0 {
		return errors.New("upstream is missing in the client argument")
	}

	argsUpstreamParser := ArgsUpstreamParser{}
	upstreamMap, err := argsUpstreamParser.Parse(upstream)
	if err != nil {
		return err
	}

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return errors.Wrap(err, "failed to get 'url' value.")
	}

	insecure, err := cmd.Flags().GetBool("insecure")
	if err != nil {
		return errors.Wrap(err, "failed to get 'insecure' value.")
	}

	tokenFile, err := cmd.Flags().GetString("token-from")
	if err != nil {
		return errors.Wrap(err, "failed to get 'token-from' value.")
	}

	var token string
	if len(tokenFile) > 0 {
		fileData, err := os.ReadFile(tokenFile)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to load file: %s", tokenFile))
		}

		// new-lines will be stripped, this is not configurable and is to
		// make the code foolproof for beginners
		token = strings.TrimRight(string(fileData), "\n")
	} else {
		tokenVal, err := cmd.Flags().GetString("token")
		if err != nil {
			return errors.Wrap(err, "failed to get 'token' value.")
		}
		token = tokenVal
	}

	printToken, err := cmd.Flags().GetBool("print-token")
	if err != nil {
		return errors.Wrap(err, "failed to get 'print-token' value.")
	}

	strictForwarding, err := cmd.Flags().GetBool("strict-forwarding")
	if err != nil {
		return errors.Wrap(err, "failed to get 'strict-forwarding' value.")
	}

	if len(token) > 0 && printToken {
		log.Printf("Token: %q", token)
	}

	if len(url) == 0 {
		return fmt.Errorf("--url is required")
	}

	if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
		return fmt.Errorf("--url should be prefixed with ws:// (insecure) or wss:// (secure)")
	}

	if strings.HasPrefix(url, "ws://") {
		if !insecure {
			fmt.Print(`[================================== Warning ==================================]

You are trying to connect to an poglets server without any form of encryption.

You can disable this warning with --insecure, but be aware that your data 
could be read by a third-party.

You may benefit from using poglets PRO which has options for automatic 
encryption.

[=============================================================================]
`)
			os.Exit(1)
		} else {
			fmt.Printf("Warning: running in insecure mode, without encryption.\n")
		}
	}

	log.Printf("Starting client - version %s", getVersion())
	for k, v := range upstreamMap {
		log.Printf("Upstream: %s => %s\n", k, v)
	}

	pogletsClient := client.Client{
		Remote:           url,
		UpstreamMap:      upstreamMap,
		Token:            token,
		StrictForwarding: strictForwarding,
	}

	if err := pogletsClient.Connect(); err != nil {
		return err
	}

	return nil
}
