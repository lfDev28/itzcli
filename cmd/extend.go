package cmd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var extendDate string

func generateJSONRequest(email, extensionDate, reservationId string) string {
	return fmt.Sprintf(`{"IBMID":"%s","requestType":"ibmcloud-2","extensionDate":"%s","reservationId":"%s","id":"%s","environmentId":"%s"}`, email, extensionDate, reservationId, reservationId, reservationId)

}




func getDefaultExtendDate(endDate string) (string, error) {
    // Parse the endDate string into a time.Time value
    t, err := time.Parse("2006-01-02 15:04:05", endDate)
    if err != nil {
        return "", fmt.Errorf("unable to parse date: %v", err)
    }

    // Add 24 hours to the end date, then subtract 1 minute to ensure it's strictly less than 2 days
    t = t.Local().Add(time.Hour*48 - time.Minute)

    // Format the new date back into a string with milliseconds
    newEndDate := t.UTC().Format("2006-01-02T15:04:05.000Z")

    return newEndDate, nil
}

var extendCmd = &cobra.Command{
	Use:    ExtendAction,
	Short:  "Allows you to Extend environments",
	Long:   "Allows you to Extend environments",
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Extending your environment...")
		// Need to load all of the filter flags and then use them when reserving the environment

		apiConfig, err := LoadApiClientConfig(configuration.TechZone)
		if err != nil {
			return err
		}

		svc, err := techzone.NewReservationWebServiceClient(apiConfig)
		if err != nil {
			return err
		}

		// Allowing for a default extend date of 48 hours from the current end date
		if extendDate == "" {
			rez, err := svc.Get(reservationID)
			if err != nil {
				return err
			}

			// Print the extend date in readable format
			logger.Debugf("Current end date: %v", rez.ProvisionUntil)

			extendDate, err = getDefaultExtendDate(rez.ProvisionUntil)
			if err != nil {
				return err
			}

			logger.Debugf("Default extend date: %v", extendDate)
		}

		body := generateJSONRequest(email, extendDate, reservationID)

		extension, err := svc.Extend(reservationID, body)
		if err != nil {
			return err
		}




		logger.Debugf("Extension: %v", extension)

		w := techzone.NewModelWriter(reflect.TypeOf(techzone.Extension{}).Name(), GetFormat(cmd))

		logger.Debugf("Writing extension to output...")
		

		return w.WriteOne(cmd.OutOrStdout(), extension)
	
	},

}

func init() {
	extendCmd.Flags().StringVarP(&email, "email", "e", "", "The email of the user to reserve the environment for")
	extendCmd.Flags().StringVarP(&reservationID, "reservation-id", "r", "", "The ID of the reservation to extend")
	extendCmd.Flags().StringVarP(&extendDate, "end", "n", "", "The end time of the reservation")

	rootCmd.AddCommand(extendCmd)
}

