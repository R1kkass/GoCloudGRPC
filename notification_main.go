package main

import (
	"encoding/json"
	"fmt"
	"log"
	db "mypackages/db"
	"mypackages/helpers"
	"strconv"

	"mypackages/proto/notification"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type notificationServer struct {
	notification.UnimplementedNotificationGreeterServer
}


func (s *notificationServer) GetNotification(in *notification.Empty, responseStream notification.NotificationGreeter_GetNotificationServer) error {
	ctx := responseStream.Context()
	user, err := helpers.GetUserFormMd(ctx)
	channel := make(chan map[string]any)

	defer func(){
		if r:=recover(); r != nil {
			fmt.Println("Error get notification: ", r)
		}
	}()

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	go func() {

			key := strconv.Itoa(int(user.ID)) + "_notification"
			res := db.ConnectRedisNotificationDB.Subscribe(ctx, key)

			for {
				var jsonDecodeMsg map[string]any
				message, err := res.ReceiveMessage(ctx)
				json.Unmarshal([]byte(message.Payload), &jsonDecodeMsg)

				if err != nil {
					log.Println("Can not create subscribe")
					return
				}
			
				channel<-jsonDecodeMsg
			}

	}()

	for {
		val:=<-channel

		var options = make(map[string]string)
		
		if _, ok := val["options"]; ok {
			for k, v := range val["options"].(map[string]any) {
				options[k] = fmt.Sprintf("%v", v)
			}
		}

		if _, ok := val["description"]; !ok {
			val["description"] = ""
		}

		err := responseStream.Send(
			&notification.NotificationMessage{
				Type:          val["type"].(string),
				Description: val["description"].(string),
				Title:       val["title"].(string),
				Options: options,
			},
		)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}

}
