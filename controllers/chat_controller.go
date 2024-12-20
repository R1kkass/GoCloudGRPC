package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	chat_actions "github.com/R1kkass/GoCloudGRPC/actions/chat"
	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/proto/chat"
	"github.com/R1kkass/GoCloudGRPC/structs"
	"github.com/R1kkass/GoCloudGRPC/validate"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ChatServer struct {
	chat.UnimplementedChatGreeterServer
	Conns map[string]structs.DataStreamConnect
}

func (s *ChatServer) CreateChat(ctx context.Context, in *chat.CreateRequestChat) (*chat.CreateResponseChat, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error CreateChat: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"otherId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetOtherId(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	if err := chat_actions.CheckChatExist(ctx, in); err != nil {
		return nil, status.Error(codes.AlreadyExists, "Чат уже существует")
	}

	p, g, modelChat, err := chat_actions.CreateChatTransaction(user.ID, uint(in.GetOtherId()))

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	objectMessage := map[string]any{
		"title": "Новый запрос на отправку сообщений",
		"type":  "New_ChatRequest",
	}

	go chat_actions.NotificationChatCreate(int(in.GetOtherId()), objectMessage)

	return &chat.CreateResponseChat{
		Message: "Чат создан",
		Keys: &chat.Keys{
			P: p,
			G: g,
		},
		ChatId: uint32(modelChat.ID),
	}, nil
}

func (s *ChatServer) StreamGetChat(in *chat.Empty, requestStream chat.ChatGreeter_StreamGetChatServer) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error StreamGetChat: ", r)
		}
	}()

	var chats []*chat.ChatUsersCount

	ctx := requestStream.Context()
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return status.Error(codes.Unauthenticated, "Токен не найден")
	}

	jwtToken, ok := md["authorization"]

	if !ok {
		return status.Error(codes.Unauthenticated, "Токен не найден")
	}

	user, err := helpers.GetUser(jwtToken)
	channel := make(chan map[string]any)

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
		key := strconv.Itoa(int(user.ID)) + "_messages"

		for {
			res := db.ConnectRedisNotificationDB.Subscribe(ctx, key)
			var jsonDecodeMsg map[string]any
			message, err := res.ReceiveMessage(ctx)
			json.Unmarshal([]byte(message.Payload), &jsonDecodeMsg)

			if err != nil {
				log.Println("can't get chats")
				return
			}

			channel <- jsonDecodeMsg
		}

	}()

	chats, err = chat_actions.StreamGetChat(user.ID)

	if err != nil {
		return status.Error(codes.Unauthenticated, "Невозможно получить сообщения")
	}

	err = requestStream.Send(&chat.StreamGetResponseChat{
		Chats: chats,
	})
	if err != nil {
		log.Println("error while sending chats:", err)
		return err
	}

	for {
		<-channel

		var chats []*chat.ChatUsersCount

		chats, err = chat_actions.StreamGetChat(user.ID)

		if err != nil {
			return status.Error(codes.Unauthenticated, "Невозможно получить сообщения")
		}

		err := requestStream.Send(&chat.StreamGetResponseChat{
			Chats: chats,
		})
		if err != nil {
			log.Println("error while sending chats:", err)
			return err
		}

	}
}

func (s *ChatServer) GetUnSuccessChats(ctx context.Context, in *chat.Empty) (*chat.GetUnSuccessChatsResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error GetUnSuccessChats: ", r)
		}
	}()

	var chats []*chat.ChatUsers
	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	r := db.DB.Model(&Model.ChatUser{}).
		Preload("User").Preload("Chat").Preload("Chat.ChatUsers.User").
		Where("chat_users.user_id = ? AND submit_create = ?", user.ID, Model.UnSuccessChat).
		Find(&chats)
	if r.Error != nil {
		return nil, status.Error(codes.Unknown, "Пользователь не найден")
	}
	return &chat.GetUnSuccessChatsResponse{
		Chats: chats,
	}, nil
}

func (s *ChatServer) CreateSecondaryKey(ctx context.Context, in *chat.CreateSecondaryKeyRequest) (*chat.CreateSecondaryKeyResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error CreateSecondaryKey: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
			"key": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetKey(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}
	_, err = chat_actions.CheckSecondaryKey(user.ID, in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.AlreadyExists, "Ключ уже создан")
	}

	db.DB.Create(&Model.KeysSecondary{
		UserID: user.ID,
		ChatID: uint(in.GetChatId()),
		Key:    in.GetKey(),
	})

	return &chat.CreateSecondaryKeyResponse{
		Message: "Успешно",
	}, nil
}

