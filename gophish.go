package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
)

//Print coloured output
func gopishPrintColor(s string, c string) {
	if c == "warn" {
		fmt.Println("[ gophish ]", color.YellowString(s))
	} else if c == "info" {
		fmt.Println("[ gophish ]", color.BlueString(s))
	} else if c == "error" {
		fmt.Println("[ gophish ]", color.RedString(s))
	} else if c == "details" {
		fmt.Println("[ gophish ]", color.GreenString(s))
	}
}

//Wrapper for gophishPrintColor in case of colorless terminal (cmd)
func gPrint(s string, c string) {
	switch runtime.GOOS {
	case "windows":
		fmt.Println("[ gophish ]", s)
	default:
		gopishPrintColor(s, c)

	}
}

//Check errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Get directory name
func dirname() string {
	dir, err := os.Getwd()
	checkErr(err)
	return dir
}

//Retrieve the tracking url from ki.tc
func getTrackingURL(url string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	b := []byte(`{"url":"http://` + url + `"}`)
	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	//Request hooked URL
	response, err := client.Post("http://ki.tc/", "application/json", bytes.NewReader(b))
	checkErr(err)

	//Convert byte stream to JSON
	b, err = ioutil.ReadAll(response.Body)
	checkErr(err)
	err = json.Unmarshal(b, &res)
	checkErr(err)
	return res
}

//Monitor clicks from a specific link e.g when resuming a session
func resumeSessionFromAdminLink(timeout int, link string) {
	gPrint("ADMIN url => "+link, "details")
	quit := make(chan bool)
	c := make(chan os.Signal, 1)

	//Exit on interrupt
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			gPrint(sig.String(), "error")
			quit <- true
			os.Exit(0)
		}
	}()

	//Time counter
	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	status := "clicked 0 times"

	//Asynchronous go routine for checking link status
	go func() {
		for t := range ticker.C {
			var b []byte

			//Request options
			client := http.Client{
				Timeout: 30 * time.Second,
			}

			//Get time
			hour := strconv.Itoa(t.Hour())
			min := strconv.Itoa(t.Minute())
			sec := strconv.Itoa(t.Second())
			if len(min) == 1 {
				min = "0" + min
			}
			if len(hour) == 1 {
				hour = "0" + hour
			}
			if len(sec) == 1 {
				sec = "0" + sec
			}

			//Retrieve link status
			res, err := client.Get(link)
			checkErr(err)
			b, err = ioutil.ReadAll(res.Body)
			checkErr(err)
			OBJ := make(map[string]interface{})
			err = json.Unmarshal(b, &OBJ)
			checkErr(err)

			//Get captured user keys [Hashes]
			var keys []string
			for k := range OBJ {
				if k != "_id" && k != "url_build" && k != "url_id" {
					keys = append(keys, k)
				}
			}

			//For each key write out .yaml file containing the information
			for i, id := range keys {
				b, err = yaml.Marshal(OBJ[id])
				checkErr(err)
				separator := "\\"
				if runtime.GOOS != "windows" {
					separator = "/"
				}

				//Output filename
				fn := dirname() + separator + "target" + OBJ["_id"].(string) + "#" + strconv.Itoa(i) + ".yaml"
				if _, err := os.Stat(fn); os.IsNotExist(err) {
					gPrint("writing data to "+fn, "details")
					_, err = os.Create(fn)
					checkErr(err)
					err = ioutil.WriteFile(fn, b, 0644)
					checkErr(err)
				}
			}

			//Inform user on number of clicks
			status = "clicked " + strconv.Itoa(len(keys)) + " times"
			gPrint("["+hour+":"+min+":"+sec+"] link_status ("+status+")", "info")
		}
	}()

	<-quit
}

