package cmd

import (
	"fmt"
	"reflect"

	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pipelineID string
var reservationID string

// showCmd represents the version command
var showCmd = &cobra.Command{
	Use:   ShowAction,
	Short: "Shows the details of the requested single object",
	Long:  `Shows the details of the requested single object.`,
}

var showReservationCmd = &cobra.Command{
	Use:   ReservationResource,
	Short: "Shows the details of the specific reservation",
	Long: `
Shows the details of the specific IBM Technology Zone reservation.
`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
        logger.Debug("Listing your reservations...")
        rez, err := GetReservationDetails(reservationID)
        if err != nil {
            return err
        }
        w := techzone.NewModelWriter(reflect.TypeOf(techzone.Reservation{}).Name(), GetFormat(cmd))
        return w.WriteOne(cmd.OutOrStdout(), rez)
    },
}


var showPipelinesCmd = &cobra.Command{
	Use:   PipelineResource,
	Short: fmt.Sprintf("Shows the details of the specific %s from the %s catalog", PipelineResource, TechZoneShort),
	Long: `
Shows the details of the IBM Technology Zone pipeline specified by the
--pipeline-id flag.
`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("Getting your pipeline...")
		if err := AssertFlag(pipelineID, NotNull, "you must specify a valid pipeline ID using --pipeline-id"); err != nil {
			return err
		}
		apiConfig, err := LoadApiClientConfig(configuration.Backstage)
		if err != nil {
			return err
		}
		svc, err := solutions.NewWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		w := solutions.NewSolutionWriter(GetFormat(cmd))
		sol, err := svc.Get(pipelineID)
		if err != nil {
			return err
		}
		return w.Write(cmd.OutOrStdout(), sol)
	},
}

var showEnvironmentCmd = &cobra.Command{
	Use:   EnvironmentResource,
	Short: fmt.Sprintf("Shows the details of the %s %s", TechZoneShort, EnvironmentResource),
	Long: `
Shows the details of the IBM Technology Zone environments.
`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Listing the %s environments...", TechZoneFull)
		return nil
	},
	Hidden: true,
}

func init() {

	// Add the parameters to the show commands...
	showPipelinesCmd.Flags().StringVar(&pipelineID, "pipeline-id", "", "ID of the build in the catalog")
	showReservationCmd.Flags().StringVar(&reservationID, "reservation-id", "", "ID of the reservation")

	showCmd.AddCommand(showReservationCmd)
	showCmd.AddCommand(showEnvironmentCmd)
	showCmd.AddCommand(showPipelinesCmd)

	rootCmd.AddCommand(showCmd)
}

// Helper function to be used by the show command and deploy command
func GetReservationDetails(reservationID string) (*techzone.Reservation, error) {
    if len(reservationID) == 0 {
        return nil, fmt.Errorf("the reservation id is empty")
    }
    apiConfig, err := LoadApiClientConfig(configuration.TechZone)
    if err != nil {
        return nil, err
    }
    svc, err := techzone.NewReservationWebServiceClient(apiConfig)
    if err != nil {
        return nil, errors.Wrap(err, "could not create web service client")
    }
    rez, err := svc.Get(reservationID)
    if err != nil {
        return nil, err
    }
    return rez, nil
}
