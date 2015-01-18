package main

import ( 
	"fmt"
	"bytes"
	"encoding/hex"
	"net/http"
	"os"
	"os/exec"
	"time"
	"strconv"
	"flag"
)

func main() {
	//
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var Result string
	var Log string
	//
	var script = flag.String("script", "/path/to/foo.sh", "Set script")
	var jenkins_host = flag.String("host", "http://localhost:8080", "Set Jenkins hostname")
	var jenkins_job = flag.String("job", "hoge", "Set Jenkins job name")
	flag.Parse()
		
	start := time.Now().Unix()
	cmd := exec.Command(*script)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
    	if err != nil {
		Log = "\n" + string(stderr.String()) + "\n\n"
		Result = "1"
    	}else{
		Log = "\n" + string(stdout.String()) + "\n\n"
		Result = "0"
	}
	end := time.Now().Unix()
	Duration := (end - start) * 1000
	//
	sEnc := hex.EncodeToString([]byte(Log))
	body := bytes.NewBufferString("<run><log encoding=\"hexBinary\">" + sEnc + "</log><result>" + Result + "</result><duration>" + strconv.FormatInt(Duration, 10) + "</duration></run>")
	res, err := http.Post(*jenkins_host + "/job/" + *jenkins_job + "/postBuildResult", "text/xml; charset=utf-8", body)
	//
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("status: ", res.Status)
}
