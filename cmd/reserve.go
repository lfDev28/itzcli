package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cloud-native-toolkit/itzcli/pkg/configuration"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// For now, we can hard code a users array with their desired amount of environments and the deployment configurations / MAS Versions for each deployment

// Later, we can make this customizable with flags etc...

var postUrl = "https://api.techzone.ibm.com/api/reservation/ibmcloud-2"

var email string
var region string
var purpose string
var description string
var start string
var end string
var api_key string


//TODO: Figure out how to add an auto extension to the limit of the reservation, or figure out how to do it in a separate script / command.


var reserveCmd = &cobra.Command{
	Use:    ReserveAction,
	Short:  "Allows you to reserve environments",
	Long:   "Allows you to reserve environments",
	PreRun: SetLoggingLevel,
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Debug("Reserving your environment...")
		// Need to load all of the filter flags and then use them when reserving the environment


		apiConfig, err := LoadApiClientConfig(configuration.TechZone)
		if err != nil {
			return err
		}

		if api_key == "" {
			api_key = apiConfig.Token
		}

		// We need to generate a few things to fulfill the JSON request such as the start and end times, the user, the region, the template etc...
		// We will use the users array to generate the JSON request for each user and then send the request to the TechZone API to reserve the environment
		// We will need to loop through the users array and generate the JSON request for each user and then send the request to the TechZone API to reserve the environment

		logger.Debug(fmt.Sprintf("Reserving environment for %s in %s region", email, region))
		logger.Debug(fmt.Sprintf("Purpose: %s", purpose))
		logger.Debug(fmt.Sprintf("Description: %s", description))
		logger.Debug(fmt.Sprintf("ApiKey: %s", api_key))

		JSON_Request := getJSONRequest(purpose, description, email, region)


		// We will need to send a POST request to the TechZone API to reserve the environment and use the header with Bearer and the API Key
		
		r, err := http.NewRequest("POST", postUrl, bytes.NewBuffer([]byte(JSON_Request)))

		if err != nil {
			logger.Error("Error creating request")
			return err
		}

		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Authorization", "Bearer " + api_key)

		logger.Debug(fmt.Sprintf("Sending request to %s", postUrl))

		client := &http.Client{}
		resp, err := client.Do(r)

		if err != nil {
			logger.Error(fmt.Sprintf("Error sending request: %s", err))
			return err
		}
		

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)


		if err != nil {
			logger.Error(fmt.Sprintf("Error reading response body: %s", err))
			return err
		}

		// Need to write the response to stdout so that other scripts can use the response
		fmt.Println(string(body))
		logger.Debug(fmt.Sprintf("Response: %s", body))

		logger.Debug("Environment reservation has started")

		// This function can be run on a Friday for example, it will reserve all the environments for our team, and then on Sunday, we can run the ITZ Deploy command to deploy the environments.
		// This will allow us to have the environments ready for the week ahead. When we create the ITZ Deploy command, we can make it so that it will deploy the environments that have been reserved for us.		

		return nil
	},
}

func init() {
	reserveCmd.Flags().StringVarP(&email, "email", "e", "", "The email of the user to reserve the environment for")
	reserveCmd.Flags().StringVarP(&region, "region", "r", "", "The region to reserve the environment in")
	reserveCmd.Flags().StringVarP(&purpose, "purpose", "p", "", "The purpose of the environment")
	reserveCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the environment")
	reserveCmd.Flags().StringVarP(&start, "start", "s", "", "The start time of the reservation")
	reserveCmd.Flags().StringVarP(&end, "end", "n", "", "The end time of the reservation")
	reserveCmd.Flags().StringVarP(&api_key, "api_key", "a", "", "The API Key to use for the request")

	rootCmd.AddCommand(reserveCmd)
}