func (s *ChatServer) GetSecondaryKey(ctx context.Context, in *chat.GetSecondaryKeyRequest) (*chat.GetSecondaryKeyResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error GetSecondaryKey: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}

	err, _ = chat_actions.CheckChat(in.GetChatId(), user.ID)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	key, err := chat_actions.GetSecondaryKey(user.ID, in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	keys, err := chat_actions.GetPublicKey(in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &chat.GetSecondaryKeyResponse{
		Key: key.Key,
		P:   keys.P,
	}, nil
}

func (s *ChatServer) GetPublicKey(ctx context.Context, in *chat.GetPublicKeyRequest) (*chat.GetPublicKeyResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error GetPublicKey: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	keys, err := chat_actions.GetPublicKey(in.GetChatId())

	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &chat.GetPublicKeyResponse{
		G: keys.G,
		P: keys.P,
	}, nil
}

func (s *ChatServer) AcceptChat(ctx context.Context, in *chat.AcceptChatRequest) (*chat.AcceptChatResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error AcceptChat: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	db.DB.Model(&Model.ChatUser{}).Where("chat_id = ?", in.GetChatId()).Update("submit_create", Model.CreatedChat)

	return &chat.AcceptChatResponse{
		Message: "Успех",
	}, nil
}

func (s *ChatServer) DissalowChat(ctx context.Context, in *chat.DissalowChatRequest) (*chat.DissalowChatResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error DissalowChat: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	var chatUser Model.ChatUser
	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("error ", r)
		}
	}()
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.NotFound, "Пользователь не найден")
	}

	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", in.GetChatId(), user.ID).First(&chatUser)

	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Error(codes.Aborted, "Чат не найден")
	}

	err = chat_actions.DissalowChatTransaction(chatUser.ChatID)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Не удалось отказать в доступе")
	}

	return &chat.DissalowChatResponse{
		Message: "Успех",
	}, nil
}

func (s *ChatServer) GetMessages(ctx context.Context, in *chat.GetMessagesRequest) (*chat.GetMessagesResponse, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get messages: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
			"page": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetPage(),
			},
		},
	)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	var messages []*chat.Message

	if err := helpers.CheckChat(ctx, in.GetChatId()); err != nil {
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Пользователь не найден")
	}
	var page = int64(in.GetPage())
	var countNotReaded int64 = 0
	countNotReaded, err = chat_actions.GetCountNotReadedMessages(int(in.GetChatId()), int(user.ID))

	if err != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	if in.GetInit() {
		page = countNotReaded / 10
		if page != 0 {
			page += 1
		}
	}

	r := db.DB.Model(&Model.Message{}).Preload("User").
		Select("messages.*, SUM(CASE WHEN un_readed_messages.id IS NULL THEN 0 ELSE 1 END) AS un_readed_message").
		Joins("LEFT JOIN un_readed_messages ON un_readed_messages.message_id = messages.id AND un_readed_messages.user_id = ?", user.ID).
		Preload("ChatFiles").
		Where("messages.chat_id = ? AND status_message = 'success'", in.GetChatId()).
		Group("messages.id").
		Order("id DESC").Offset(10 * int(page)).
		Limit(10).
		Find(&messages)

	if r.Error != nil {
		return nil, status.Error(codes.Unknown, "Не удалось получить сообщения")
	}

	return &chat.GetMessagesResponse{
		Messages:     messages,
		Page:         int32(page),
		CountNotRead: int32(countNotReaded),
	}, nil
}

