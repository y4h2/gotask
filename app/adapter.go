package app

import (
	"errors"
	"fmt"
	"net/http"
)

type InternalAdapter struct {
}

func (ad *InternalAdapter) SendCancelRequest(host string, taskID string) error {
	req, err := http.NewRequest("", fmt.Sprintf("%s/task/%s", host, taskID), nil)
	if err != nil {
		return err
	}
	// @todo add header to mark the internal call
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("invalid call with response: %d", resp.StatusCode))
	}
	return nil
}
