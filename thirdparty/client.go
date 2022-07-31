package thirdparty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/thoniwutr/schedule-school-teachning-bsd13-backend/util"
)

// ApiClient is a struct that calls internal APIs
type ApiClient struct {
	log             *util.Logger
	client          *http.Client
	urlOrganisation string
	urlRecipient    string
}

// NewApiClient returns the ApiClient struct
func NewApiClient(log *util.Logger, urlOrganisation string, urlRecipient string) *ApiClient {
	return &ApiClient{log: log, client: &http.Client{}, urlOrganisation: urlOrganisation, urlRecipient: urlRecipient}
}

func (ac ApiClient) CreateOrganisation(req *OrganisationRequest, userInfo string) (*OrganisationResponse, error) {

	resp, err := ac.httpRequest(http.MethodPost, ac.urlOrganisation, userInfo, req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if isHTTPSuccess(resp.StatusCode) {

		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		organisationRes := &OrganisationResponse{}
		if err := json.NewDecoder(ioutil.NopCloser(bytes.NewBuffer(buf))).Decode(organisationRes); err != nil {
			return nil, err
		}

		return organisationRes, nil
	}

	return nil, fmt.Errorf("backend: failed to create organisation with status code: %v", resp.StatusCode)

}

func (ac ApiClient) FollowEntity(source string, organisationId string) error {

	resp, err := ac.httpRequest(http.MethodPost, fmt.Sprintf("%v/%v/follow/%v", ac.urlRecipient, source, organisationId), "", nil)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if isHTTPSuccess(resp.StatusCode) {
		return nil
	}

	return fmt.Errorf("backend: failed to follow entity with status code: %v", resp.StatusCode)
}

func (ac ApiClient) CreateRecipient(crReq *CreateRecipientRequest) error {

	resp, err := ac.httpRequest(http.MethodPost, ac.urlRecipient, "", crReq)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if isHTTPSuccess(resp.StatusCode) {
		return nil
	}

	return fmt.Errorf("backend: failed to create recipient with status code: %v", resp.StatusCode)
}

func (ac *ApiClient) httpRequest(httpMethod, apiPath string, userInfo string, reqStruct interface{}) (*http.Response, error) {

	jsonReq, err := json.Marshal(reqStruct)
	if err != nil {
		ac.log.Error().Err(err).Msg("failed to marshal request to json")
		return nil, err
	}

	req, err := http.NewRequest(httpMethod, apiPath, bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	if len(userInfo) > 0 {
		req.Header.Set("X-Endpoint-API-UserInfo", userInfo)
	}

	if err != nil {
		return nil, err
	}

	resp, err := ac.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func isHTTPSuccess(code int) bool {
	return code >= 200 && code < 300
}
