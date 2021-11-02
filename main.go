package main

import (
	csv1 "encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
)


type postResponse struct {
	UserId            int64            `json:"userId"`
	Id                int64            `json:"id"`
	Title           string            `json:"title"`
	Body           string           `json:"body"`
}

type commentResponse struct {
	PostId            int64            `json:"postId"`
	Id                int64            `json:"id"`
	Name           string            ` json:"name"`
	Email            string            `json:"email"`
	Body           string           `json:"body"`
}

type completePostObj struct {
	UserId            int            `json:"userId"`
	Id                int            `json:"id"`
	Title           string            `json:"title"`
	Body           string           `json:"body"`
}

var Result []completePostObj

func main() {

	// getPost Data
	postDataList := createPostRequest()

	// getCommentData

	getDataList := createGetRequest()

	for _, pV := range postDataList{
		count := 0
		for _, cV := range getDataList {
			if pV.Id == cV.PostId {
					if count == 0 {
						c := completePostObj{}
						c.UserId = int(pV.UserId)
						c.Id = int(pV.Id)
						c.Title = pV.Title
						c.Body = pV.Body + "|" + cV.Body
						Result = append(Result, c)
					} else {
						for _, reV := range Result {
							if int(pV.Id) == reV.Id {
								pV.Body = pV.Body + "|" + cV.Body
							}
						}
					}
					count ++

			}
		}
	}
	// writeToCSV
	writeToCsv(Result)

	// return requests of post of all comments
	handleRequests()


}

func createPostRequest() []postResponse{
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	 body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var postArray []postResponse

	if err = json.Unmarshal(body, &postArray); err != nil {
		fmt.Printf("error %s", err)
	}

	return postArray

}

func createGetRequest() []commentResponse{
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var commentArray []commentResponse
	if err = json.Unmarshal(body, &commentArray); err != nil {
		fmt.Printf("error %s", err)
	}

	return commentArray

}

func writeToCsv(result []completePostObj) {
	csvFile, err := os.Create("./source.csv")
	if err != nil {
		fmt.Print("not able to open the file")
	}
	//
	//
	csvwriter := csv1.NewWriter(csvFile)

	for i:= -1; i < len(result); i++ {
		var row []string
		if i == -1 {
			row = append(row, "userId")
			row = append(row, "id")
			row = append(row, "title")
			row = append(row, "body")
			csvwriter.Write(row)
			continue
		}
		row = append(row, strconv.Itoa(result[i].UserId))
		row = append(row, strconv.Itoa(result[i].Id))
		row = append(row, result[i].Title)
		row = append(row, result[i].Body)
		csvwriter.Write(row)
	}

	csvwriter.Flush()
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllPostComments)
	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func returnAllPostComments(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllPostComments")
	json.NewEncoder(w).Encode(Result)
}
