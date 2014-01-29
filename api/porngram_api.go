package main

import "fmt"
import s "strings"
import "bytes"
import "io/ioutil"
import "encoding/json"
import "net/http"
import "log"

type qword struct {
    qname string
    qstats map[string]float32
}

func main() {

	db, stats := make_dataset("../data/xhamster_sample.json")
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Inside HelloServer handler")
		query := s.Split(s.ToLower(s.Replace(req.URL.Path[1:], " ", "", -1)), ",")
		chResp := make(chan qword)

		for _, word := range query {
			a_word := qword{qname: word} 
			go query_word(a_word, &db, &stats, chResp)
		}

		res := make(map[string]map[string]float32)
		for i := len(query); i > 0; i-- {
			tmp := <- chResp
			res[tmp.qname] = tmp.qstats
		}
		js, _ := json.Marshal(res)
		fmt.Fprint(w, string(js))
	})
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}	
	
}

func query_word(w qword, db *map[string]*bytes.Buffer, stats *map[string]int, chResp chan qword) {
	
	count := make(map[string]int)
	res := make(map[string]float32)

	for year, txt := range *db {
		count[year] = s.Count(txt.String(), w.qname)
	}

	for year, total := range *stats {
		res[year] = (float32(count[year]) / float32(total)) * 100
	}

	w.qstats = res

	chResp <- w
}

func make_dataset(path string) (map[string]*bytes.Buffer, map[string]int) {
	
	dat, _ := ioutil.ReadFile(path)
	var jdat map[string]interface{}
	json.Unmarshal(dat, &jdat)

	res := make(map[string]*bytes.Buffer)
	stats := make(map[string]int)
	
	for _, v := range jdat {
	
		content := v.(map[string]interface{})
		date := s.Split(content["upload_date"].(string), "-")[0]
		valiDate := date != "NA" &&	date != "2007" && date != "2013"
	
		if valiDate {
			_, prs := res[date]
	
			if prs{
				res[date].WriteString(content["title"].(string))
				stats[date] += 1
	
				} else {
					res[date] = bytes.NewBufferString(content["title"].(string))
					stats[date] = 1
			}
		}
	}
	return res, stats
}
