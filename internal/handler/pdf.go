package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/models"
	"github.com/Engineering-Association-UofK/Technical-Office-Bot/internal/service"
)

type HelperHandler struct {
	PDFService service.PDFService
}

func NewHelper(pdf service.PDFService) *HelperHandler {
	return &HelperHandler{
		PDFService: pdf,
	}
}

func (hh *HelperHandler) HandlePDFRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)

	var pdfReq models.PDFRequest
	if err := json.Unmarshal(body, &pdfReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	pdfBytes, err := hh.PDFService.GenerateFromHTML(r.Context(), pdfReq.HTMLContent)
	if err != nil {
		http.Error(w, "PDF Generation Failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=certificate.pdf")
	w.Write(pdfBytes)
}
