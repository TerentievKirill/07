package main

import (
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
    "testing"
)

var cafeList = map[string][]string{
    "moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
    countStr := req.URL.Query().Get("count")
    if countStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("count missing"))
        return
    }

    count, err := strconv.Atoi(countStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong count value"))
        return
    }

    city := req.URL.Query().Get("city")

    cafe, ok := cafeList[city]
    if !ok {
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte("wrong city value"))
        return
    }

    if count > len(cafe) {
        count = len(cafe)
    }

    answer := strings.Join(cafe[:count], ",")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(answer))
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/cafe?count=10&city=moscow", nil)
    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    if responseRecorder.Code != http.StatusOK {
        t.Fatalf("Ожидался код 200, но получен %d", responseRecorder.Code)
    }

    expected := strings.Join(cafeList["moscow"], ",")
    result := responseRecorder.Body.String()
    if result != expected {
        t.Errorf("неожиданный ответ: %s, ожидался: %s", result, expected)
    }
}
