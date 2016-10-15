package issue

import (
	"net/http"
	"encoding/json"
)

func Handle(w http.ResponseWriter, err error, statusCode int) {
	if err != nil {
		e, _ := json.Marshal(Error{Error: err.Error()})

		w.WriteHeader(statusCode)
		w.Write(e)
	}
}