func (s *ChatServer) UploadChatFile(requestStream chat.ChatGreeter_UploadChatFileServer) error {
	ctx := requestStream.Context()
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	var completedParts []types.CompletedPart
	var resp *s3.CreateMultipartUploadOutput
	partNumber := 1
	var messageId uint
	var chatFile *Model.ChatFile
	var size int
	defer func() {
		if r := recover(); r != nil {
			chat_actions.Rollback(messageId)
		}
	}()

	for {
		req, err := requestStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error UploadChatFile: ", err)
			chat_actions.Rollback(messageId)
			return err
		}

		err = validate.Valid(
			validate.ValidType{
				"fileName": validate.ValidateStruct{
					Rule:  "required|string",
					Value: req.GetFileName(),
				},
				"messageId": validate.ValidateStruct{
					Rule:  "required|uint32",
					Value: req.GetMessageId(),
				},
				"chunk": validate.ValidateStruct{
					Rule:  "required",
					Value: len(req.GetChunk()),
				},
			},
		)

		if err != nil {
			fmt.Println("Error UploadChatFile: ", err)
			chat_actions.Rollback(messageId)
			return status.Error(codes.Unknown, "Неизвестная ошибка")
		}

		chatUser, err := chat_actions.CheckChatByMessageId(req.GetMessageId(), user.ID)

		if err != nil {
			fmt.Println("Error UploadChatFile: ", err)
			chat_actions.Rollback(messageId)
			return status.Error(codes.Unknown, "Чат не найден")
		}

		if resp == nil {
			messageId = uint(req.GetMessageId())
			err := chat_actions.GetCountChatFile(messageId)
			if err != nil {
				fmt.Println("Error UploadChatFile: ", err)
				chat_actions.Rollback(messageId)
				return status.Error(codes.OutOfRange, "Переполнение файлов")
			}

			chatFile, err = chat_actions.UploadChatFileTransaction(user.ID, chatUser.ChatID, messageId, req.GetFileName())

			if err != nil {
				fmt.Println("Error UploadChatFile: ", err)
				chat_actions.Rollback(messageId)
				return status.Error(codes.Unknown, "Неизвестная ошибка")
			}

			resp, err = chat_actions.UploadChatFile(ctx, chatFile.ID)

			if err != nil {
				fmt.Println(err)
				chat_actions.Rollback(messageId)
				return status.Error(codes.Unknown, "Неизвестная ошибка")
			}
		}
		size += len(req.GetChunk())
		completedPart, err := chat_actions.UploadPart(ctx, resp, req.GetChunk(), partNumber)
		if err != nil {
			chat_actions.AbortMultipartUpload(ctx, resp)
			fmt.Println("Error UploadChatFile:", err)
			chat_actions.Rollback(messageId)
			return status.Error(codes.Unknown, "Неизвестная ошибка")
		}
		completedParts = append(completedParts, *completedPart)
		partNumber++
	}
	_, err = chat_actions.CompleteMultipartUpload(ctx, resp, completedParts)
	if err != nil {
		fmt.Println("Error UploadChatFile:", err)
		chat_actions.Rollback(messageId)
		return status.Error(codes.Unknown, "Не удалось загрузить файл")
	}

	r := db.DB.Model(&Model.ChatFile{}).Where("id = ?", chatFile.ID).Update("size", size)

	if r.Error != nil {
		fmt.Println("Error UploadChatFile:", err)
		chat_actions.Rollback(messageId)
		return status.Error(codes.Unknown, "Не удалось загрузить файл")
	}

	return requestStream.SendAndClose(&chat.UploadFileChatResponse{Message: "Успешно загружено"})
}

func (s *ChatServer) StreamGetMessagesGeneral(in *chat.Empty, responseStream chat.ChatGreeter_StreamGetMessagesGeneralServer) error {
	ctx := responseStream.Context()
	user, err := helpers.GetUserFormMd(ctx)
	channel := make(chan bool)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error get general messages: ", r)
		}
	}()

	if err != nil {
		return status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	var count int64
	r := db.DB.Model(&Model.UnReadedMessage{}).Where("user_id = ?", user.ID).Distinct("chat_id").Count(&count)

	if r.Error != nil {
		log.Println("error while sending count messages:", err)
		return err
	}

	err = responseStream.Send(&chat.StreamGetMessagesGeneralResponse{
		Count: int32(count),
	})

	if err != nil {
		log.Println("error while sending count messages:", err)
		return err
	}

	go chat_actions.NotificationObserver(ctx, int(user.ID), &channel)

	for {
		<-channel
		var count int64
		r := db.DB.Model(&Model.UnReadedMessage{}).Where("user_id = ?", user.ID).Distinct("chat_id").Count(&count)
		if r.Error != nil {
			log.Println("error while sending count messages:", err)
			return err
		}
		err := responseStream.Send(&chat.StreamGetMessagesGeneralResponse{
			Count: int32(count),
		})

		if err != nil {
			log.Println("error while sending count messages:", err)
			return err
		}
	}
}

