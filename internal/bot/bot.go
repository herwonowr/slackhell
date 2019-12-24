package bot

import (
	"context"
	"log"

	"github.com/herwonowr/slacker"
	"github.com/herwonowr/slackhell/internal/helper"
	"github.com/herwonowr/slackhell/internal/repository"
	"github.com/herwonowr/slackhell/internal/service"
	"github.com/nlopes/slack"
)

// Service ...
type Service interface {
	InitBot(DBVersion int, slackID string, slackRealName string) error
	Listen() error
}

type botService struct {
	slack   *slacker.Slacker
	service service.Service
	helper  helper.Service
}

// NewService ...
func NewService(srv service.Service, token string, debug bool) Service {
	client := slacker.NewClient(token, slacker.WithDebug(debug))

	return &botService{
		slack:   client,
		service: srv,
	}
}

func (s *botService) InitBot(DBVersion int, slackID string, slackRealName string) error {
	version, err := s.service.GetVersion()
	if err != nil && err.Error() == "version not found" {
		version = 0
	} else if err != nil {
		return err
	}

	//TODO: database migrator
	if version < DBVersion {
		log.Printf("migrating database version: %d to %d\n", version, DBVersion)

		err := s.service.PutVersion(DBVersion)
		if err != nil {
			return err
		}

		role := repository.AdminRole
		return s.service.CreateAccount(role, slackID, slackRealName)
	}

	return nil
}

func (s *botService) Listen() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.runBot()
	err := s.slack.Listen(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *botService) runBot() {
	s.slack.Init(func() {
		log.Println("Slackhell connected")
	})

	s.slack.Err(func(err string) {
		log.Println(err)
	})

	s.slack.DefaultCommand(func(request slacker.Request, response slacker.ResponseWriter) {
		s.sendReply("Hi!, need something?, try type `help` for available commands", request, response, true)
	})

	s.slack.DefaultEvent(func(event interface{}) {
		log.Printf("Slackhell event: %v\n", event)
	})

	s.cmdHelp()
	s.cmdUser()
	s.cmdUsers()
	s.cmdAddUser()
	s.cmdUpdateUser()
	s.cmdDeleteUser()
	s.cmdExecute()
	s.cmdGenerateShellcode()
	s.cmdGetShellcodes()
	s.cmdUpdateShellcode()
	s.cmdDeleteShellcode()
	s.cmdGetSession()
	s.cmdSetSession()
	s.cmdKillSession()
}

func (s *botService) authorized(account string, admin bool) bool {
	a, err := s.service.GetAccountBySlackID(account)
	if a == nil || err != nil {
		return false
	}

	if admin {
		if a.Role != repository.AdminRole {
			return false
		}
	}

	return true
}

func (s *botService) sendReply(msg string, request slacker.Request, response slacker.ResponseWriter, mention bool) {
	agent := request.Event().User
	client := response.Client()
	conInfo, err := client.GetConversationInfo(request.Event().Channel, false)
	if err != nil {
		log.Println("error get conversation info")
		return
	}

	if !conInfo.IsIM && mention {
		response.Reply("*Agent:* <@" + agent + ">" + "\n_Response only visible to agent_")
	}

	_, err = client.PostEphemeral(
		request.Event().Channel,
		request.Event().User,
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		log.Println("error sending reply message")
	}
}

func (s *botService) sendAttachmentReply(attachment []slack.Attachment, request slacker.Request, response slacker.ResponseWriter, mention bool) {
	agent := request.Event().User
	client := response.Client()
	conInfo, err := client.GetConversationInfo(request.Event().Channel, false)
	if err != nil {
		log.Println("error get conversation info")
		return
	}

	if !conInfo.IsIM && mention {
		response.Reply("*Agent:* <@" + agent + ">" + "\n_Response only visible to agent_")
	}

	_, err = client.PostEphemeral(
		request.Event().Channel,
		request.Event().User,
		slack.MsgOptionAsUser(true),
		slack.MsgOptionAttachments(attachment...),
	)
	if err != nil {
		log.Println("error sending reply message")
	}
}

func (s *botService) sendError(msg error, request slacker.Request, response slacker.ResponseWriter, mention bool) {
	agent := request.Event().User
	client := response.Client()
	conInfo, err := client.GetConversationInfo(request.Event().Channel, false)
	if err != nil {
		log.Println("error get conversation info")
		return
	}

	if !conInfo.IsIM && mention {
		response.Reply("*Agent:* <@" + agent + ">" + "\n_Response only visible to agent_")
	}

	_, err = client.PostEphemeral(
		request.Event().Channel,
		request.Event().User,
		slack.MsgOptionText("*Error:* "+"_"+msg.Error()+"_", false),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		log.Println("error sending error message")
	}
}

func (s *botService) mentionAgent(request slacker.Request, response slacker.ResponseWriter) {
	agent := request.Event().User
	client := response.Client()
	conInfo, err := client.GetConversationInfo(request.Event().Channel, false)
	if err != nil {
		log.Println("error get conversation info")
		return
	}

	if !conInfo.IsIM {
		response.Reply("*Agent:* <@" + agent + ">" + "\n_Response only visible to agent_")
	}
}
