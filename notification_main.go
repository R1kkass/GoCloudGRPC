package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mypackages/db"
	"mypackages/helpers"
	"mypackages/proto/notification"
	"strconv"
)

type notificationServer struct {
	notification.UnimplementedNotificationGreeterServer
}

func (s *notificationServer) GetNotification(in *notification.Empty, responseStream notification.NotificationGreeter_GetNotificationServer) error {
	ctx := responseStream.Context()
	user, _ := helpers.GetUserFormMd(ctx)
	channel := make(chan map[string]string)

	go func() {
		for {

			key := strconv.Itoa(int(user.ID)) + ":" + "*"
			keys, _ := db.ConnectRedisNotificationDB.Keys(ctx, key).Result()
			db.ConnectRedisNotificationDB.Watch(ctx, )
			for _, v := range keys {
				objects, _ := db.ConnectRedisNotificationDB.LRange(ctx, v, 0, -1).Result()
				for _, o := range objects {
					fmt.Println(o)
					var mapMessage map[string]string
					json.Unmarshal([]byte(o), &mapMessage)
					channel <- mapMessage
				}
				db.ConnectRedisNotificationDB.Del(ctx, v)
			}
		}
	}()

	for {
		val := <-channel
		err := responseStream.Send(
			&notification.NotificationMessage{
				Type:          val["type"],
				Description: val["description"],
				Title:       val["title"],
			},
		)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}

	return nil
}
