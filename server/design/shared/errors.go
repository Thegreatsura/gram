package shared

import (
	. "goa.design/goa/v3/dsl"

	"github.com/speakeasy-api/gram/internal/oops"
)

// DeclareErrorResponses declares the error responses at the service and method
// level including transport mappings (i.e. HTTP status codes).
func DeclareErrorResponses() {
	Error(string(oops.CodeUnauthorized), func() { Description(oops.CodeUnauthorized.UserMessage()) })
	Error(string(oops.CodeForbidden), func() { Description(oops.CodeForbidden.UserMessage()) })
	Error(string(oops.CodeBadRequest), func() { Description(oops.CodeBadRequest.UserMessage()) })
	Error(string(oops.CodeNotFound), func() { Description(oops.CodeNotFound.UserMessage()) })
	Error(string(oops.CodeConflict), func() { Description(oops.CodeConflict.UserMessage()) })
	Error(string(oops.CodeUnsupportedMedia), func() { Description(oops.CodeUnsupportedMedia.UserMessage()) })
	Error(string(oops.CodeInvalid), func() { Description(oops.CodeInvalid.UserMessage()) })
	Error(string(oops.CodeInvariantViolation), func() {
		Description(oops.CodeInvariantViolation.UserMessage())
		Fault()
	})
	Error(string(oops.CodeUnexpected), func() {
		Description(oops.CodeUnexpected.UserMessage())
		Fault()
	})

	HTTP(func() {
		Response(string(oops.CodeUnauthorized), StatusUnauthorized, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeForbidden), StatusForbidden, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeBadRequest), StatusBadRequest, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeNotFound), StatusNotFound, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeConflict), StatusConflict, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeUnsupportedMedia), StatusUnsupportedMediaType, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeInvalid), StatusUnprocessableEntity, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeInvariantViolation), StatusInternalServerError, func() {
			ContentType("application/json")
		})
		Response(string(oops.CodeUnexpected), StatusInternalServerError, func() {
			ContentType("application/json")
		})
	})
}
