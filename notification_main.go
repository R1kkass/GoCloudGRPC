package main

import (
	"fmt"
	"log"

	notification_action "github.com/R1kkass/GoCloudGRPC/actions/notification"
	"github.com/R1kkass/GoCloudGRPC/helpers"

	"github.com/R1kkass/GoCloudGRPC/proto/notification"

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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get notification: ", r)
		}
	}()

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	go notification_action.NotificationSendToRedis(int(user.ID), ctx, &channel)

	for {
		val := <-channel

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
				Type:        val["type"].(string),
				Description: val["description"].(string),
				Title:       val["title"].(string),
				Options:     options,
			},
		)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}

}
