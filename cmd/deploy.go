package cmd

import (
	"fmt"

	"github.com/cloud-native-toolkit/itzcli/pkg"
	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	"github.com/cloud-native-toolkit/itzcli/pkg/solutions"
	"github.com/cloud-native-toolkit/itzcli/pkg/techzone"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var resId string



var deployCmd = &cobra.Command{
	Use:    DeployAction,
	Short:  "Deploys a pipeline in a cluster",
	Long:   "Deploys a pipeline into an existing cluster",
	PreRun: SetLoggingLevel,
}

var deployPipelineCmd = &cobra.Command{
	Use:   PipelineResource,
	Short: "Deploys the given pipeline to the specified cluster",
	Long: `
Deploys the given pipeline to the cluster specified by --cluster-api-url ("-a").
The pipeline is identified by a UUID and can be found by executing the command:

    itz list pipelines

To view the current pipelines. With the pipeline ID, you can deploy the pipeline
to a cluster with the given API endpoint ("--cluster-api-url" or "-a"), and a 
username/password of a user with permissions to create Pipelines and PipelineRuns.

Example:

    itz deploy pipeline -p c567d9bd-5f0f-4254-bce1-c40ef1fedc0c \
      -a http://cluster.api.example.com \
      -u clusteruser \
      -P mysecretpassword

`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		SetLoggingLevel(cmd, args)

		
        if err := AssertFlag(pipelineID, NotNull, "you must specify a valid pipeline ID using --pipeline-id"); err != nil {
            return err
        }

        // If a reservation ID is provided, get the cluster details from the reservation
        if resId != "" {
            rez, err := GetReservationDetails(resId)
            if err != nil {
                return err
            }

            clusterDetails, err := getClusterAdminCredentials(rez)
            if err != nil {
                return err
            }

			logger.Tracef("Cluster Details: %v", clusterDetails)
			logger.Tracef(clusterDetails.ClusterURL, clusterDetails.ClusterAdminUsername, clusterDetails.ClusterAdminPassword)

            clusterURL = clusterDetails.ClusterURL
            clusterUsername = clusterDetails.ClusterAdminUsername
            clusterPassword = clusterDetails.ClusterAdminPassword
        } else {
            if err := AssertFlag(clusterURL, ValidURL, "you must specify a valid URL using --cluster-api-url"); err != nil {
                return err
            }

            if err := AssertFlag(clusterUsername, NotNull, "you must specify a valid username using --cluster-username"); err != nil {
                return err
            }

            if err := AssertFlag(clusterPassword, NotNull, "you must specify a valid value using --cluster-password"); err != nil {
                return err
            }
        }

        return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debugf("Deploying your pipeline %s to cluster %s...", pipelineID, clusterURL)
		// Go get the pipeline component from the catalog
		apiConfig, err := LoadApiClientConfig(configuration.Backstage)
		if err != nil {
			return err
		}
		svc, err := solutions.NewWebServiceClient(apiConfig)
		if err != nil {
			return errors.Wrap(err, "could not create web service client")
		}
		sol, err := svc.Get(pipelineID)
		if err != nil {
			return err
		}

		// Look up the pipeline location, and get that.
		pipelineURI, found := LookupAnnotation(sol, PipelineAnnotation)
		if !found {
			return fmt.Errorf("could not find the pipeline location from catalog entry with id: %s", pipelineID)
		}

		pipelineRunURI, found := LookupAnnotation(sol, PipelineRunAnnotation)
		if !found {
			// try guesssing...
			logger.Infof("No pipeline run location was found, attempting to guess...")
			pipelineRunURI, err = pkg.AppendToFilename(pipelineURI, "-run")
			if err != nil {
				return nil
			}
			logger.Debugf("Guessed %s as the pipeline run location.", pipelineRunURI)
		}

		execArgs := PipelineExecArgs{
			PipelineURI:     pipelineURI,
			PipelineRunURI:  pipelineRunURI,
			ClusterURL:      clusterURL,
			ClusterUsername: clusterUsername,
			ClusterPassword: clusterPassword,
			AdditionalArgs:  args,
			AcceptDefaults:  acceptDefaults,
			UseContainer:    useContainer,
		}
		return ExecutePipeline(cmd, execArgs)
	},
}

func LookupAnnotation(sol *solutions.Solution, name string) (string, bool) {
	if sol.Entity.Metadata.Annotations != nil && len(sol.Entity.Metadata.Annotations) > 0 {
		val, found := sol.Entity.Metadata.Annotations[name]
		return val, found
	}
	return "", false
}

func init() {
	deployPipelineCmd.Flags().StringVarP(&pipelineID, "pipeline-id", "p", "", "ID of the pipeline from the catalog (required)")
	deployPipelineCmd.Flags().StringVarP(&clusterURL, "cluster-api-url", "a", "", "The URL of the target cluster (required)")
	deployPipelineCmd.Flags().StringVarP(&clusterUsername, "cluster-username", "u", "", "A username to login to the target cluster (required)")
	deployPipelineCmd.Flags().StringVarP(&clusterPassword, "cluster-password", "P", "", "A password to login to the target cluster (required)")
	deployPipelineCmd.Flags().BoolVarP(&acceptDefaults, "accept-defaults", "d", false, "Accept defaults for pipeline parameters without asking (optional)")
	deployPipelineCmd.Flags().StringVarP(&resId, "reservation-id", "r", "", "The reservation id to deploy it to (optional)")

	// for _, pname := range []string{"pipeline-id", "cluster-api-url", "cluster-username", "cluster-password"} {
	// 	if err := deployPipelineCmd.MarkFlagRequired(pname); err != nil {
	// 		panic(fmt.Sprintf("could not mark %s required", pname))
	// 	}
	// }

	// deployPipelineCmd.Flags().BoolVarP(&useContainer, "use-container", "c", DefaultUseContainer, "If true, the commands run in a container")
	deployCmd.AddCommand(deployPipelineCmd)

	rootCmd.AddCommand(deployCmd)
}


// Return type interface
type ClusterDetails struct {
	ClusterURL string
	ClusterAdminUsername string
	ClusterAdminPassword string
}

// Function to get the Cluster Admin Username and Password from the response
// Loop throug the service links property and get the values from the labels: "Cluster Admin Username" and "Cluster Admin Password"}
// This function will be used in another command which will be to show and deploy in one command.
func getClusterAdminCredentials(rez *techzone.Reservation) (ClusterDetails, error) {
	var details ClusterDetails


	for _, link := range rez.ServiceLinks {
		logger.Tracef("Link: %v", link.Label)
		if link.Label == "Cluster Admin Username" {
			username, ok := link.Data.(string)
			if !ok {
				return details, errors.New("could not convert username to string")
			}
			details.ClusterAdminUsername = username
		}
		if link.Label == "Cluster Admin Password" {
			password, ok := link.Data.(string)
			if !ok {
				return details, errors.New("could not convert password to string")
			}
			details.ClusterAdminPassword = password
		} 
		if link.Label == "API URL" {
			url, ok := link.Data.(string)
			if !ok {
				return details, errors.New("could not convert url to string")
			}
			details.ClusterURL = url
		}
	}
	return details, nil
	
}

