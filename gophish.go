package main

import (
	"github.com/ghodss/yaml"
	"github.com/fatih/color"
	"encoding/json"
	"io/ioutil"
	"os/signal"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"bytes"
	"time"
	"fmt"
	"os"
)

func gopishPrintColor(s string,c string){
	if c == "warn" {
		fmt.Println("[ gophish ]",color.YellowString(s))
	} else if c == "info" {
		fmt.Println("[ gophish ]",color.BlueString(s))
	} else if c == "error" {
		fmt.Println("[ gophish ]",color.RedString(s))
	} else if c == "details" {
		fmt.Println("[ gophish ]",color.GreenString(s))
	}
}

func gPrint(s string, c string){
	switch runtime.GOOS {
	case "windows":
		fmt.Println("[ gophish ]",s)
	default:
		gopishPrintColor(s,c)
		
	}
}

func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

func arrayHas(arr1 []string, arr2 []string) bool {
	fmt.Println(arr1,arr2)
	for a1 := range arr1 {
		for a2 := range arr2 {
			if arr1[a1] == arr2[a2] {
				return true
			}
		}
	}

	return false
}

func dirname() string {
	dir, err := os.Getwd()
    checkErr(err)
    return dir
}

func getTrackingURL(url string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	b := []byte(`{"url":"http://`+url+`"}`)
	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	// fmt.Println(string(b))
	
	response, err := client.Post("http://ki.tc/","application/json",bytes.NewReader(b))

	checkErr(err)
	// fmt.Println(response.Body)
	b,err = ioutil.ReadAll(response.Body)
	checkErr(err)
	// fmt.Println(string(b))
	err = json.Unmarshal(b,&res)
	checkErr(err)
	return res
}

func checkCapturedUsersFromAdminLink(timeout int, link string){
	gPrint("ADMIN url => "+link,"details")
	// done <- true
	quit := make(chan bool)
	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			gPrint(sig.String(),"error")
			quit <- true
			os.Exit(0)
		}
	}()
	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	status := "clicked 0 times"
	func(){
		go func(){
			for t := range ticker.C {
				var b []byte
				client := http.Client{
					Timeout: 30 * time.Second,
				}
				hour := strconv.Itoa(t.Hour())
				min := strconv.Itoa(t.Minute())
				sec := strconv.Itoa(t.Second())
				if len(min) == 1{
					min = "0"+min
				}
				if len(hour) == 1 {
					hour = "0"+hour
				}
				if len(sec) == 1 {
					sec = "0"+sec
				}
				res, err := client.Get(link)
				checkErr(err)
				b,err = ioutil.ReadAll(res.Body)
				checkErr(err)
				OBJ := make(map[string]interface{})
				err = json.Unmarshal(b,&OBJ)
				checkErr(err)
				var keys []string
				for k := range OBJ {
					// gPrint(k)
					if k!="_id"&&k!="url_build"&&k!="url_id"{
						// gPrint(OBJ[k])
						keys = append(keys,k)
					}
				}
				for i,id := range keys {
					b,err = yaml.Marshal(OBJ[id])
					checkErr(err)
					separator := "\\"
					if runtime.GOOS != "windows" {
						separator = "/"
					}
					fn := dirname()+separator+"target"+OBJ["_id"].(string)+"#"+strconv.Itoa(i)+".yaml"
					if _,err := os.Stat(fn); os.IsNotExist(err) {
						gPrint("writing data to "+fn,"details")
						_,err = os.Create(fn)
						checkErr(err)
						err = ioutil.WriteFile(fn,b,0644)
						checkErr(err)
					}
				}
				status = "clicked " + strconv.Itoa(len(keys)) + " times"
				gPrint("["+hour+":"+min+":"+sec+"] link_status ("+status+")","info")

				// gPrint(string(b))
			}
		}()
	}()

	<- quit
}

