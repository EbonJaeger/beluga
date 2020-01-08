package v1

import (
	"github.com/valyala/fasthttp"
	"net/http"
)

// Stop kills the running Beluga service
func (c *Client) Stop() error {
	// Create the request
	req, err := http.NewRequest("DELETE", formURI("api/v1/daemon"), nil)
	if err != nil {
		return err
	}

	// Send the requeset
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response
	return readError(resp.Body)
}

// StopDaemon handles stopping the daemon
func (l *Listener) StopDaemon(ctx *fasthttp.RequestCtx) {
	writeErrorString(ctx, "Stopping the daemon is not yet implmented", http.StatusBadRequest)
	return
}
