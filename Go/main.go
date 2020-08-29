package main

import (
	"eduhacks2020/Go/protocol/websocket"
	_ "github.com/gorilla/websocket"
	"net/http"
)

func main() {

	//http.HandleFunc("/", helloWorld)
	//http.ListenAndServe(":8000", nil)
	websocketHandler := &protocol.Controller{}
	http.HandleFunc("/ws", websocketHandler.Run)
	if err := http.ListenAndServe(":555", nil); err != nil {
		panic(err)
	}
}

//func helloWorld(w http.ResponseWriter, r *http.Request) {
//	test := &protobuf.Student{
//		Name:   "Hao",
//		Male:   true,
//		Scores: []int32{98, 85, 88},
//	}
//	data, err := proto.Marshal(test)
//	if err != nil {
//		return
//	}
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("content-type", "application/octet-stream")
//	w.Write(data)
//}
