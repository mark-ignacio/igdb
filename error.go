package igdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// errEndOfJSON occurs when encountering an unexpected end of JSON input.
var errEndOfJSON = errors.New("unexpected end of JSON input")

// Errors returned when the IGDB responds with a problematic status code.
//
// For more information, visit: https://igdb.github.io/api/references/response-codes/
var (
	ErrAuthFailed       = errors.New("IGDB: authentication failed - need valid API key in user-key header")
	ErrBadRequest       = errors.New("IGDB: bad request - check query parameters")
	ErrInternalError    = errors.New("IGDB: internal error - report bug")
	ErrUnexpectedStatus = errors.New("IGDB: unexpected status code")
)

// Error contains information on an
// error returned from an IGDB API call.
type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// checkResponse checks the provided HTTP response
// for errors returned by the IGDB.
func checkResponse(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return ErrBadRequest
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrAuthFailed
	case http.StatusInternalServerError:
		return ErrInternalError
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var e Error

	err = json.Unmarshal(b, &e)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Status %d", e.Status)
	if e.Message != "" {
		msg += fmt.Sprintf(" - %v", e.Message)
	}
	return errors.New(msg)
}
