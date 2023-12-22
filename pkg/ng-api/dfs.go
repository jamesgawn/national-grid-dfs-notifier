package ngapi

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type DFSRequirementResultRecords struct {
	AcceptancePriceGBPPerMWh string `json:"Guaranteed Acceptance Price GBP per MWh"`
	From                     string `json:"From"`
	DeliveryDate             string `json:"Delivery Date"`
	DespatchType             string `json:"Despatch Type"`
	DFSRequestedMW           int    `json:"DFS Required MW"`
	To                       string `json:"To"`
	ServiceRequirementType   string `json:"Service Requirement Type"`
	Id                       int    `json:"_id"`
	ParticipaentBidsEligible string `json:"Participant Bids Eligible"`
}

type DFSRequirementResult struct {
	Records []DFSRequirementResultRecords `json:"records"`
	Sql     string                        `json:"sql"`
}

type DFSRequirementResponse struct {
	Help    string               `json:"help"`
	Success bool                 `json:"success"`
	Result  DFSRequirementResult `json:"result"`
}

func GetDemandFlexibilityServiceRequirements() []DFSRequirementResultRecords {
	response, err := http.Get("https://api.nationalgrideso.com/api/3/action/datastore_search_sql?sql=SELECT%20*%20FROM%20%20%227914dd99-fe1c-41ba-9989-5784531c58bb%22%20ORDER%20BY%20%22_id%22%20ASC%20LIMIT%20100")

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve DFS requirement from National Grid API")
	}

	var decodedResponse DFSRequirementResponse
	err = json.NewDecoder(response.Body).Decode(&decodedResponse)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode DFS requirement data from National Grid API")
	}

	log.Debug().Int("statusCode", response.StatusCode).Msg("Retrieved DFS data from National Grid API")

	return decodedResponse.Result.Records
}
