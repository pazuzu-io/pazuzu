package features

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetAPIFeaturesParams creates a new GetAPIFeaturesParams object
// with the default values initialized.
func NewGetAPIFeaturesParams() *GetAPIFeaturesParams {
	var ()
	return &GetAPIFeaturesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetAPIFeaturesParamsWithTimeout creates a new GetAPIFeaturesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetAPIFeaturesParamsWithTimeout(timeout time.Duration) *GetAPIFeaturesParams {
	var ()
	return &GetAPIFeaturesParams{

		timeout: timeout,
	}
}

// NewGetAPIFeaturesParamsWithContext creates a new GetAPIFeaturesParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetAPIFeaturesParamsWithContext(ctx context.Context) *GetAPIFeaturesParams {
	var ()
	return &GetAPIFeaturesParams{

		Context: ctx,
	}
}

/*GetAPIFeaturesParams contains all the parameters to send to the API endpoint
for the get API features operation typically these are written to a http.Request
*/
type GetAPIFeaturesParams struct {

	/*Limit
	  maximum number of features to return.

	*/
	Limit *int64
	/*Names
	  value, that must present in feature name.

	*/
	Names []string
	/*Offset
	  the offset to start from.

	*/
	Offset *int64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get API features params
func (o *GetAPIFeaturesParams) WithTimeout(timeout time.Duration) *GetAPIFeaturesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get API features params
func (o *GetAPIFeaturesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get API features params
func (o *GetAPIFeaturesParams) WithContext(ctx context.Context) *GetAPIFeaturesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get API features params
func (o *GetAPIFeaturesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithLimit adds the limit to the get API features params
func (o *GetAPIFeaturesParams) WithLimit(limit *int64) *GetAPIFeaturesParams {
	o.SetLimit(limit)
	return o
}

// SetLimit adds the limit to the get API features params
func (o *GetAPIFeaturesParams) SetLimit(limit *int64) {
	o.Limit = limit
}

// WithNames adds the names to the get API features params
func (o *GetAPIFeaturesParams) WithNames(names []string) *GetAPIFeaturesParams {
	o.SetNames(names)
	return o
}

// SetNames adds the names to the get API features params
func (o *GetAPIFeaturesParams) SetNames(names []string) {
	o.Names = names
}

// WithOffset adds the offset to the get API features params
func (o *GetAPIFeaturesParams) WithOffset(offset *int64) *GetAPIFeaturesParams {
	o.SetOffset(offset)
	return o
}

// SetOffset adds the offset to the get API features params
func (o *GetAPIFeaturesParams) SetOffset(offset *int64) {
	o.Offset = offset
}

// WriteToRequest writes these params to a swagger request
func (o *GetAPIFeaturesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	r.SetTimeout(o.timeout)
	var res []error

	if o.Limit != nil {

		// query param limit
		var qrLimit int64
		if o.Limit != nil {
			qrLimit = *o.Limit
		}
		qLimit := swag.FormatInt64(qrLimit)
		if qLimit != "" {
			if err := r.SetQueryParam("limit", qLimit); err != nil {
				return err
			}
		}

	}

	valuesNames := o.Names

	joinedNames := swag.JoinByFormat(valuesNames, "")
	// query array param names
	if err := r.SetQueryParam("names", joinedNames...); err != nil {
		return err
	}

	if o.Offset != nil {

		// query param offset
		var qrOffset int64
		if o.Offset != nil {
			qrOffset = *o.Offset
		}
		qOffset := swag.FormatInt64(qrOffset)
		if qOffset != "" {
			if err := r.SetQueryParam("offset", qOffset); err != nil {
				return err
			}
		}

	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}