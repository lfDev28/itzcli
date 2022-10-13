package cmd

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.ibm.com/skol/atkcli/cmd/dr"
)

var fixDoctorIssues bool = false

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Checks the environment and configuration",
	Long: `If using the init sub-command, the doctor command will initialize the
environment for first run.
`,
	PreRun: SetLoggingLevel,
	// Perform the checks on the system to make sure that ATK is OK to run
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Checking the environment...")
		configChecks := dr.AllConfigChecks
		fileChecks := dr.FileChecks
		errs := dr.DoChecks(append(fileChecks, configChecks...), fixDoctorIssues)
		if len(errs) > 0 {
			logger.Error("One or more requirements unmet; consider using --auto-fix to try to resolve them")
			return fmt.Errorf("found %d errors", len(errs))
		}
		return nil
	},
	// The usage usually prints if there is an error, but in this case we do not
	// want to print the usage.
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVarP(&fixDoctorIssues, "auto-fix", "f", false, "If true, makes best attempt to fix the issues")
}
