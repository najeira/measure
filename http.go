package measure

import (
	"fmt"
	"net/http"
)

func HandleStats(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	key := q.Get("key")

	stats := GetStats()
	stats.SortDesc(key)

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "key,count,sum,min,max,avg,rate,p95\n")
	for _, s := range stats {
		fmt.Fprintf(w, "%s,%d,%f,%f,%f,%f,%f,%f\n",
			s.Key, s.Count, s.Sum, s.Min, s.Max, s.Avg, s.Rate, s.P95)
	}
}