func (s *ChatServer) StreamGetMessages(stream chat.ChatGreeter_StreamGetMessagesServer) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic error: ", r)
		}
	}()

	ctx := stream.Context()
	md, _ := metadata.FromIncomingContext(ctx)
	user, err := helpers.GetUserFormMd(ctx)
	channel := make(chan *chat.StreamGetMessagesResponse)
	if err != nil {
		return status.Error(codes.PermissionDenied, "пользователь не найден")
	}
	chatId, err := strconv.Atoi(md["chat_id"][0])
	if err != nil {
		return status.Error(codes.PermissionDenied, "чат не найден ошибка")
	}
	err, chatUser := chat_actions.CheckChat(uint32(chatId), user.ID)

	if err != nil {
		return status.Error(codes.PermissionDenied, "чат не найден ошибка")
	}

	go func() {
		for {
			messageResponse := <-channel
			for _, v := range s.Conns {
				if v.ChatId == uint(chatId) {
					v.Stream.Send(messageResponse)
				}
			}
		}
	}()
	jwtToken, _ := md["authorization"]
	jwtToken = strings.Split(jwtToken[0], " ")
	s.Conns[jwtToken[1]] = structs.DataStreamConnect{
		ChatId: chatUser.ChatID,
		UserID: user.ID,
		Stream: stream,
	}

	defer chat_actions.CloseConnect(s.Conns, jwtToken[1])

	return chat_actions.StreamGetMessages(stream, s.Conns, uint(chatId), user.ID, &channel, jwtToken[1])
}

func (s *ChatServer) DownloadChatFile(in *chat.DownloadFileChatRequest, responseStream chat.ChatGreeter_DownloadChatFileServer) error {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic error: ", r)
		}
	}()

	ctx := responseStream.Context()
	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		fmt.Println("Error DownloadFileChat: ", err)
		return status.Error(codes.NotFound, "Пользователь не найден")
	}

	err = chat_actions.CheckChatFile(user, in.GetChatFileId())
	if err != nil {
		fmt.Println("Error DownloadFileChat: ", err)
		return status.Error(codes.NotFound, "Файл не найден")
	}

	size, err := chat_actions.GetFileSize(ctx, strconv.Itoa(int(in.GetChatFileId())))
	var rangeInt int64 = 0
	if err != nil {
		fmt.Println("Error DownloadFileChat: ", err)
		return status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	for *size > rangeInt {
		bytes, err := chat_actions.DownloadChunk(ctx, strconv.Itoa(int(in.GetChatFileId())), rangeInt, int(*size))
		if err != nil {
			fmt.Println("Error DownloadFileChat: ", err)
			return status.Error(codes.Unknown, "Неизвестная ошибка")
		}
		rangeInt += 256 * 1024
		responseStream.Send(&chat.DownloadFileChatResponse{
			Chunk:    bytes,
			Progress: float32(rangeInt / (*size) * 100),
		})
	}

	return nil
}

func (s *ChatServer) CreateFileMessage(ctx context.Context, in *chat.CreateFileMessageRequest) (*chat.CreateFileMessageResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic error: ", r)
		}
	}()

	err := validate.Valid(
		validate.ValidType{
			"chatId": validate.ValidateStruct{
				Rule:  "required",
				Value: in.GetChatId(),
			},
		},
	)

	if err != nil {
		fmt.Println("Error CreateFileMessage: ", err)
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	user, err := helpers.GetUserFormMd(ctx)
	if err != nil {
		fmt.Println("Error CreateFileMessage: ", err)
		return nil, status.Error(codes.NotFound, "Пользователь не найден")
	}

	err, _ = chat_actions.CheckChat(in.GetChatId(), user.ID)

	if err != nil {
		fmt.Println("Error CreateFileMessage: ", err)
		return nil, status.Error(codes.NotFound, "Чат не найден")
	}

	var message = &Model.Message{
		Text: in.GetText(),
		UserRelation: Model.UserRelation{
			UserID: user.ID,
		},
		ChatRelations: Model.ChatRelations{
			ChatID: uint(in.GetChatId()),
		},
		TypeMessage:   Model.FileMessage,
		StatusMessage: Model.Uploading,
	}

	r := db.DB.Create(&message)

	if r.Error != nil || r.RowsAffected == 0 {
		fmt.Println("Error CreateFileMessage: ", err)
		return nil, status.Error(codes.Unknown, "Неизвестная ошибка")
	}

	return &chat.CreateFileMessageResponse{
		MessageId: uint32(message.ID),
		CreatedAt: message.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: message.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
