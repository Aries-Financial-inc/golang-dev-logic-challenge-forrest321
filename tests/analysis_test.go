package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-forrest321/models"
	"github.com/Aries-Financial-inc/golang-dev-logic-challenge-forrest321/routes"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// getTestJson loads the appropriate test data as a byte slice of json
func getTestJson(valid bool) ([]byte, error) {
	filename := "../testdata/testdata.json"
	if !valid {
		filename = "../testdata/invalidtestdata.json"
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// buildTestData builds a slice of OptionsContracts
func buildTestData(valid bool) ([]models.OptionsContract, error) {
	b, err := getTestJson(valid)
	if err != nil {
		return nil, err
	}
	var contracts []models.OptionsContract
	if err = json.Unmarshal(b, &contracts); err != nil {
		return nil, err
	}
	return contracts, nil
}

// TestOptionContractModelValidation tests the validation to ensure each contract is valid
func TestOptionsContractModelValidation(t *testing.T) {
	validContracts, err := buildTestData(true)
	if err != nil {
		t.Logf("test data error")
	}

	errorsFound := false
	for i, c := range validContracts {
		errs := c.Validate()
		if len(errs) > 0 {
			t.Log(fmt.Sprintf("error in contract #%v, err:%v", i, errs))
			errorsFound = true
		}
	}
	if errorsFound {
		t.Logf("Errors found during validation of testdata.json")
	}

	invalidContracts, err := buildTestData(false)
	if err != nil {
		t.Logf("test data error: %v", err)
	}

	errorsFound = false
	errorCount := 0
	for _, c := range invalidContracts {
		errs := c.Validate()
		if len(errs) > 0 {
			errorsFound = true
			errorCount++
		}
	}
	if !errorsFound && errorCount != 4 {
		t.Logf("Errors not found during validation of invalidtestdata.json")
	}
}

// TestAnalysisEndpoint ensures that the endpoint properly handles valid and invalid requests
func TestAnalysisEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	b, err := getTestJson(true)
	if err != nil {
		t.Logf("error loading test json")
		return
	}

	req, err := http.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(b))
	if err != nil {
		t.Logf("error building http request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Logf(fmt.Sprintf("error in endpoint /analyze: %v", w.Body.String()))
	}

	badTestData, err := getTestJson(false)
	if err != nil {
		t.Logf("error loading test json")
		return
	}
	badreq, err := http.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(badTestData))
	if err != nil {
		t.Logf("error building http request")
		return
	}
	badreq.Header.Set("Content-Type", "application/json")

	badW := httptest.NewRecorder()
	router.ServeHTTP(badW, badreq)
	if badW.Code != http.StatusBadRequest {
		t.Logf(fmt.Sprintf("error in endpoint /analyze validation. expected status 400, got: %v", w.Code))
	}
}

// TestIntegration validates that the endpoint runs and provides correct results
func TestIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := routes.SetupRouter()

	b, err := getTestJson(true)
	if err != nil {
		t.Logf("error loading test json")
		return
	}

	req, err := http.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(b))
	if err != nil {
		t.Logf("error building http request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Logf(fmt.Sprintf("error in endpoint /analyze: %v", w.Body.String()))
	}

	validResults := "{\"xy_values\":[{\"x\":100,\"y\":10.05},{\"x\":102.5,\"y\":12.1},{\"x\":103,\"y\":15.5},{\"x\":105,\"y\":18}],\"max_profit\":90.4,\"max_loss\":87.5,\"break_even_points\":[110.05,114.6,118.5,87]}"
	if w.Body.String() != validResults {
		t.Logf("incorrect results: %s", w.Body.String())
	}
}
