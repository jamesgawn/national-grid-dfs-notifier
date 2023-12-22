package natgridapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type DFSRequirementResultRecordsDto struct {
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

type DFSRequirementResultDto struct {
	Records []DFSRequirementResultRecordsDto `json:"records"`
	Sql     string                           `json:"sql"`
}

type DFSRequirementResponseDto struct {
	Help    string                  `json:"help"`
	Success bool                    `json:"success"`
	Result  DFSRequirementResultDto `json:"result"`
}

func GetDemandFlexibilityServiceRequirements() []DFSRequrement {
	response, err := http.Get("https://api.nationalgrideso.com/api/3/action/datastore_search_sql?sql=SELECT%20*%20FROM%20%20%227914dd99-fe1c-41ba-9989-5784531c58bb%22%20ORDER%20BY%20%22_id%22%20ASC%20LIMIT%20100")

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve DFS requirement from National Grid API")
	}

	var decodedResponse DFSRequirementResponseDto
	err = json.NewDecoder(response.Body).Decode(&decodedResponse)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode DFS requirement data from National Grid API")
	}

	log.Debug().Int("statusCode", response.StatusCode).Msg("Retrieved DFS data from National Grid API")

	return adaptDFSRRDtoToDomain(decodedResponse.Result.Records)
}

type DFSRequrement struct {
	Id                       int
	AcceptancePriceGBPPerMWh int
	RequestedMW              int
	DispatchType             string
	ServiceRequirementType   string
	EligibleSuppliers        []string
	From                     time.Time
	To                       time.Time
}

// Convert raw reuirement dto to a domain object more adept to be utilised in business logic
func adaptDFSRRDtoToDomain(requirementDtos []DFSRequirementResultRecordsDto) []DFSRequrement {
	requirements := make([]DFSRequrement, 0)

	for _, requirementDto := range requirementDtos {
		id := requirementDto.Id
		eligibleSuppliers := strings.Split(requirementDto.ParticipaentBidsEligible, ",")

		acceptancePrice, err := strconv.Atoi(requirementDto.AcceptancePriceGBPPerMWh)
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to convert acceptance price in requirementDto: %d", id)
		}

		from, err := time.Parse(time.DateTime, requirementDto.DeliveryDate+" "+requirementDto.From)
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to convert from time in requirementDto: %d", id)
		}

		to, err := time.Parse(time.DateTime, requirementDto.DeliveryDate+" "+requirementDto.To)
		if err != nil {
			log.Fatal().Err(err).Msgf("Unable to convert to time in requirementDto: %d", id)
		}

		requirement := DFSRequrement{
			Id:                       id,
			AcceptancePriceGBPPerMWh: acceptancePrice,
			RequestedMW:              requirementDto.DFSRequestedMW,
			DispatchType:             requirementDto.DespatchType,
			ServiceRequirementType:   requirementDto.ServiceRequirementType,
			EligibleSuppliers:        eligibleSuppliers,
			From:                     from,
			To:                       to,
		}
		requirements = append(requirements, requirement)
	}
	return requirements
}
