package server

import (
	"gAD-System/services/gad-manager/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRepoCalculate struct{}

func (m MockRepoCalculate) DoCalculate([]string) ([]string, error) {
	return []string{}, nil
}

func TestHandlers_calculate_valid(t *testing.T) {
	var mock MockRepoCalculate
	calc := domain.NewCalculator(mock)
	hs := InitServerHandlers(calc)
	router := newRouter(hs)

	validJSON := `{
		"exprs": ["100+100", "20+30"]
	}`
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/calc", strings.NewReader(validJSON))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, r)

	assert.Equal(t, 200, w.Code)
}

func TestHandlers_calculate_invalid(t *testing.T) {
	var mock MockRepoCalculate
	calc := domain.NewCalculator(mock)
	hs := InitServerHandlers(calc)

	invalidJSON := `{
		"asd": ["100+100"]
	}`
	w := httptest.NewRecorder()
	w.Body.WriteString(invalidJSON)
	c, _ := gin.CreateTestContext(w)

	hs.calculate(c)

	assert.Equal(t, 400, w.Code)
}
