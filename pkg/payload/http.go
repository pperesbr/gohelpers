package payload

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// MakeRequest cria um http.Request com payload gerado
func MakeRequest(t *testing.T, method, url string, fields []FieldDef) (*http.Request, error) {
	t.Helper()

	reader, err := GenerateReader(fields)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// MakeRequestWithRecorder cria request + ResponseRecorder
func MakeRequestWithRecorder(t *testing.T, method, url string, fields []FieldDef) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	req, err := MakeRequest(t, method, url, fields)
	require.NoError(t, err, "failed to create request")

	rec := httptest.NewRecorder()

	return req, rec
}

// MakeJSONRequest cria request com JSON arbitrário (não usa FieldDef)
func MakeJSONRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()

	var reader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		require.NoError(t, err, "failed to marshal JSON")
		reader = bytes.NewReader(jsonBytes)
	}

	req := httptest.NewRequest(method, url, reader)
	req.Header.Set("Content-Type", "application/json")

	return req
}

// ParseJSONResponse decodifica response JSON
func ParseJSONResponse(t *testing.T, rec *httptest.ResponseRecorder, target interface{}) {
	t.Helper()

	err := json.NewDecoder(rec.Body).Decode(target)
	require.NoError(t, err, "failed to decode response")
}

// ParseErrorResponse decodifica response de erro padrão
func ParseErrorResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()

	var response map[string]interface{}
	ParseJSONResponse(t, rec, &response)
	return response
}

// AssertErrorDetails verifica detalhes de erro de validação
func AssertErrorDetails(t *testing.T, rec *httptest.ResponseRecorder, field, expectedMessage string) {
	t.Helper()

	response := ParseErrorResponse(t, rec)

	details, ok := response["details"].(map[string]interface{})
	require.True(t, ok, "response should have details field")

	message, ok := details[field].(string)
	require.True(t, ok, "field %s should exist in details", field)
	require.Contains(t, message, expectedMessage, "error message mismatch for field %s", field)
}

// gohelpers/pkg/payload/http.go

// AssertStatusCode verifica status com mensagem clara
func AssertStatusCode(t *testing.T, rec *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if rec.Code != expected {
		var body interface{}
		json.NewDecoder(rec.Body).Decode(&body)
		t.Fatalf("expected status %d but got %d. Body: %+v", expected, rec.Code, body)
	}
}

// AssertCreated verifica 201 + ID no response
func AssertCreated(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()

	AssertStatusCode(t, rec, 201)

	var response map[string]string
	ParseJSONResponse(t, rec, &response)

	id, ok := response["id"]
	require.True(t, ok, "response should have id field")
	require.NotEmpty(t, id, "id should not be empty")

	return id
}

// AssertNoContent verifica 204
func AssertNoContent(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	AssertStatusCode(t, rec, 204)
}

// AssertError verifica erro com mensagem específica
func AssertError(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedMessage string) {
	t.Helper()

	AssertStatusCode(t, rec, expectedStatus)

	response := ParseErrorResponse(t, rec)
	errorMsg, ok := response["error"].(string)
	require.True(t, ok, "response should have error field")
	require.Contains(t, errorMsg, expectedMessage)
}
