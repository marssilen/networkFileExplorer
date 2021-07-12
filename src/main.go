package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	_ "github.com/lib/pq"
	_ "google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

var err error

var JSON_FILE string
var CONFIG_FILE string
var validate *validator.Validate
type Config struct {
	ServerId 		int64			`json:"server_id"`
	Location 		string			`json:"location"`
}
var config Config

func main() {
	dir := UserHomeDir()+"/vpnfiles"
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(dir)

	userOS := runtime.GOOS
	switch userOS {
	case "windows":
		JSON_FILE = dir+"/servers.json"
		CONFIG_FILE = dir+"/vpngo_config.json"
		fmt.Println("Windows")
		break
	case "darwin":
	case "linux":
		JSON_FILE = dir+"/servers.json"
		CONFIG_FILE = dir+"/vpngo_config.json"
		fmt.Println("Linux")
	default:
		JSON_FILE = dir+"/servers.json"
		CONFIG_FILE = dir+"/vpngo_config.json"
	}
	if len(os.Args)==3 {
		config.ServerId,err = strconv.ParseInt( os.Args[1],10, 64)
		if err!=nil{
			vpngoUsage()
			panic(err)
		}
		config.Location = os.Args[2]
		file, _ := json.MarshalIndent(config, "", " ")
		_ = ioutil.WriteFile(CONFIG_FILE, file, 0644)
	}else{
		jsonFile, err := os.Open(CONFIG_FILE)
		if err != nil {
			vpngoUsage()
			panic(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		err=json.Unmarshal(byteValue, &config)
		if err!=nil {
			vpngoUsage()
			panic(err)
		}
	}
	
	
	fmt.Println("Server Started")
	app := iris.Default()
	validate = validator.New()

	tmpl := iris.HTML("view", ".html").
		Layout("shared/layout.html").
		Reload(true)
	app.RegisterView(tmpl)
	app.HandleDir("/public", "view/public")
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().
			GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.View("shared/error.html")
	})
	WebRoutes(app)
//	app.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))
	
	
	
	
	validate = validator.New()
	app := iris.New()
	ApiRoutes(app)
	getOnlineUsers()
	go backgroundTask()
	go udpServer()
	app.Run(iris.Addr(":8081"), iris.WithoutServerError(iris.ErrServerClosed))
}
func backgroundTask() {
	getServersFiles()
	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		getServersFiles()
		fmt.Println("Getting servers lists every 1 minutes")
	}
}

func getServersFiles() {
	updateServerJsonFile(getThisServer())
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var serversArray []Server
	json.Unmarshal(byteValue, &serversArray)
	ser := getThisServer()
	for i := 0; i < len(serversArray); i++ {
		if ser.Ip != serversArray[i].Ip {
			fmt.Println("getting address " + serversArray[i].Ip)
			getAndSendInformation(serversArray[i].Ip)
		}
	}
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
func udpServer(){
	PORT := ":12345"
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	buffer := make([]byte, 1024)
	for {
		_, addr, err := connection.ReadFromUDP(buffer)
		//fmt.Print("-> ", string(buffer[0:n-1]))

		//if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
		//	fmt.Println("Exiting UDP server!")
		//	return
		//}
		//strconv.Itoa(random(1, 1001))
		data := []byte("CONNECTED")
		//fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}