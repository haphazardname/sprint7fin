package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) { // проверяет работу сервера при разных значениях параметра count
	city := "moscow"
	r := "/cafe?city=moscow&count="
	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList[city])},
	}
	handler := http.HandlerFunc(mainHandle)

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		s := strings.TrimSpace(response.Body.String())
		sep := ","
		result := strings.Split(s, sep)
		filteredResult := make([]string, 0, len(result))
		for _, str := range result {
			if str != "" {
				filteredResult = append(filteredResult, str)
			}
		}
		l := len(filteredResult)
		assert.Equal(t, v.want, l)

	}

}

func TestCafeSearch(t *testing.T) {
	r := "/cafe?city=moscow&"
	requests := []struct {
		search string // передаваемое значение
		want   int    // ожидаемое количество кафе в ответе
	}{
		{"search=фасоль", 0},
		{"search=кофе", 2},
		{"search=вилка", 1},
	}

	handler := http.HandlerFunc(mainHandle)

	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", r+v.search, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
		s := strings.TrimSpace(response.Body.String())
		sep := ","
		result := strings.Split(s, sep)
		filteredResult := make([]string, 0, len(result))
		for _, str := range result {
			if str != "" {
				filteredResult = append(filteredResult, str)
			}
		}
		l := len(filteredResult)
		assert.Equal(t, v.want, l)

	}
}