//Start new capture session
func newCaptureSession(o map[string]map[string]string, timeout int) {
	gPrint("ADMIN url    => "+o["url_short"]["admin_link"], "details")
	gPrint("TRANSER_ID   => "+o["url_short"]["_id"], "details")
	gPrint("Logging link => "+o["url_short"]["link"]+" (Send this to target)", "details")
	quit := make(chan bool)
	c := make(chan os.Signal, 1)

	//Exit on interrupt
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			gPrint(sig.String(), "error")
			quit <- true
			os.Exit(0)
		}
	}()

	//Time counter
	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	status := "clicked 0 times"

	//Asyncronous go routine for checking link status
	go func() {
		for t := range ticker.C {
			var b []byte

			//Request options
			client := http.Client{
				Timeout: 30 * time.Second,
			}

			//Get time
			hour := strconv.Itoa(t.Hour())
			min := strconv.Itoa(t.Minute())
			sec := strconv.Itoa(t.Second())
			if len(min) == 1 {
				min = "0" + min
			}
			if len(hour) == 1 {
				hour = "0" + hour
			}
			if len(sec) == 1 {
				sec = "0" + sec
			}

			//Retrieve link status
			res, err := client.Get(o["url_short"]["admin_link"])
			checkErr(err)
			b, err = ioutil.ReadAll(res.Body)
			checkErr(err)
			OBJ := make(map[string]interface{})
			err = json.Unmarshal(b, &OBJ)
			checkErr(err)

			//Get captured users keys [Hashes]
			var keys []string
			for k := range OBJ {
				if k != "_id" && k != "url_build" && k != "url_id" {
					keys = append(keys, k)
				}
			}

			//For each key write out .yaml file containing the information
			for i, id := range keys {
				b, err = yaml.Marshal(OBJ[id])
				checkErr(err)
				separator := "\\"
				if runtime.GOOS != "windows" {
					separator = "/"
				}

				//Output filename
				fn := dirname() + separator + "target" + OBJ["_id"].(string) + "#" + strconv.Itoa(i) + ".yaml"
				if _, err := os.Stat(fn); os.IsNotExist(err) {
					gPrint("writing data to "+fn, "details")
					_, err = os.Create(fn)
					checkErr(err)
					err = ioutil.WriteFile(fn, b, 0644)
					checkErr(err)
				}
			}

			//Inform user on number of clicks
			status = "clicked " + strconv.Itoa(len(keys)) + " times"
			gPrint("["+hour+":"+min+":"+sec+"] link_status ("+status+")", "info")
		}
	}()

	<-quit
}

//Convert os.Args into a map of string to string
func genArgsMap(s []string) map[string]string {
	m := make(map[string]string)
	for _, i := range s {
		if strings.HasPrefix(i, "--") {
			kvp := strings.Split(i, "=")
			if len(kvp) == 1 {
				kvp = append(kvp, "true")
			}
			m[kvp[0][2:]] = kvp[1]
		}
	}
	return m
}

//Help
func help() {
	str := "gophish <url> [options?] " + "\n\t\turl -> url to redirect target after getting information (url should not start with http:// )" +
		"\n\t\toptions\n\t\t\t--admin -> Admin link provided by first launch to continue (url not needed)" +
		"\n\t\t\t--timeout -> optional interval for checking if the linked was clicked"
	gPrint(str, "info")
	os.Exit(0)
}

func main() {
	//Greetings
	fmt.Print("Gophish v0.0.1\n\n")
	args := os.Args[1:]
	if len(args) == 0 {
		help()
	}
	m := genArgsMap(args)
	timeout := 10
	if m["help"] == "" || m["help"] != "" {
		help()
	}
	if m["timeout"] != "" {
		var err error
		timeout, err = strconv.Atoi(m["timeout"])
		checkErr(err)
	}
	if m["admin"] != "" {
		gPrint("Resuming from Admin url "+m["admin"], "info")
		resumeSessionFromAdminLink(timeout, m["admin"])
	} else {
		jsonObject := getTrackingURL(args[0])
		gPrint("Creating ip logging url", "info")
		newCaptureSession(jsonObject, timeout)
	}
}
