package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetURLParam(r *http.Request, key string) string {
	v, ok := mux.Vars(r)[key]

	if !ok {
		return ""
	}

	return v
}

func GetURLParamWithDefault(r *http.Request, key, defaultValue string) string {
	v := GetURLParam(r, key)

	if v == "" {
		return defaultValue
	}

	return v
}

func GetURLParamHasInt(r *http.Request, key string) int {
	s := GetURLParam(r, key)

	if s == "" {
		return 0
	}

	i, err := strconv.Atoi(s)

	if err != nil {
		return 0
	}

	return i
}

func GetURLParamHasIntWithDefault(r *http.Request, key string, defaultValue int) int {
	i := GetURLParamHasInt(r, key)

	if i == 0 {
		return defaultValue
	}

	return i
}
