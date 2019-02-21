package zendesk

import (
	"encoding/json"
	"fmt"
	"time"

	resty "gopkg.in/resty.v0"
)

// Client struct defines common api access
type Client struct {
	domain     string
	client     *resty.Client
	apiVersion string
	reqi       *func(*resty.Request)
	resi       *func(*resty.Response)
}

// User func gets Users API
func (c *Client) User() UserAPIHandler {
	return UserAPIHandler{c}
}

// Ticket func gets Tickets API
func (c *Client) Ticket() TicketAPIHandler {
	return TicketAPIHandler{c}
}

// InterceptRequestResponse calls interceptors
func (c *Client) InterceptRequestResponse(request *func(*resty.Request), response *func(*resty.Response)) {
	if request == nil {
		r := func(req *resty.Request) {
			fmt.Printf("Request size: %d\n", req.RawRequest.ContentLength)
		}

		request = &r
	}

	if response == nil {
		r := func(res *resty.Response) {
			fmt.Printf("Response size: %d\n", res.RawResponse.ContentLength)
		}

		response = &r
	}

	c.reqi = request
	c.resi = response
}

// toFullURL parses path to endpoint URL
func (c *Client) toFullURL(path string) string {
	return fmt.Sprintf("https://%v.zendesk.com/api/%s/%s", c.domain, c.apiVersion, path)
}

// get handles GET endpoint
func (c *Client) get(path string, params map[string]string) (*resty.Response, error) {
	resp, err := c.client.R().SetQueryParams(params).Get(c.toFullURL(path))

	if c.reqi != nil {
		req := *c.reqi
		req(resp.Request)
	}

	if c.resi != nil {
		res := *c.resi
		res(resp)
	}

	if resp != nil && resp.StatusCode() == 429 {
		duration, err := time.ParseDuration(fmt.Sprintf("%ss", resp.Header().Get("Retry-After")))

		if err == nil {
			fmt.Println(fmt.Sprintf("Sleeping %s", duration.String()))
			time.Sleep(duration)
			return c.get(path, params)
		}
	}

	if err != nil {
		return resp, err
	}

	return resp, nil
}

// post handles POST endpoint
func (c *Client) post(path string, params interface{}, v interface{}) (interface{}, error) {
	resp, err := c.client.R().SetBody(params).SetResult(&v).Post(c.toFullURL(path))

	if c.reqi != nil {
		req := *c.reqi
		req(resp.Request)
	}

	if c.resi != nil {
		res := *c.resi
		res(resp)
	}

	if resp != nil && resp.StatusCode() == 429 {
		duration, err := time.ParseDuration(fmt.Sprintf("%sm", resp.Header().Get("Retry-After")))

		if err == nil {
			fmt.Println(fmt.Sprintf("Sleeping %s", duration.String()))
			time.Sleep(duration)
			return c.post(path, params, v)
		}
	}

	if err != nil {
		return v, err
	}

	return v, nil
}

// put handles PUT endpoint
func (c *Client) put(path string, params interface{}, v interface{}) (interface{}, error) {
	resp, err := c.client.R().SetBody(params).SetResult(&v).Put(c.toFullURL(path))

	if c.reqi != nil {
		req := *c.reqi
		req(resp.Request)
	}

	if c.resi != nil {
		res := *c.resi
		res(resp)
	}

	if resp != nil && resp.StatusCode() == 429 {
		duration, err := time.ParseDuration(fmt.Sprintf("%sm", resp.Header().Get("Retry-After")))

		if err == nil {
			fmt.Println(fmt.Sprintf("Sleeping %s", duration.String()))
			time.Sleep(duration)
			return c.put(path, params, v)
		}
	}

	if err != nil {
		return v, err
	}

	return v, nil
}

// delete handles DELETE endpoint
func (c *Client) delete(path string) (*resty.Response, error) {
	resp, err := c.client.R().Delete(c.toFullURL(path))

	if c.reqi != nil {
		req := *c.reqi
		req(resp.Request)
	}

	if c.resi != nil {
		res := *c.resi
		res(resp)
	}

	if resp != nil && resp.StatusCode() == 429 {
		duration, err := time.ParseDuration(fmt.Sprintf("%sm", resp.Header().Get("Retry-After")))

		if err == nil {
			fmt.Println(fmt.Sprintf("Sleeping %s", duration.String()))
			time.Sleep(duration)
			return c.delete(path)
		}
	}

	return resp, err
}

// parseResponseToInterface func convert json response into an interface
func (c *Client) parseResponseToInterface(response *resty.Response, v interface{}) {
	json.Unmarshal(response.Body(), &v)
}
