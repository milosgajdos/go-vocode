package vocode

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/milosgajdos/go-vocode/request"
)

type CallStatus string

const (
	CallNotStarted CallStatus = "not_started"
	CallInProgress CallStatus = "in_progress"
	CallError      CallStatus = "error"
	CallEnded      CallStatus = "ended"
)

type CallStage string

const (
	CallCreated         CallStage = "created"
	CallPickedUp        CallStage = "picked_up"
	CallTransferStart   CallStage = "transfer_started"
	CallTransferSuccess CallStage = "transfer_successful"
)

type CallStageOutcome string

const (
	CallStageHumanUnAnswer      CallStageOutcome = "human_unanswered"
	CallStageHumanDisconnect    CallStageOutcome = "human_disconnected"
	CallStageDidNotConnect      CallStageOutcome = "call_did_not_connect"
	CallStageBotDisconnect      CallStageOutcome = "bot_disconnected"
	CallStageTransferUnAnswer   CallStageOutcome = "transfer_unanswered"
	CallStageTransferDisconnect CallStageOutcome = "transfer_disconnected"
)

type CallHumanDetection string

const (
	CallHumanDetected   CallHumanDetection = "human"
	CallNoHumanDetected CallHumanDetection = "no_human"
)

type CallOnNoHumanAnswer string

const (
	ContinueCallOnNoHumanAnswer CallOnNoHumanAnswer = "continue"
	HangupCallOnNoHumanAnswer   CallOnNoHumanAnswer = "hangup"
)

type Calls struct {
	Items []Call `json:"items"`
	*Paging
}

type Call struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	Status          CallStatus          `json:"status"`
	ErrorMsg        string              `json:"error_message"`
	RecordAvailable bool                `json:"recording_available"`
	Transcript      string              `json:"transcript"`
	HumanDetection  CallHumanDetection  `json:"human_detection_result"`
	DNC             bool                `json:"do_not_call_result"`
	TelID           string              `json:"telephony_id"`
	Stage           CallStage           `json:"stage"`
	StageOutcome    CallStageOutcome    `json:"stage_outcome"`
	TelMetadata     *TelMetadata        `json:"telephony_metadata"`
	FromNumber      string              `json:"from_number"`
	ToNumber        string              `json:"to_number"`
	Agent           *Agent              `json:"agent"`
	TelProvider     TelProvider         `json:"telephony_provider"`
	AgentPhoneNr    string              `json:"agent_phone_number"`
	StartTime       string              `json:"start_time"`
	EndTime         string              `json:"end_time"`
	HIPAACompliant  bool                `json:"hipaa_compliant"`
	OnNoHumanAnswer CallOnNoHumanAnswer `json:"on_no_human_answer"`
	Context         map[string]any      `json:"context,omitempty"`
	RunDNC          bool                `json:"run_do_not_call_detection"`
	TelAccountConn  *TelAccountConn     `json:"telephony_account_connection"`
	TelParams       map[string]any      `json:"telephony_params"`
}

type CreateCallReq struct {
	FromNr          string              `json:"from_number"`
	ToNr            string              `json:"to_number"`
	Agent           string              `json:"agent"`
	OnHumanNoAnswer CallOnNoHumanAnswer `json:"on_no_human_answer"`
	RunDNC          bool                `json:"run_do_not_call_detection"`
	HIPAACompliant  bool                `json:"hipaa_compliant"`
	Context         map[string]any      `json:"context,omitempty"`
}

type EndCallReq struct {
	ID string `json:"id"`
}

func (c *Client) ListCalls(ctx context.Context, paging *PageParams) (*Calls, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/calls/list")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}
	if paging != nil {
		request.WithPageParams(paging.Encode())
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		actions := new(Calls)
		if err := json.NewDecoder(resp.Body).Decode(actions); err != nil {
			return nil, err
		}
		return actions, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) GetCall(ctx context.Context, callID string) (*Call, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/calls")
	if err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", callID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Call)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) CreateCall(ctx context.Context, createReq *CreateCallReq) (*Call, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/calls/create")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(createReq); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return nil, err
	}

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Call)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) EndCall(ctx context.Context, callID string) (*Call, error) {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/calls/end")
	if err != nil {
		return nil, err
	}

	var body = &bytes.Buffer{}
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)
	// NOTE: this is empty body
	if err := enc.Encode(nil); err != nil {
		return nil, err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodPost, u.String(), body, options...)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("id", callID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		action := new(Call)
		if err := json.NewDecoder(resp.Body).Decode(action); err != nil {
			return nil, err
		}
		return action, nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return nil, errors.Join(err, jsonErr)
		}
		return nil, apiErr
	case http.StatusTooManyRequests:
		return nil, ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return nil, ErrUnprocessableEntity
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}

func (c *Client) GetCallRecording(ctx context.Context, callID string, w io.Writer) error {
	u, err := url.Parse(c.opts.BaseURL + "/" + c.opts.Version + "/calls/recording")
	if err != nil {
		return err
	}

	options := []request.HTTPOption{
		request.WithBearer(c.opts.APIKey),
	}

	req, err := request.NewHTTP(ctx, http.MethodGet, u.String(), nil, options...)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("id", callID)
	req.URL.RawQuery = q.Encode()

	resp, err := request.Do[APIParamError](c.opts.HTTPClient, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if _, err := io.Copy(w, resp.Body); err != nil {
			return err
		}
		return nil
	case http.StatusForbidden, http.StatusBadRequest:
		var apiErr APIGenError
		if jsonErr := json.NewDecoder(resp.Body).Decode(&apiErr); jsonErr != nil {
			return errors.Join(err, jsonErr)
		}
		return apiErr
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusUnprocessableEntity:
		return ErrUnprocessableEntity
	default:
		return fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}
}
