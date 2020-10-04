package api

import (
	"NewsAggregator/pkg/aggregator"
	"encoding/json"
	"fmt"
	"net/http"
)

func StartServer(p int, a aggregator.Aggregator) error {
	http.Handle("/latest", latestHandler(a))

	return http.ListenAndServe(fmt.Sprintf(":%d", p), nil)
}

func latestHandler(agg aggregator.Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := agg.GetLatestData()
		bytes, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(bytes)
	}
}
