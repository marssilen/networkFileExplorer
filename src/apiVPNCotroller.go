package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	_ "strings"
	"time"
)

type apiVPNController struct {
	PAGINATION_LIMIT uint16
}
type Server struct {
	Id          int64  `json:"id" validate:"required"`
	Ip          string `json:"ip" validate:"required"`
	Location    string `json:"location" validate:"required"`
	Status      string `json:"status" validate:"required"`
	OnlineUsers string `json:"online_users" validate:"required"`
	LastUpdate  int64  `json:"last_update" validate:"required"`
	IdServer    string `json:"id_server"`
	Port        string `json:"port"`
}
type ServerJsonResponse struct {
	Code    int64    `json:"code" validate:"required"`
	Data    []Server `json:"data" validate:"required"`
	Message string   `json:"message" validate:"required"`
}

//const TENMINUTES = 600000
const TENMINUTES = 10

var myClient = &http.Client{Timeout: 10 * time.Second}

func (self *apiVPNController) servers(ctx iris.Context) {
	//type Server struct {
	//	Id 			int64				`json:"id"`
	//	Ip 			string				`json:"ip"`
	//	Location 	string				`json:"location"`
	//	Id_Server 	string				`json:"id_server"`
	//	Port		string				`json:"port"`
	//}
	var serversArray []Server
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serversArray)

	msg(ctx, 0, serversArray, "servers")
}

func (self *apiVPNController) status(ctx iris.Context) {
	///////////////////////// Update Sender Server File
	var arg Server
	if err := ctx.ReadForm(&arg); err != nil {
		invalidMsg(ctx)
		return
	}
	if err := validate.Struct(arg); err != nil {
		invalidMsg(ctx)
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			fmt.Println(e)
		}
		return
	}
	updateServerJsonFile(arg)
	//////////////////////// END Update Sender Server File
	//////////////////////// Update Self Server File
	updateServerJsonFile(getThisServer())
	//////////////////////// END Update Self Server File
	//////////////////////// Send Server List File
	var serversArray []Server
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serversArray)
	//ctx.JSON(serversArray)
	msg(ctx, 0, serversArray, "done")
	//////////////////////// END Send Server List File
}

func getThisServer() Server {
	//var onlineUsers string
	//userOS := runtime.GOOS
	//switch userOS {
	//case "windows":
	//	//TODO CHANGE FOR WINDOWS, CURRENT NUMBER IS FAKE
	//	onlineUsers = fmt.Sprintf("%d", 52)
	//	break
	//default:
	//	onlineUsers = getOnlineUsers()
	//}
	onlineUsers := getOnlineUsers()
	return Server{config.ServerId, getOutboundIP(), config.Location, "1400", onlineUsers, getTimestampInMinutes(),
		"1400","54371"}
}

func getTimestampInMinutes() int64 {
	return time.Now().UnixNano() / int64(time.Minute)
	//return time.Now().UnixNano() / int64(time.Millisecond)
}
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
func getOnlineUsers() string {
	//ec Server status | grep "local sessions*" | grep -Eo "32m*[0-9]{1,}" | grep -Eo "m[0-9]{1,}" | grep -Eo "[0-9]{1,}"
	terminalCommand := "ec Server status | grep \"local sessions*\" | grep -Eo \"32m*[0-9]{1,}\" | grep -Eo \"m[0-9]{1,}\" | grep -Eo \"[0-9]{1,}\""
	result, err := exec.Command("bash", "-c", terminalCommand).Output()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
		return ""
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(result))
	// Base64 Standard Decoding
	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		panic(err)
		return ""
	}
	fmt.Println("online users"+ strings.TrimSpace(string(sDec)))
	return strings.TrimSpace(string(sDec))
}