func checkCapturedUsers(o map[string]map[string]string, timeout int){
	gPrint("ADMIN url    => "+o["url_short"]["admin_link"],"details")
	gPrint("TRANSER_ID   => "+o["url_short"]["_id"],"details")
	gPrint("Logging link => "+o["url_short"]["link"]+" (Send this to target)","details")
	// done <- true
	quit := make(chan bool)
	c := make(chan os.Signal,1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			gPrint(sig.String(),"error")
			quit <- true
			os.Exit(0)
		}
	}()
	ticker := time.NewTicker(time.Duration(timeout) * time.Second)
	status := "clicked 0 times"
	func(){
		go func(){
			for t := range ticker.C {
				var b []byte
				client := http.Client{
					Timeout: 30 * time.Second,
				}
				hour := strconv.Itoa(t.Hour())
				min := strconv.Itoa(t.Minute())
				sec := strconv.Itoa(t.Second())
				if len(min) == 1{
					min = "0"+min
				}
				if len(hour) == 1 {
					hour = "0"+hour
				}
				if len(sec) == 1 {
					sec = "0"+sec
				}
				res, err := client.Get(o["url_short"]["admin_link"])
				checkErr(err)
				b,err = ioutil.ReadAll(res.Body)
				checkErr(err)
				OBJ := make(map[string]interface{})
				err = json.Unmarshal(b,&OBJ)
				checkErr(err)
				var keys []string
				for k := range OBJ {
					// gPrint(k)
					if k!="_id"&&k!="url_build"&&k!="url_id"{
						// gPrint(OBJ[k])
						keys = append(keys,k)
					}
				}
				for i,id := range keys {
					b,err = yaml.Marshal(OBJ[id])
					checkErr(err)
					separator := "\\"
					if runtime.GOOS != "windows" {
						separator = "/"
					}
					fn := dirname()+separator+"target"+OBJ["_id"].(string)+"#"+strconv.Itoa(i)+".yaml"
					if _,err := os.Stat(fn); os.IsNotExist(err) {
						gPrint("writing data to "+fn,"details")
						_,err = os.Create(fn)
						checkErr(err)
						err = ioutil.WriteFile(fn,b,0644)
						checkErr(err)
					}
				}
				status = "clicked " + strconv.Itoa(len(keys)) + " times"
				gPrint("["+hour+":"+min+":"+sec+"] link_status ("+status+")","info")

				// gPrint(string(b))
			}
		}()
	}()

	<- quit
}

func genArgsMap(s []string) map[string]string {
	m := make(map[string]string)
	for _,i := range s {
		if strings.HasPrefix(i,"--") {
			kvp := strings.Split(i,"=")
			if len(kvp) == 1 {
				kvp = append(kvp,"true")
			}
			m[kvp[0][2:]] = kvp[1]
		}
	}
	return m
}

func main(){
	fmt.Println("Gophish v0.0.1\n")
	args := os.Args[1:]
	if len(args) == 0 {
		str := "gophish <url> " + "\n\t\turl -> url to redirect target after getting information (url should not start with http:// )" + 
				"\n\t\toptions\n\t\t\t--admin -> Admin link provided by first launch to continue (url not needed)" +
				"\n\t\t\t--timeout -> optional interval for checking if the linked was clicked"
		gPrint(str,"info")
		os.Exit(0)
	}
	m := genArgsMap(args)
	timeout := 10
	if m["timeout"] != "" {
		var err error
		timeout, err = strconv.Atoi(m["timeout"])
		checkErr(err)
	}
	if m["admin"] != "" {
		gPrint("Resuming from Admin url "+m["admin"],"info")
		checkCapturedUsersFromAdminLink(timeout,m["admin"])
	} else {
		jsonObject := getTrackingURL(args[0])
		gPrint("Creating ip logging url","info")
		// ticker := time.NewTicker(1 * time.Second)
		// gPrint("["+strconv.Itoa(t.Hour())+":"+strconv.Itoa(t.Minute())+"] check status")
		checkCapturedUsers(jsonObject,timeout)
		// gPrint(jsonString)
		// gPrint(jsonObject)
	}
}