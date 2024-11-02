package notification_action

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/db"
)

func NotificationSendToRedis(userId int, ctx context.Context, channel *chan map[string]any) {
	key := strconv.Itoa(userId) + "_notification"
	res := db.ConnectRedisNotificationDB.Subscribe(ctx, key)

	for {
		var jsonDecodeMsg map[string]any
		message, err := res.ReceiveMessage(ctx)
		json.Unmarshal([]byte(message.Payload), &jsonDecodeMsg)

		if err != nil {
			log.Println("Can not create subscribe")
			return
		}

		*channel <- jsonDecodeMsg
	}

}
