package cmd

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var reservationId string

var statusCmd = &cobra.Command{
	Use: StatusAction,
	Short: "Shows the status of the requested single object",
	Long:  `Shows the status of the requested single object.`,
}


var reservationStatusCmd = &cobra.Command{
	Use:   ReservationResource,
	Short: "Shows the status of the specific reservation",
	Long: `
		Shows the status of the specific IBM Technology Zone reservation.
		`,
	PreRun: SetLoggingLevel,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Getting your reservation status...")
		rez, err := GetReservationDetails(reservationId)
		if err != nil {
			return err
		}
		
		// Just print to stdout for now, as the vm script just parses the output.
		fmt.Println(rez.Status)

		return nil
	},
}

// TODO: Implement the pipelineStatusCmd command


func init() { 
	reservationStatusCmd.Flags().StringVarP(&reservationId, "reservation-id", "r", "", "The ID of the reservation to show the status of")
	statusCmd.AddCommand(reservationStatusCmd)

	rootCmd.AddCommand(statusCmd)
}