func updateServerJsonFile(newServer Server) {
	currentTime := getTimestampInMinutes()
	if newServer.LastUpdate <= currentTime && newServer.LastUpdate >= (currentTime-TENMINUTES) {
		jsonFile, err := os.Open(JSON_FILE)
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var serversArray []Server
		json.Unmarshal(byteValue, &serversArray)
		for i := 0; i < len(serversArray); i++ {
			if serversArray[i].Ip == newServer.Ip {
				fmt.Println("Server Json File Updated")
				if serversArray[i].LastUpdate < newServer.LastUpdate {
					serversArray[i] = newServer
					saveServerJsonFile(serversArray)
				}
				return
			}
		}
		serversArray = append(serversArray, newServer)
		fmt.Println("Server Json File Appended")
		saveServerJsonFile(serversArray)
	} else {
		//remove server here
		//removeServerJsonFile(newServer)
		fmt.Println("expired server "+fmt.Sprintf("Ip is: %s \n LastUpdate:%d",newServer.Ip,newServer.LastUpdate))
		fmt.Println("expired server "+fmt.Sprintf("current: %d \n TENMINUTES:%d",currentTime,currentTime-TENMINUTES))
		fmt.Println("expired server "+fmt.Sprintf("newServer.LastUpdate <= currentTime: %d \n newServer.LastUpdate >= (currentTime-TENMINUTES):%d",
			newServer.LastUpdate <= currentTime,
			newServer.LastUpdate >= (currentTime-TENMINUTES)))
	}
}
func RemoveIndex(s []Server, index int) []Server {
	return append(s[:index], s[index+1:]...)
}
func removeServerJsonFile(newServer Server) {
	jsonFile, err := os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var serversArray []Server
	json.Unmarshal(byteValue, &serversArray)
	for i := 0; i < len(serversArray); i++ {
		if serversArray[i].Ip == newServer.Ip {
			serversArray = RemoveIndex(serversArray, i)
			fmt.Println("Server Json File Removed")
			saveServerJsonFile(serversArray)
			return
		}
	}

}
func saveServerJsonFile(serversArray []Server) {
	file, _ := json.MarshalIndent(serversArray, "", " ")
	_ = ioutil.WriteFile(JSON_FILE, file, 0644)
	fmt.Println("Server Json File Saved")
}

func getJson(serverUrl string, target interface{}) error {
	//send self
	ser := getThisServer()
	response, err := myClient.PostForm(serverUrl,
		url.Values{
			"id":          {fmt.Sprintf("%d", ser.Id)},
			"ip":          {ser.Ip},
			"location":    {ser.Location},
			"Status":      {ser.Status},
			"onlineUsers": {ser.OnlineUsers},
			"LastUpdate":  {fmt.Sprintf("%d", ser.LastUpdate)},
			"IdServer":	   {ser.IdServer},
			"Port": 	   {ser.Port},
		})

	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}
func getAndSendInformation(targetIp string) {
	//servers :=  new([]Server)//&[]Server{}
	var serverJsonResponse ServerJsonResponse
	err2 := getJson("http://"+targetIp+":8081/api/v1/status", &serverJsonResponse)
	if err2 != nil {
		removeServerJsonFile(Server{
			Id: 0,
			Ip: targetIp,
			Location:    "",
			Status:      "",
			OnlineUsers: "",
			LastUpdate:  0,
		})
		//panic(err2)
	}
	for i := 0; i < len(serverJsonResponse.Data); i++ {
		fmt.Println("getAndSendInformation :" + serverJsonResponse.Data[i].Ip)
		//Gets from remote server and validate if its online
		var serverJsonResponse2 ServerJsonResponse
		err2 = getJson("http://"+serverJsonResponse.Data[i].Ip+":8081/api/v1/servers", &serverJsonResponse2)
		if err2 != nil {
			fmt.Println("server not accessible")
		} else {
			updateServerJsonFile(serverJsonResponse.Data[i])
		}
	}
}
func (self *apiVPNController) getServersFiles() {
	//jsonFile, err := os.Open(JSON_FILE)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer jsonFile.Close()
	//byteValue, _ := ioutil.ReadAll(jsonFile)
	//var serversArray []Server
	//json.Unmarshal(byteValue, &serversArray)
	//ser := getThisServer()
	//for i := 0; i < len(serversArray); i++ {
	//	if ser.Ip != serversArray[i].Ip {
	//		fmt.Println("getting adress " + serversArray[i].Ip)
	//		getAndSendInformation(serversArray[i].Ip)
	//	}
	//}
	/*
	//reload
	jsonFile, err = os.Open(JSON_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ = ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &serversArray)

	msg(ctx, 0, serversArray, "servers")*/
}
