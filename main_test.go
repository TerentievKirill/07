package main

import (
    "net/http"
    "net/http/httptest"
    "net/url"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMainHandler_ValidRequest(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=moscow", nil)
    resp := httptest.NewRecorder()

    mainHandle(resp, req)

    require.Equal(t, http.StatusOK, resp.Code)
    assert.NotEmpty(t, resp.Body.String())
}

func TestMainHandler_WrongCity(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/cafe?count=2&city=paris", nil)
    resp := httptest.NewRecorder()

    mainHandle(resp, req)

    require.Equal(t, http.StatusBadRequest, resp.Code)
    assert.Equal(t, "wrong city value", resp.Body.String())
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    totalCount := 4 // Всего 4 кафе в Москве

    q := url.Values{}
    q.Set("count", "10")
    q.Set("city", "moscow")

    req := httptest.NewRequest(http.MethodGet, "/cafe?"+q.Encode(), nil)
    resp := httptest.NewRecorder()

    mainHandle(resp, req)

    require.Equal(t, http.StatusOK, resp.Code)
    result := resp.Body.String()

    cafes := cafeList["moscow"]
    assert.Len(t, cafes, totalCount)
    assert.Equal(t, cafes, splitResult(result))
}

// вспомогательная функция для сравнения списков
func splitResult(s string) []string {
    if s == "" {
        return []string{}
    }
    return strings.Split(s, ",")
}
