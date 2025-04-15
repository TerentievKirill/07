package main

import (
    "net/http"
    "net/http/httptest"
    "net/url"
    "strings"
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

func TestMainHandlerWhenCountMoreThanTotal_Implemented(t *testing.T) {
    q := url.Values{}
    q.Set("count", "10")
    q.Set("city", "moscow")

    req := httptest.NewRequest(http.MethodGet, "/cafe?"+q.Encode(), nil)
    resp := httptest.NewRecorder()

    mainHandle(resp, req)

    require.Equal(t, http.StatusOK, resp.Code)
    expected := strings.Join(cafeList["moscow"], ",")
    assert.Equal(t, expected, resp.Body.String())
}
