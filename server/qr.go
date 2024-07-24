package server

import "github.com/boombuler/barcode/qr"

type QRBody struct {
	Height  int
	Width   int
	Level   qr.ErrorCorrectionLevel
	Mode    qr.Encoding
	Content string
}

func (q *QRBody) SetSize(height, width int) *QRBody {
	q.Height = height
	q.Width = width

	return q
}

func NewQRCode(content string) QRBody {
	return QRBody{
		Height:  300,
		Width:   300,
		Level:   qr.L,
		Mode:    qr.Auto,
		Content: content,
	}
}