func getJSONRequest(purpose, description, email, region string) string {
    // Map that maps the purpose to the number of hours to add to current time.
	// TODO: Change to add the current time - 1 minute rather than a whole hour. For example 11Hrs 59 minutes rather than 12 hours.
    purposeHours := map[string]int{
        "Test": 11,
        "Practice / Self-Education": 47,
    }


	// Replace with .UTC() for production
	now := time.Now()
	
    // Get the number of hours to add for the given purpose
    hours, exists := purposeHours[purpose]
    if !exists {
        hours = 0 // Default to 0 if the purpose is not in the map
    }

	// We need to add the hours to the current time to get the end time

	// Format the start and end times
	start := now.Format("2006-01-02T15:04:05.000Z")
	end = now.Add(time.Hour * time.Duration(hours)).Format("2006-01-02T15:04:05.000Z")

	


	


    JSON_Request := fmt.Sprintf(`{
		"id": null,
		"name": "OpenShift Cluster (VMware on IBM Cloud) - UPI - Public",
		"purpose": "%s",
		"customer": "",
		"opportunity": [],
		"opportunityProduct": "",
		"description": "%s",
		"start": "%s",
		"end": "%s",
		"notes": "Notes",
		"user": "%s",
		"template": "vmware-openshift-upi",
		"infrastructure": "ibmcloud-2",
		"type": "now",
		"requestMethod": "vmware-openshift-upi",
		"region": "%s",
		"maxMemory": "",
		"startCpus": "",
		"postInstallScript": "",
		"approvalGroupId": "",
		"customerData": "false",
		"customerDataTypes": [],
		"iui": "6960009WPF",
		"platform": {
		  "id": "63a3a25a3a4689001740dbb3",
		  "createdAt": 1670281452745,
		  "updatedAt": 1713095738558,
		  "oid": "63a3a1de838d0e5a30d649ce",
		  "name": "OpenShift Cluster (VMware on IBM Cloud) - UPI - Public",
		  "description": "Self-Managed UPI OpenShift cluster (VMware on IBM Cloud) with NFS or ODF (OCS) storage.\n\n**Environment details:**\n\n*   UPI - User Provided Infrastructure\n*   Cloud platform independent (Bare-metal simulation)\n*   This environment can **NOT** be templated.\n*   Use this environment if you need a public ingress.\n\n**Self-service options:**\n\n*   OCP version: 4.8 - 4.12\n*   Cluster size: 3, 4, or 5 worker nodes\n*   Worker node ‘flavor’ (vCPUs x Memory): 4x16, 8x32, 16x64, 32x128\n*   FIPS: yes/no. (FIPS security features enabled on the nodes of the cluster)\n*   Managed NFS 2TB\n*   ODF (local disks) - either no (leave unchosen), or a storage size and options are 500GB, 2TB, 5TB\n*   Ephemeral Storage: up to 300GB per worker node\n*   SSH and RDP bastion access",
		  "automation": "",
		  "automationBucket": "public-solutions",
		  "infrastructure": "ibmcloud-2",
		  "idleRuntimeLimit": 300,
		  "totalRuntimeLimit": 0,
		  "timeoutAction": "shutdown",
		  "autoStart": true,
		  "smartrdp": true,
		  "linkPatterns": [],
		  "regions": [
			{
			  "name": "us-east",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "us-east",
			  "geo": "americas",
			  "datacenter": "wdc04",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 1,
			  "status": "Enabled",
			  "accountPool": "itzvmware",
			  "cloudAccount": "itzvmware",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "eu-gb",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "eu-gb",
			  "geo": "europe",
			  "datacenter": "lon06",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 2,
			  "status": "Enabled",
			  "accountPool": "itzvmware",
			  "cloudAccount": "itzvmware",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "eu-de",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "eu-de",
			  "geo": "europe",
			  "datacenter": "fra04",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 3,
			  "status": "Enabled",
			  "accountPool": "itzvmware",
			  "cloudAccount": "itzvmware",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "jp-tok",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "jp-tok",
			  "geo": "ap",
			  "datacenter": "tok02",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 4,
			  "status": "Enabled",
			  "accountPool": "itzvmware",
			  "cloudAccount": "itzvmware",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "---",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "any",
			  "geo": "any",
			  "datacenter": "any",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 0,
			  "status": "Disabled",
			  "accountPool": "itzvmware-itzna6-sjc03",
			  "cloudAccount": "any",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "---",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "any",
			  "geo": "any",
			  "datacenter": "any",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 0,
			  "status": "Disabled",
			  "accountPool": "itzvmware-itzeu6-mad02",
			  "cloudAccount": "any",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			},
			{
			  "name": "---",
			  "infrastructure": "ibmcloud-2",
			  "template": "vmware-openshift-upi",
			  "region": "any",
			  "geo": "any",
			  "datacenter": "any",
			  "requestMethod": "vmware-openshift-upi",
			  "description": "",
			  "weight": 0,
			  "status": "Disabled",
			  "accountPool": "itzvmware-itzap4-syd01",
			  "cloudAccount": "any",
			  "variables": [
				{
				  "name": "shared_datastore_cluster",
				  "description": "Share datastore cluster",
				  "type": "bool",
				  "default": "",
				  "value": "true"
				}
			  ]
			}
		  ],
		  "audience": [],
		  "visibility": ["IBMers", "Business Partners"],
		  "adminOptions": ["Allow Workshops"],
		  "weight": 8,
		  "status": "Enabled",
		  "links": [],
		  "reservationChain": [],
		  "startCpus": null,
		  "maxMemory": null,
		  "postInstallScript": null,
		  "disableWorkshops": false,
		  "cloudTarget": null,
		  "cloudAccount": null,
		  "emailTemplate": null,
		  "options": null,
		  "url": null,
		  "bom": [],
		  "collection": {
			"id": "5fb3200cec8dd00017c57f20",
			"createdAt": 1605574668668,
			"updatedAt": 1713215521645,
			"oid": null,
			"simple": false,
			"name": "TechZone Certified Base Images",
			"slug": "tech-zone-certified-base-images",
			"synopsis": "These environments represent your best starting point for building new content, showing customers how easy it is to deploy IBM Technology from scratch or test custom configurations.",
			"description": "This collection is owned by the TechZone team and provides a centralized location for all available environments which have been certified for our users. These environments represent your best starting point for building new content, showing customers how easy it is to deploy IBM Technology from scratch, or testing custom configurations.\n\n**TechZone Certified** means that these environments have been tested, secured, and are fully supported by the TechZone team, plus regularly maintained to ensure reliability and stability.\n\n**Getting Started**\n\nThe left navigation outlines the preferred and premium environment offerings. Before you select your environment, read the page descriptions to understand reservation duration and available purpose. You can also find this information in the [Reservation Duration Policies Runbook](https://github.com/IBM/itz-support-public/blob/main/IBM-Technology-Zone/IBM-Technology-Zone-Runbooks/reservation-duration-policy.md).\n\n**Preferred** environments are the lowest cost options with the longest duration. It is strongly recommended that you leverage these environments, such as **Pre-installed Cloud Paks** and **VMware on IBM Cloud** which provide the ability to reserve for up to 49 days for a customer facing opportunity and 8 days for self-education or testing.\n\n**Premium** environments are the most expensive and provide the shortest duration, plus they are not available for self-education or testing. These environments are strictly for customer facing opportunities and must be reserved with a valid opportunity code or Gainsight ID.\n\n**Note:** Additional vetting will be required if the infrastructure costs 10 of the opportunity size or up to $10K maximum.",
			"url": null,
			"sales": null,
			"support": "techzone.help@ibm.com",
			"type": [],
			"cover": "https://dte2.s3.us-east.cloud-object-storage.appdomain.cloud/tech-zone-certified-base-images-certified.jpeg",
			"brands": [
			  "Automation",
			  "Data and AI",
			  "Public Cloud",
			  "Security",
			  "Sustainability Software"
			],
			"categories": [
			  { "name": "Public Cloud Platform", "utcode": "153QH", "utlevel": "15" },
			  { "name": "Data and AI", "utcode": "15ANP", "utlevel": "15" },
			  { "name": "Automation", "utcode": "15IGO", "utlevel": "15" },
			  { "name": "Red Hat", "utcode": "15EDG", "utlevel": "15" }
			],
			"businessUnits": [
			  { "name": "IBM Software w/o TPS", "utcode": "10A00", "utlevel": "10" }
			],
			"products": [
			  { "name": "Redhat" },
			  {
				"name": "6950-33X IBM Cloud infrastructure on Cloud BU paper (Redhat)",
				"description": null,
				"utcode": "30QAB",
				"utlevel": "30",
				"utcodes": {
				  "10": "10A00",
				  "15": "153QH",
				  "17": "17ISM",
				  "20": "20D0C",
				  "30": "30QAB"
				},
				"utdescriptions": {
				  "10": "IBM Software w/o TPS",
				  "15": "Public Cloud Platform",
				  "17": "Public Cloud IaaS",
				  "20": "Public Cloud Infrastructure Services",
				  "30": "6950-33X IBM Cloud infrastructure on Cloud BU paper (Redhat)"
				}
			  },
			  {
				"name": "OpenShift",
				"shortname": "OpenShift",
				"longname": "OpenShift",
				"description": null,
				"utcode": "30A6M",
				"utlevel": "30",
				"utcodes": {
				  "10": "10A00",
				  "15": "15EDG",
				  "17": "17JLM",
				  "20": "20ANW",
				  "30": "30A6M"
				},
				"utdescriptions": {
				  "10": "IBM Software w/o TPS",
				  "15": "Red Hat",
				  "17": "Red Hat Market",
				  "20": "Red Hat Portfolio",
				  "30": "OpenShift"
				}
			  },
			  {
				"name": "watsonx.data",
				"shortname": "watsonx.data",
				"longname": "watsonx.data",
				"description": null,
				"utcode": "30AW0",
				"utlevel": "30",
				"utcodes": {
				  "10": "10A00",
				  "15": "15ANP",
				  "17": "17W1X",
				  "20": "20AQR",
				  "30": "30AW0"
				},
				"utdescriptions": {
				  "10": "IBM Software w/o TPS",
				  "15": "Data and AI",
				  "17": "Fabric Market",
				  "20": "watsonx",
				  "30": "watsonx.data"
				}
			  },
			  {
				"name": "watsonx.ai",
				"shortname": "watsonx.ai",
				"longname": "watsonx.ai",
				"description": null,
				"utcode": "30AX6",
				"utlevel": "30",
				"utcodes": {
				  "10": "10A00",
				  "15": "15ANP",
				  "17": "17W1X",
				  "20": "20AQR",
				  "30": "30AX6"
				},
				"utdescriptions": {
				  "10": "IBM Software w/o TPS",
				  "15": "Data and AI",
				  "17": "Fabric Market",
				  "20": "watsonx",
				  "30": "watsonx.ai"
				}
			  }
			],
			"portfolios": [
			  {
				"name": "Data Integration & Governance Portfolio",
				"utcode": "20A11",
				"utlevel": "20"
			  },
			  {
				"name": "IT Automation Portfolio",
				"utcode": "20APO",
				"utlevel": "20"
			  },
			  {
				"name": "Business Automation Portfolio",
				"utcode": "20A0Y",
				"utlevel": "20"
			  },
			  { "name": "Red Hat Portfolio", "utcode": "20ANW", "utlevel": "20" },
			  { "name": "watsonx", "utcode": "20AQR", "utlevel": "20" }
			],
			"industries": [],
			"tags": [
			  "TechZone",
			  "Certified",
			  "Base Image",
			  "Base Images",
			  "TechZone Certified Base Images",
			  "Certified Base Images",
			  "Certified Images",
			  "TechZone Images",
			  "Openshift cluster",
			  "IPI ",
			  "UPI",
			  "VM",
			  "Linux",
			  "Windows",
			  "Ubuntu",
			  "RHEL",
			  "OCP",
			  "FIPS",
			  "Multizone",
			  "multi zone",
			  "request",
			  "vcenter",
			  "ocp gym",
			  "VPC",
			  "ROKS",
			  "Azure",
			  " AWS",
			  "Fyre",
			  "VMWare",
			  "IBM Cloud",
			  "Cloud",
			  "Power Systems",
			  "Power",
			  "Systems",
			  "Storage",
			  "patterns",
			  "premium",
			  "Classic",
			  "gen2",
			  "watsonx",
			  "watsonx governance,",
			  "watsonx assistant",
			  "watsonx.data",
			  "watsonx.ai",
			  "watsonx.governance",
			  "saas"
			],
			"flags": ["Featured", "Platinum"],
			"audience": ["Developer"],
			"visibility": ["IBMers", "Business Partners"],
			"language": "English",
			"status": "Active",
			"owner": "techzone.help@ibm.com",
			"collaborators": [
			  "ben.foulkes@ca.ibm.com",
			  "brooke.jones@ibm.com",
			  "\tjaericks@us.ibm.com",
			  "david.stacy@ibm.com",
			  "craig.cooper@ibm.com"
			],
			"expireAt": null,
			"verifiedAt": null,
			"publishedAt": null
		  },
		  "approvalGroup": null,
		  "approvalGated": false,
		  "automationGated": false,
		  "approvalGroupId": null
		},
		"terms": true,
		"type-0": "now",
		"dynamicOutputs": [
		  {
			"name": "shared_datastore_cluster",
			"description": "Share datastore cluster",
			"type": "bool",
			"default": "",
			"value": "true"
		  },
		  { "name": "ocs", "value": "5Ti" },
		  { "name": "worker_node_count", "value": "5" },
		  { "name": "worker_node_flavor", "value": "32x128.300gb" }
		],
		"serviceLinks": [],
		"accountPool": "itzvmware",
		"geo": "ap",
		"datacenter": "jp-tok-1",
		"cloudAccount": "itzvmware",
		"end_time": "9:21",
		"ocs": "5Ti",
		"worker_node_count": "5",
		"worker_node_flavor": "32x128.300gb",
		"policy": {
		  "id": "63e3aa31bcbabc0017092c1b",
		  "name": "VMWare-Education-Policy",
		  "origin": "",
		  "geography": [],
		  "infrastructure": ["ibmcloud-2"],
		  "dataCenter": [],
		  "collection": [],
		  "user": [],
		  "asset": [
			"61f848ce227cdf001ec688c1",
			"6239e8adaa08bd001870c845",
			"6241d306a81132001fcfe0d1",
			"627abf3c5f0e75001fb4f2e6",
			"628bc766da758e001e195109",
			"629a50dfe951a7001e9057b8",
			"629a517c337248001edeefa4",
			"62c346830280430017722f43",
			"62c58a024dad4e001896835a",
			"62e92ae09f395200173aaa05",
			"632b1ac79df3b90017c93689",
			"6332070a01fdec0018de2baa",
			"634afebd4de3d9001816ca57",
			"6362cb2ec0d65c00185a7d09",
			"6362d110faa5690019153c9b",
			"6363e405e9e3770019fbec25",
			"636986b6f115890019662ff9",
			"636a03a12ce6940018d602a0",
			"63725bef4f612c0019e6cf45",
			"6372812b4f612c0019e6cf51",
			"637283a1c9cdd500190e643f",
			"6372e278c9cdd500190e649c",
			"637399157c0bee00180d36c9",
			"6373eda26fa5fe0018a83cb0",
			"6374053851f5be00198a2f82",
			"637411ed4ce4080018907ddd",
			"6374d5ff51f5be00198a2f9d",
			"6374e046f680be001836a364",
			"6374e8c36fa5fe0018a83ce0",
			"637556d1f680be001836a388",
			"637636b8ffd03d0019b8c0f0",
			"637bcaa5cd24bc00190d29a2",
			"637cdc066b7a060019ede0a2",
			"637e310ecd24bc00190d29f0",
			"637f9a15cd24bc00190d2a32",
			"63851c239508a40019c9203c",
			"6385f8c337f8a600183c7317",
			"638626bf9508a40019c9206b",
			"63867b6a4cd2a3001961ea2c",
			"63876509eb61a800183d69c7",
			"638768e537f8a600183c7378",
			"63876b18eb61a800183d69ca",
			"63876cd4bbddcf0019599427",
			"63876e3d00770b0019d9f47d",
			"63877af037f8a600183c737b",
			"6387aa80eb61a800183d69d3",
			"6388238e4a17fd00198f1e62",
			"6388b0114a17fd00198f1e85",
			"6388e372f381050018ba354f",
			"6389feba71c87d001831c119",
			"638a047471c87d001831c11e",
			"638a1f6a308f5500188808a7",
			"638e319647db6b00195a08b5",
			"638ed3ba5cad290018089476",
			"638f5f8467b96f00194cbf3e",
			"63909560aec20e0018fd5b9c",
			"6390fa02c7f5ac001abadc9c",
			"63920cf22323f100193228c4",
			"6392347d7206a7001886da80",
			"6392386ff80e0100183fd7d9",
			"63923e1ccaa30c0019f8910d",
			"6392411910e072001846ea55",
			"63924429f80e0100183fd7df",
			"6392f6d20f8a690018e13a1b",
			"639321f64c9c1700196e7c7a",
			"6393be27d282e700186def99",
			"639795546535e60017f29baf",
			"6397a8b86535e60017f29bb2",
			"63989b063ea77b00185598e7",
			"6398e2e09df2c80018d4adf3",
			"639a4de66535e60017f29cc1",
			"639cad5454b4ce0018719ec9",
			"639cbd4c85bbf40018d8fd8d",
			"63a0ba390644e60019079a7a",
			"63a12bbb04473c001a352155",
			"63a2466053636400183ace2b",
			"63a373783a4689001740dbad",
			"63a3761704eb640017c9b067",
			"63a3765f5f0cda0017d3dc05",
			"63a47c4de556e00017ef3b8a",
			"63a47dac04eb640017c9b086",
			"63a47ee57772160018803e70",
			"63adcfd9b496f00017e85efe",
			"63af62f9284b1600170d816b",
			"63af66e4ded18e001966b0bc",
			"63af68f2ded18e001966b0bf",
			"63af6a5a9849610018a7788a",
			"63af6f2c284b1600170d8171",
			"63af703bded18e001966b0c3",
			"63af745dded18e001966b0c5",
			"63af7541ded18e001966b0c8",
			"63b16c99284b1600170d817b",
			"63b2ae8a284b1600170d8184",
			"63b2bd1fb3c21c0018330ac5",
			"63b2d6b0284b1600170d8189",
			"63b4b4f83e0b0e00173bd998",
			"63b6f8225e31ee00177054a2",
			"63bc2756e3a24a001780018d",
			"63bc7ea6ec1358001858eda7",
			"63beddd3e5f3800018a2b32e",
			"63bee2453e31020017c6c6ab",
			"63bee758a4b8dc001741eab7",
			"63beed95a4b8dc001741eab9",
			"63c0040afb0c8400188e5bf1",
			"63c03fc7610ba600185c8da8",
			"63c18ff06563460018361910",
			"63c1ba7ad4f6a700174de92f",
			"63c585c618d23d00182aa34c",
			"63c8294bd5bd7f001708196f",
			"63c981b085aa110017839ee0",
			"63cddc2e63026a0017d5f99c",
			"63d13db0af86280017ce26b6",
			"63d41f9832ab5e00181f701c",
			"63da96861bb6640018053793",
			"63dae88f522b7e00187a897f",
			"63db9ea3a3b642001821f680",
			"63dba0237896b400170f308b",
			"63dba0d0cc19150018af0848",
			"63dba171cdd28a00178ba0ca",
			"63dba359cc19150018af084f",
			"63dc05298b134800173dcce4",
			"63dd14f6c7920900171f0950",
			"5fd383938d8480001f437a22",
			"60093b01cf63a20018dbc904",
			"63db9f61fd02a70017a2e64d",
			"63d957bbbb11c2001704e65c",
			"636c1075ed3d04001837f3de",
			"63b735d3186dad00181dc351",
			"63b6c0c3e05d350017a2d114",
			"638e535667b96f00194cbf1d",
			"63dd88f090a8f20018fc27e5",
			"63a6024eb4fedf001808c07c",
			"6377e0887b9158001886b30f",
			"63d04aa64e26be0017d9faaa",
			"6393a6cf1a0c1f0018adb6c3",
			"63a32cbdb766ac0018dfd709",
			"627506721721b0001fb3ab31",
			"633b6d7f71efc80018e69d66",
			"63d4558a3cb1f20018b39836",
			"63bebdd0e5f3800018a2b318",
			"63e2c4a7e2a7bf0018430cc8",
			"63c15d75e8a0e900174f1c71",
			"63d24b58d95bd8001902e9c1",
			"63da62b0bfa7190018c20105",
			"63cefa8a50d70900186db96c",
			"63e36757e2a7bf0018430ccd",
			"63fcf346113069a1fb42ac54",
			"63ed14ef7c001900174b6032",
			"63a3a25a3a4689001740dbb3",
			"6400e60a502aba0018ef8929",
			"63ebedd694411800184aaced",
			"609561071dfddb001f59e773",
			"64071f935d66300018b64414",
			"6413928d8654ff0017946bd6",
			"63f44a13c1bafa00188ac6e0",
			"64189742a023cd00194fa703",
			"64023309970a4f0017501680",
			"63ee798afb861100199f8f22",
			"6418cd173155c30017d36ba1",
			"63a0cf6c61dc1d0018b71c67",
			"640b14f790c20500182c0af8",
			"6419c672c5ba610017a4441b",
			"63f6f13b4ea4020017706f9e",
			"64188d021177a30019fbc038",
			"6414cfe4f1008a0018850d9c",
			"6414bbf8dc2d750017c9a0e7",
			"6408a645bc7e3f001721cb93",
			"641ae96d450cf9001771e10e",
			"631b6a3b845753001803d600",
			"63be1385b08286001889511c",
			"63e3db7838f9f600196eab23",
			"63e4d1ae38f9f600196eab39",
			"63e51a4a38f9f600196eab44",
			"63e51c9660a2f70019770e74",
			"63e56ee32580420017be3132",
			"63e65b2b1aeb6d00182ac305",
			"63e6afda27b6b9001743fab0",
			"63e767122580420017be315e",
			"63ed44687c001900174b6042",
			"63ef936b65678700186e066f",
			"63f8e94bda35890018a97bc4",
			"63f932285d7bcc00188c04c1",
			"63fd41525dd0f9001958cc83",
			"63fe944f0757a7001788309f",
			"63fefbd0c655700018c5dc1b",
			"63ff78cafcef5c0017f59f7d",
			"63ffa94ab6b41a00178377f5",
			"63ffed9657eefa0017b10080",
			"6400b16c3a06140018785b6f",
			"64041241eddc870018edadcb",
			"640429b1eddc870018edadce",
			"6406d15676bf7100177f7356",
			"640733e55d66300018b64418",
			"64076beaf8b18000177141d0",
			"6408f30ff92d410019b7ac09",
			"640b6a53d6521e001a015a72",
			"640b6dd30e25720017542e81",
			"640bc34a0e25720017542e91",
			"6412129e76e78200188ffd3e",
			"6413246dea3dc8001816a79f",
			"641374a05c146400187997a3",
			"6413855d5c146400187997b0",
			"6413855d5c146400187997b1",
			"641388e284742a0017a2ee09",
			"64145b914514070018b1cefa",
			"64146bcd64991900187a021b",
			"64146e235f1b500018322e75",
			"641d9ec1ad1b27001822140f",
			"641c3e85218f0e0017ad610e",
			"6441aa8a198c4332a346e408",
			"6441aa2d198c4332a346e406",
			"6441aa64198c4332a346e407",
			"6441aaad198c4332a346e409",
			"6441aae4198c4332a346e40a",
			"64acd6435d1077001755cbc2",
			"65366cbbc0d4aa0017e23fb8"
		  ],
		  "purposes": ["Practice / Self-Education"],
		  "roles": [],
		  "owner": "brooke.jones@ibm.com",
		  "quota": {},
		  "inherit": false,
		  "status": "Enabled",
		  "weight": 0,
		  "pattern": [],
		  "groups": [],
		  "policies": {
			"defaultLength": 172800,
			"offsetLength": 0,
			"maximumLength": 172800,
			"extensionLength": 172800,
			"extensionLimit": 2,
			"opportunity": false
		  },
		  "updatedAt": 1711575025425,
		  "createdAt": 1675864625835,
		  "type": "Duration",
		  "defaultLength": 172800,
		  "offsetLength": 0,
		  "maximumLength": 172800,
		  "extensionLength": 172800,
		  "extensionLimit": 2,
		  "opportunity": false,
		  "score": 100010100,
		  "start": "2024-04-15T23:21:00.000Z",
		  "end": "2024-04-17T23:20:00.000Z",
		  "extend": "2024-04-17T23:20:00.000Z",
		  "startMinDate": "2024-04-15T23:21:00.000Z",
		  "startMaxDate": "2025-04-10T23:21:00.000Z",
		  "endMinDate": "2024-04-15T23:20:00.000Z",
		  "endDefaultDate": "2024-04-17T23:20:00.000Z",
		  "endMaxDate": "2024-04-17T23:20:00.000Z",
		  "reserveMin": 172800,
		  "reserveMax": 172800,
		  "extensionMaxDate": "2024-04-19T23:20:00.000Z",
		  "isExtendable": true,
		  "validation": {
			"start": true,
			"end": true,
			"extension": null,
			"opportunity": null,
			"opportunityProduct": null
		  },
		  "inPolicy": true
		}
	  }
	  `, purpose, description, start, end, email, region)

		logger.Debug("Successfully generated JSON request...")
		logger.Debug(JSON_Request)

    return JSON_Request
}
