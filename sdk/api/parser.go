package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

func ParseResponse(response *http.Response, responseType any) (*ResponseStatus, error) {
	if response == nil {
		return nil, errors.New("response is nil")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errResponse ResponseError
		if response.Body != nil {
			err := json.NewDecoder(response.Body).Decode(&errResponse)
			if err != nil {
				return nil, err
			}
			return &ResponseStatus{
				ErrorCode:  errResponse.Code,
				HttpStatus: response.StatusCode,
			}, nil
		}

		return &ResponseStatus{
			ErrorCode:  response.StatusCode,
			HttpStatus: response.StatusCode,
		}, nil
	}
	if responseType == nil {
		return nil, errors.New("response type is nil")
	}
	err := json.NewDecoder(response.Body).Decode(responseType)
	if err != nil {
		return nil, err
	}
	return &ResponseStatus{
		ErrorCode:  0,
		HttpStatus: response.StatusCode,
	}, nil
}
