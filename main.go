package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type UserDetail struct {
	Username string `json:"username"`
	Follower int    `json:"followers"`
}

type UserFollower struct {
	Follower int `json:"followers"`
}

type ResponseError struct {
	Error string `json:"error"`
}

type UserID map[string]UserDetail

func main() {
	router := httprouter.New()
	router.GET("/follower/:username", FindByUsername)
	router.GET("/user/:userId/detail", FindByUserId)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func ConsumeApi() UserID {
	//get data json
	response, err := http.Get("https://jsonkeeper.com/b/DMXK")
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	//read all data json and string
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	//initial map
	var userid UserID

	//decode data json string
	dec := json.NewDecoder(strings.NewReader(string(responseData)))

	//decode assign to map
	for {
		if err := dec.Decode(&userid); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	return userid
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	if err != nil {
		panic(err)
	}
}

func FindByUserId(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userId := params.ByName("userId")

	map_data := ConsumeApi()
	//checking map
	if v, found := map_data[userId]; found {
		responseApi := UserDetail{
			Username: v.Username,
			Follower: v.Follower,
		}
		WriteToResponseBody(writer, responseApi)

	} else {
		responseError := ResponseError{
			Error: "Not Found",
		}
		WriteToResponseBody(writer, responseError)
	}
}

func FindByUsername(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	username := params.ByName("username")

	map_data := ConsumeApi()

	//find followers
	var follower_user int
	for _, v := range map_data {
		if v.Username == username {
			follower_user = v.Follower
		}
	}

	if follower_user != 0 {
		responseApi := UserFollower{
			Follower: follower_user,
		}

		WriteToResponseBody(writer, responseApi)

	} else {
		responseError := ResponseError{
			Error: "Not Found",
		}
		WriteToResponseBody(writer, responseError)
	}

}
