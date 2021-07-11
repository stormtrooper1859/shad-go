package firewall

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var _ http.RoundTripper = (*firewall)(nil)

type firewall struct {
	config *RulesConfig
}

func NewFirewall(config *RulesConfig) http.RoundTripper {
	f := firewall{
		config: config,
	}
	return &f
}

func (f *firewall) RoundTrip(request *http.Request) (*http.Response, error) {
	rule := f.getRule(request.RequestURI)

	allowed := rule.allowedRequest(request)
	if !allowed {
		return getBlockedQuery(), nil
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		return nil, err
	}

	allowed = rule.allowedResponse(response)
	if !allowed {
		return getBlockedQuery(), nil
	}

	return response, nil
}

func (f *firewall) getRule(requestURI string) *Rule {
	index := 0
	for ; index < len(f.config.Rules); index++ {
		if f.config.Rules[index].Endpoint == requestURI {
			break
		}
	}

	if index == len(f.config.Rules) {
		return nil
	}

	return &f.config.Rules[index]
}

func (r *Rule) allowedRequest(request *http.Request) bool {
	if r == nil {
		return true
	}

	// check UserAgent
	UserAgent := request.UserAgent()
	for _, v := range r.ForbiddenUserAgents {
		matchString, err := regexp.MatchString(v, UserAgent)
		if matchString || err != nil {
			return false
		}
	}

	// check RequiredHeaders
	for _, v := range r.RequiredHeaders {
		if request.Header.Get(v) == "" {
			return false
		}
	}

	// check ForbiddenHeaders
	for _, v := range r.ForbiddenHeaders {
		fhPair := strings.SplitN(v, ": ", 2)
		fieldName, fieldRegex := fhPair[0], fhPair[1]
		matchString, err := regexp.MatchString(fieldRegex, request.Header.Get(fieldName))
		if matchString || err != nil {
			return false
		}
	}

	// get body
	if request.Body != nil {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return false
		}

		request.Body = ioutil.NopCloser(bytes.NewReader(body))

		// check MaxRequestLengthBytes
		if r.MaxRequestLengthBytes != 0 && len(body) > r.MaxRequestLengthBytes {
			return false
		}

		// check ForbiddenRequestRe
		for _, v := range r.ForbiddenRequestRe {
			matchString, err := regexp.MatchString(v, string(body))
			if matchString || err != nil {
				return false
			}
		}
	}

	return true
}

func (r *Rule) allowedResponse(response *http.Response) bool {
	if r == nil {
		return true
	}

	// check ForbiddenResponseCodes
	for _, v := range r.ForbiddenResponseCodes {
		if response.StatusCode == v {
			return false
		}
	}

	// get body
	if response.Body != nil {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return false
		}

		response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		// check MaxRequestLengthBytes
		if r.MaxResponseLengthBytes != 0 && len(body) > r.MaxResponseLengthBytes {
			return false
		}

		// check ForbiddenRequestRe
		for _, v := range r.ForbiddenResponseRe {
			matchString, err := regexp.MatchString(v, string(body))
			if matchString || err != nil {
				return false
			}
		}
	}

	return true
}

func getBlockedQuery() *http.Response {
	return &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("Forbidden")),
		StatusCode: 403,
	}
}