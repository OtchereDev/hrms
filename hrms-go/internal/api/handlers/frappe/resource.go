package frappe

import (
	"encoding/json"
	"strconv"

	"github.com/gofiber/fiber/v2"

	frappeCore "github.com/OtchereDev/hrms-go/internal/core/frappe"
)

// ResourceHandler handles Frappe-compatible resource (DocType) operations
// This mimics Frappe's /api/resource/:doctype endpoints
type ResourceHandler struct {
	doctypeService *frappeCore.DocTypeService
}

// NewResourceHandler creates a new resource handler
func NewResourceHandler(doctypeService *frappeCore.DocTypeService) *ResourceHandler {
	return &ResourceHandler{
		doctypeService: doctypeService,
	}
}

// Get retrieves a single document
// GET /api/resource/:doctype/:name
func (h *ResourceHandler) Get(c *fiber.Ctx) error {
	doctype := c.Params("doctype")
	name := c.Params("name")

	if doctype == "" || name == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
			{Field: "name", Message: "Document name is required"},
		})
	}

	doc, err := h.doctypeService.Get(c.Context(), doctype, name)
	if err != nil {
		if err.Error() == "document not found: "+doctype+" "+name {
			return frappeCore.SendNotFound(c, err.Error())
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, doc)
}

// GetList retrieves a list of documents with filters
// GET /api/resource/:doctype
func (h *ResourceHandler) GetList(c *fiber.Ctx) error {
	doctype := c.Params("doctype")

	if doctype == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
		})
	}

	// Parse filters from query parameter
	filtersStr := c.Query("filters", "{}")
	var filters map[string]interface{}
	if err := json.Unmarshal([]byte(filtersStr), &filters); err != nil {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "filters", Message: "Invalid filters format"},
		})
	}

	// Parse fields
	fieldsStr := c.Query("fields", "[]")
	var fields []string
	if err := json.Unmarshal([]byte(fieldsStr), &fields); err != nil {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "fields", Message: "Invalid fields format"},
		})
	}

	// Parse pagination
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	offset := (page - 1) * limit

	// Also support limit_start and limit_page_length (Frappe's pagination params)
	if limitStart := c.Query("limit_start"); limitStart != "" {
		offset, _ = strconv.Atoi(limitStart)
	}
	if limitPageLength := c.Query("limit_page_length"); limitPageLength != "" {
		limit, _ = strconv.Atoi(limitPageLength)
	}

	// Get documents
	docs, total, err := h.doctypeService.GetList(c.Context(), doctype, filters, fields, limit, offset)
	if err != nil {
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendList(c, docs, total, page, limit)
}

// Insert creates a new document
// POST /api/resource/:doctype
func (h *ResourceHandler) Insert(c *fiber.Ctx) error {
	doctype := c.Params("doctype")

	if doctype == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
		})
	}

	// Parse request body
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "body", Message: "Invalid request body"},
		})
	}

	// Create document
	doc, err := h.doctypeService.Insert(c.Context(), doctype, data)
	if err != nil {
		// Check for duplicate errors
		if err.Error() == "UNIQUE constraint failed" || err.Error() == "duplicate key value" {
			return frappeCore.SendDuplicateError(c, "Document already exists")
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, doc)
}

// Update updates an existing document
// PUT /api/resource/:doctype/:name
func (h *ResourceHandler) Update(c *fiber.Ctx) error {
	doctype := c.Params("doctype")
	name := c.Params("name")

	if doctype == "" || name == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
			{Field: "name", Message: "Document name is required"},
		})
	}

	// Parse request body
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "body", Message: "Invalid request body"},
		})
	}

	// Update document
	doc, err := h.doctypeService.Save(c.Context(), doctype, name, data)
	if err != nil {
		if err.Error() == "document not found: "+doctype+" "+name {
			return frappeCore.SendNotFound(c, err.Error())
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, doc)
}

// Delete deletes a document
// DELETE /api/resource/:doctype/:name
func (h *ResourceHandler) Delete(c *fiber.Ctx) error {
	doctype := c.Params("doctype")
	name := c.Params("name")

	if doctype == "" || name == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
			{Field: "name", Message: "Document name is required"},
		})
	}

	// Delete document
	if err := h.doctypeService.Delete(c.Context(), doctype, name); err != nil {
		if err.Error() == "document not found: "+doctype+" "+name {
			return frappeCore.SendNotFound(c, err.Error())
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, map[string]string{"message": "Document deleted successfully"})
}

// Submit submits a document (workflow action)
// POST /api/resource/:doctype/:name/submit
func (h *ResourceHandler) Submit(c *fiber.Ctx) error {
	doctype := c.Params("doctype")
	name := c.Params("name")

	if doctype == "" || name == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
			{Field: "name", Message: "Document name is required"},
		})
	}

	// Submit document
	doc, err := h.doctypeService.Submit(c.Context(), doctype, name)
	if err != nil {
		if err.Error() == "document not found: "+doctype+" "+name {
			return frappeCore.SendNotFound(c, err.Error())
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, doc)
}

// Cancel cancels a document (workflow action)
// POST /api/resource/:doctype/:name/cancel
func (h *ResourceHandler) Cancel(c *fiber.Ctx) error {
	doctype := c.Params("doctype")
	name := c.Params("name")

	if doctype == "" || name == "" {
		return frappeCore.SendValidationError(c, []frappeCore.ValidationError{
			{Field: "doctype", Message: "DocType is required"},
			{Field: "name", Message: "Document name is required"},
		})
	}

	// Cancel document
	doc, err := h.doctypeService.Cancel(c.Context(), doctype, name)
	if err != nil {
		if err.Error() == "document not found: "+doctype+" "+name {
			return frappeCore.SendNotFound(c, err.Error())
		}
		return frappeCore.SendInternalError(c, err.Error())
	}

	return frappeCore.SendSuccess(c, doc)
}
