package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type Response struct {
	StatusCode int
	Body       any `json:"body,omitempty"`
}

func (r *Response) SendHasFile(w http.ResponseWriter) error {
	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		return errors.New("response body is empty")
	}

	if r.StatusCode >= 400 {
		body, ok := r.Body.(*ErrorBody)

		if !ok {
			return errors.New("response body is not an error")
		}

		_, err := w.Write([]byte(body.Message))

		return err
	}

	_, err := w.Write(r.Body.([]byte))

	return err
}

func (r *Response) SendHasJson(w http.ResponseWriter) error {
	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		return errors.New("response body is empty")
	}

	return json.NewEncoder(w).Encode(r.Body)
}

func (r *Response) SendHasQRCode(w http.ResponseWriter) error {
	if r.Body == nil {
		err := errors.New("response body is empty")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if r.StatusCode >= 400 {
		body, ok := r.Body.(*ErrorBody)

		if !ok {
			err := errors.New("response body is not an error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return err
		}

		http.Error(w, body.Message, r.StatusCode)
		return errors.New(body.Message)
	}

	qrBody, ok := r.Body.(QRBody)

	if !ok {
		err := errors.New("response body is not a QR code")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	bc, err := qr.Encode(qrBody.Content, qrBody.Level, qrBody.Mode)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	qrc, err := barcode.Scale(bc, qrBody.Width, qrBody.Height)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(r.StatusCode)
	w.Header().Set("Content-Type", "image/png")

	return png.Encode(w, qrc)
}

func (r *Response) SendHasTemplate(w http.ResponseWriter) error {
	if r.Body == nil {
		return fmt.Errorf("response body is empty")
	}

	//render template using html/template
	//todo implemente template rendering

	return nil
}

func NewResponseError(statusCode int, err error) (*Response, error) {
	return NewError(statusCode, err), err
}

func NewResponse(body any) (*Response, error) {
	return &Response{
		StatusCode: http.StatusOK,
		Body:       body,
	}, nil
}
