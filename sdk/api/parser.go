package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// ResponseStatus validate the response and parse the response body
func ValidateHttpResponse(err error, response *http.Response, responseType any) (*ResponseStatus, error) {
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, errors.New("response is nil")
	}
	// parse body into LoginResponse
	if response.Body == nil {
		return nil, errors.New("response body is nil")
	}
	return parseResponse(response, responseType)
}

func ParseResponseAsString(response *http.Response) (*ResponseStatus, string, error) {
	if response == nil {
		return nil, "", errors.New("response is nil")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errResponse ResponseError
		if response.Body != nil {
			err := json.NewDecoder(response.Body).Decode(&errResponse)
			if err != nil {
				return nil, "", err
			}
			return &ResponseStatus{
				ErrorCode:  errResponse.Code,
				HttpStatus: response.StatusCode,
			}, "", nil
		}

		return &ResponseStatus{
			ErrorCode:  response.StatusCode,
			HttpStatus: response.StatusCode,
		}, "", nil
	}
	var bodyString string
	if response.Body != nil {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, "", err
		}
		bodyString = string(bodyBytes)
	}
	return &ResponseStatus{
		ErrorCode:  0,
		HttpStatus: response.StatusCode,
	}, bodyString, nil
}

func parseResponse(response *http.Response, responseType any) (*ResponseStatus, error) {
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
