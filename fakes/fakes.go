package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/manveru/faker"
)

type envelope struct {
	Data    interface{}
	Message string
}

type user struct {
	ID       string
	Username string `json:"username"`
	Password string
	Email    string
	Address  string
}
type stock struct {
	ID       string
	Name     string
	Quantity int
}

var BASEURL string
var TOKEN string

func main() {
	args := os.Args[1:]
	BASEURL = args[0]
	TOKEN = args[1]
	userArray := addFakeUser(13)
	addFakeStock(5, userArray)
}

func addFakeUser(count int) []*user {
	var returnArray []*user
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	for i := 0; i < count; i++ {
		newUser := &user{Username: fake.UserName(), Password: fake.SafeEmail(), Email: fake.Email(), Address: fake.StreetAddress()}
		newUserJson, err := json.Marshal(newUser)
		if err != nil {
			fmt.Println(err)
			continue
		}
		responseByte := postRequest(BASEURL+"users", newUserJson)
		env := envelope{}
		err = json.Unmarshal(responseByte, &env)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonString, _ := json.Marshal(env.Data)
		s := user{}
		json.Unmarshal(jsonString, &s)
		s.Password = newUser.Password
		returnArray = append(returnArray, &s)
	}
	return returnArray
}
func addFakeStock(count int, userArray []*user) []*stock {
	var returnArray []*stock
	fake, err := faker.New("en")
	if err != nil {
		panic(err)
	}
	for i := 0; i < count; i++ {
		newStock := &stock{Name: fake.CompanyName(), Quantity: fake.Rand.Intn(100)}
		newStockJson, err := json.Marshal(newStock)
		if err != nil {
			fmt.Println(err)
			continue
		}
		responseByte := postRequest(BASEURL+"stocks/users/"+userArray[i].ID, newStockJson)
		env := envelope{}
		err = json.Unmarshal(responseByte, &env)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonString, _ := json.Marshal(env.Data)
		s := stock{}
		json.Unmarshal(jsonString, &s)
		returnArray = append(returnArray, &s)
	}
	return returnArray
}

func postRequest(url string, jsonStr []byte) []byte {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Bearer " + TOKEN)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return body
}

func beautifyJson(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}
