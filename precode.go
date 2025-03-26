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
    totalCount := 4 

    req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
    require.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // Проверяем, что код ответа 200
    require.Equal(t, http.StatusOK, responseRecorder.Code, "Ожидался код 200")

    // Проверяем, что количество кафе в ответе соответствует totalCount
    cafes := strings.Split(responseRecorder.Body.String(), ",")
    require.Len(t, cafes, totalCount, "Ожидалось, что в ответе будет 4 кафе")
}

func TestMainHandlerWhenCityNotSupported(t *testing.T) {
    req, err := http.NewRequest("GET", "/cafe?count=2&city=unknown", nil)
    require.NoError(t, err)

    responseRecorder := httptest.NewRecorder()
    handler := http.HandlerFunc(mainHandle)
    handler.ServeHTTP(responseRecorder, req)

    // Проверяем, что код ответа 400
    require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Ожидался код 400")

    // Проверяем, что сообщение об ошибке корректное
    require.Equal(t, "wrong city value", responseRecorder.Body.String(), "Ожидалось сообщение 'wrong city value'")
}