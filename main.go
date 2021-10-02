package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

type ExecuteRequest struct {
	Sql string `json:"sql"`
}

func main() {
	http.HandleFunc("/execute", func(rw http.ResponseWriter, r *http.Request) {
		bodyRaw, err := ioutil.ReadAll(r.Body)
		checkErr(err)
		fmt.Printf("SQL execute request: %v\n", string(bodyRaw))
		body := &ExecuteRequest{}
		err = json.Unmarshal(bodyRaw, body)
		checkErr(err)

		cmd := exec.Command("sqlite3", "./target.db")
		inPipe, err := cmd.StdinPipe()
		checkErr(err)

		_, err = io.WriteString(inPipe, body.Sql)
		checkErr(err)

		err = inPipe.Close()
		checkErr(err)

		output, _ := cmd.CombinedOutput()
		// checkErr(err)

		// rw.Header().Add("Content-Type", "application/json")
		rw.Write(output)
	})
	http.ListenAndServe(":8080", nil)
}
