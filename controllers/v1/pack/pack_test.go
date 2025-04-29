package pack

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dekib/rePacks/internal/errors"
	"github.com/dekib/rePacks/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPacksHandler(t *testing.T) {
	// Setup test scenarios
	tests := []struct {
		name          string
		itemsNo       string
		packSizes     []int
		wantStatus    int
		wantPacks     []gin.H
		wantError     bool
		wantErrorCode int
	}{
		{
			name:       "Valid request - 1 item",
			itemsNo:    "1",
			packSizes:  []int{250, 500, 1000, 2000, 5000},
			wantStatus: http.StatusOK,
			wantPacks: []gin.H{
				{"pack": 250, "quantity": 1},
			},
		},
		{
			name:       "Valid request - 251 items",
			itemsNo:    "251",
			packSizes:  []int{250, 500, 1000, 2000, 5000},
			wantStatus: http.StatusOK,
			wantPacks: []gin.H{
				{"pack": 500, "quantity": 1},
			},
		},
		{
			name:       "Custom sizes - 500000 items",
			itemsNo:    "500000",
			packSizes:  []int{23, 31, 53},
			wantStatus: http.StatusOK,
			wantPacks: []gin.H{
				{"pack": 23, "quantity": 2},
				{"pack": 31, "quantity": 7},
				{"pack": 53, "quantity": 9429},
			},
		},
		{
			name:          "Invalid input - non-number",
			itemsNo:       "abc",
			packSizes:     []int{250, 500},
			wantStatus:    http.StatusBadRequest,
			wantError:     true,
			wantErrorCode: errors.CodeMalformedRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup service with test data
			svc := &services.PackService{}
			svc.SetPackSizes(tt.packSizes)

			// Create controller
			ctrl := NewPackController(svc)

			// Setup router
			router := gin.Default()
			router.GET("/packs", ctrl.GetPacks)

			// Create request
			req := httptest.NewRequest("GET", "/packs?items_no="+tt.itemsNo, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantError {
				// Verify error response
				var response struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}
				json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tt.wantErrorCode, response.Code)
			} else {
				// Verify success response
				var response struct {
					Packs []map[string]interface{} `json:"packs"`
				}
				json.Unmarshal(w.Body.Bytes(), &response)

				// Convert actual response to comparable format
				actualPacks := make([]gin.H, len(response.Packs))
				for i, p := range response.Packs {
					actualPacks[i] = gin.H{
						"pack":     int(p["pack"].(float64)),
						"quantity": int(p["quantity"].(float64)),
					}
				}

				assert.Equal(t, tt.wantPacks, actualPacks)
			}
		})
	}
}

func TestUpdatePackSizes(t *testing.T) {
	tests := []struct {
		name        string
		requestBody interface{}
		wantStatus  int
		wantCode    int
		wantMessage string
		wantSizes   []int
	}{
		{
			name: "Valid sizes update",
			requestBody: gin.H{
				"sizes": []int{23, 31, 53},
			},
			wantStatus:  http.StatusOK,
			wantMessage: "Pack sizes updated successfully",
			wantSizes:   []int{53, 31, 23},
		},
		{
			name: "Empty sizes array",
			requestBody: gin.H{
				"sizes": []int{},
			},
			wantStatus:  http.StatusBadRequest,
			wantCode:    errors.CodeValidationError,
			wantMessage: "Validation error",
		},
		{
			name: "Invalid sizes (non-positive)",
			requestBody: gin.H{
				"sizes": []int{0, -1},
			},
			wantStatus:  http.StatusBadRequest,
			wantCode:    errors.CodeValidationError,
			wantMessage: "Validation error",
		},
		{
			name:        "Malformed JSON",
			requestBody: "not a valid json",
			wantStatus:  http.StatusBadRequest,
			wantCode:    errors.CodeMalformedRequest,
			wantMessage: "Malformed request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			svc := &services.PackService{}
			ctrl := NewPackController(svc)

			router := gin.Default()
			router.PUT("/packs/sizes", ctrl.UpdatePackSizes)

			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/packs/sizes", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Verify status code
			assert.Equal(t, tt.wantStatus, w.Code)

			if tt.wantStatus == http.StatusOK {
				// Verify success response
				var response struct {
					Message   string `json:"message"`
					PackSizes []int  `json:"pack_sizes"`
				}
				json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tt.wantMessage, response.Message)
				assert.Equal(t, tt.wantSizes, response.PackSizes)
				assert.Equal(t, tt.wantSizes, svc.GetPackSizes())
			} else {
				// Verify error response
				var response struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				}
				json.Unmarshal(w.Body.Bytes(), &response)

				assert.Equal(t, tt.wantCode, response.Code)
				assert.Equal(t, tt.wantMessage, response.Message)
			}
		})
	}
}
