//nolint:noctx,funlen
package tests_test

import (
	"context"
	"dish_as_a_service/assembly"
	"dish_as_a_service/conf"
	"dish_as_a_service/entity"
	"time"

	"dish_as_a_service/repository"
	"testing"

	"github.com/Falokut/go-kit/telegram_bot"
	"github.com/Falokut/go-kit/test/fake"
	"github.com/Falokut/go-kit/test/telegramt"
	"github.com/google/uuid"

	"github.com/Falokut/go-kit/client/db"
	"github.com/Falokut/go-kit/test"
	"github.com/Falokut/go-kit/test/dbt"
	"github.com/stretchr/testify/suite"
	"github.com/txix-open/bgjob"
)

type UserBotSuite struct {
	suite.Suite
	test     *test.Test
	db       *db.Client
	userRepo repository.User
}

func TestUserBot(t *testing.T) {
	t.Parallel()
	suite.Run(t, &UserBotSuite{})
}

func (t *UserBotSuite) SetupTest() {
	test, _ := test.New(t.T())
	t.test = test
	db := dbt.New(test, db.WithMigrationRunner("../migrations", test.Logger()))
	t.userRepo = repository.NewUser(db.Client)
	t.db = db.Client
}

func (t *UserBotSuite) initBot() (*telegramt.BotServerMock, *telegram_bot.BotAPI) {
	bgjobDb := bgjob.NewPgStore(t.db.DB.DB)
	bgjobCli := bgjob.NewClient(bgjobDb)
	tgBot, botMockServerMock := telegramt.TestBot(t.test)
	locatorCfg, err := assembly.Locator(
		context.Background(),
		t.test.Logger(),
		t.db,
		nil,
		tgBot,
		bgjobCli,
		conf.LocalConfig{App: conf.App{AdminSecret: "secret"}},
	)
	t.Require().NoError(err)
	tgBot.Upgrade(locatorCfg.BotRouter)
	go tgBot.Serve(context.Background(), telegram_bot.UpdatesConfig{Timeout: 1}) // nolint:errcheck
	return botMockServerMock, tgBot
}

func (t *UserBotSuite) Test_RegisterUser_HappyPath() {
	botServerMock, bot := t.initBot()
	defer bot.StopReceivingUpdates()
	username := fake.It[string]()
	firstName := fake.It[string]()
	telegramId := fake.It[int64]()
	updates := []telegram_bot.Update{
		{
			UpdateId: 1,
			Message: &telegram_bot.Message{
				Entities: []telegram_bot.MessageEntity{
					{
						Type:   "bot_command",
						Length: len("/start"),
					},
				},
				Text: "/start",
				Chat: &telegram_bot.Chat{
					Id: telegramId,
				},
				From: &telegram_bot.User{
					UserName:  username,
					FirstName: firstName,
					Id:        telegramId,
				},
			},
		},
	}
	botServerMock.SetUpdates(updates)
	time.Sleep(time.Second)
	userId, err := t.userRepo.GetUserIdByTelegramId(context.Background(), telegramId)
	t.Require().NoError(err)

	user, err := t.userRepo.GetUserInfo(context.Background(), userId)
	t.Require().NoError(err)
	t.Require().Equal(username, user.Username)
	t.Require().Equal(firstName, user.Name)
	t.Require().False(user.Admin)
}

func (t *UserBotSuite) Test_PassBySecret_HappyPath() {
	botServerMock, bot := t.initBot()
	defer bot.StopReceivingUpdates()
	username := fake.It[string]()
	firstName := fake.It[string]()
	telegramId := fake.It[int64]()
	updates := []telegram_bot.Update{
		{
			UpdateId: 1,
			Message: &telegram_bot.Message{
				Entities: []telegram_bot.MessageEntity{
					{
						Type:   "bot_command",
						Length: len("/start"),
					},
				},
				Text: "/start",
				Chat: &telegram_bot.Chat{
					Id: telegramId,
				},
				From: &telegram_bot.User{
					UserName:  username,
					FirstName: firstName,
					Id:        telegramId,
				},
			},
		},
		{
			UpdateId: 2,
			Message: &telegram_bot.Message{
				Entities: []telegram_bot.MessageEntity{
					{
						Type:   "bot_command",
						Length: len("/pass_by_secret"),
					},
				},
				Text: "/pass_by_secret secret",
				Chat: &telegram_bot.Chat{
					Id: telegramId,
				},
				From: &telegram_bot.User{
					UserName:  username,
					FirstName: firstName,
					Id:        telegramId,
				},
			},
		},
	}
	botServerMock.SetUpdates(updates)
	time.Sleep(time.Second)
	userId, err := t.userRepo.GetUserIdByTelegramId(context.Background(), telegramId)
	t.Require().NoError(err)

	user, err := t.userRepo.GetUserInfo(context.Background(), userId)
	t.Require().NoError(err)
	t.Require().Equal(username, user.Username)
	t.Require().Equal(firstName, user.Name)
	t.Require().True(user.Admin)
}

func (t *UserBotSuite) Test_RemoveAdminByUsername_HappyPath() {
	firstUserId := uuid.NewString()
	firstChatId := fake.It[int64]()
	firstUsername := fake.It[string]()
	err := t.userRepo.Register(context.Background(), entity.RegisterUser{
		Id:       firstUserId,
		Username: firstUsername,
		Name:     fake.It[string](),
		Telegram: &entity.Telegram{
			ChatId: firstChatId,
			UserId: firstChatId,
		},
	})
	t.Require().NoError(err)
	err = t.userRepo.SetUserAdminStatus(context.Background(), firstUsername, true)
	t.Require().NoError(err)

	secondUserId := uuid.NewString()
	secondChatId := fake.It[int64]()
	secondUsername := fake.It[string]()
	err = t.userRepo.Register(context.Background(), entity.RegisterUser{
		Id:       secondUserId,
		Username: secondUsername,
		Name:     fake.It[string](),
		Telegram: &entity.Telegram{
			ChatId: secondChatId,
			UserId: secondChatId,
		},
	})
	t.Require().NoError(err)
	err = t.userRepo.SetUserAdminStatus(context.Background(), secondUsername, true)
	t.Require().NoError(err)

	botServerMock, bot := t.initBot()
	botServerMock.SetUpdates([]telegram_bot.Update{
		{
			UpdateId: 1,
			Message: &telegram_bot.Message{
				Entities: []telegram_bot.MessageEntity{
					{
						Type:   "bot_command",
						Length: len("/remove_admin"),
					},
				},
				Text: "/remove_admin " + secondUsername,
				Chat: &telegram_bot.Chat{
					Id: firstChatId,
				},
				From: &telegram_bot.User{
					UserName: firstUsername,
					Id:       firstChatId,
				},
			},
		},
	})
	time.Sleep(time.Second)
	bot.StopReceivingUpdates()
	user, err := t.userRepo.GetUserInfo(context.Background(), secondUserId)
	t.Require().NoError(err)
	t.Require().False(user.Admin)
}
