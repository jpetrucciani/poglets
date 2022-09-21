package cmd

import (
	"fmt"
	"os"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func init() {
	pogletsCmd.AddCommand(versionCmd)
}

// pogletsCmd represents the base command when called without any sub commands.
var pogletsCmd = &cobra.Command{
	Use:   "poglets",
	Short: "Expose your local ports to the Internet.",
	Long:  ``,
	Run:   runPoglets,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the poglets version information.",
	Run:   parseBaseCommand,
}

func getVersion() string {
	if len(Version) != 0 {
		return Version
	}
	return "dev"
}

func parseBaseCommand(_ *cobra.Command, _ []string) {
	printLogo()

	fmt.Println("Version:", getVersion())
	fmt.Println("Git Commit:", GitCommit)
	os.Exit(0)
}

// Execute adds all child commands to the root command(pogletsCmd) and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the pogletsCmd.
func Execute(version, gitCommit string) error {

	// Get Version and GitCommit values from main.go.
	Version = version
	GitCommit = gitCommit

	if err := pogletsCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func runPoglets(cmd *cobra.Command, args []string) {
	printLogo()
	cmd.Help()
}

func printLogo() {
	pogletsLogo := aec.WhiteF.Apply(pogletsFigletStr)
	fmt.Println(pogletsLogo)
}

const pogletsFigletStr = `                _        _      
 ___  ___  ___ | | ___ _| |_ ___
| . \/ . \/ . || |/ ._> | | <_-<
|  _/\___/\_. ||_|\___. |_| /__/
|_|       <___'                 
`
