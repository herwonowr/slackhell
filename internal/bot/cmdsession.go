package bot

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/herwonowr/slacker"
)

func (s *botService) cmdGetSession() {
	cmd := &slacker.CommandDefinition{
		Description: "Get shell session",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			session, err := s.service.GetSessionByID(userInfo.ID)
			if err != nil && err.Error() == "key not found" {
				s.sendError(errors.New("no active session"), request, response, true)
				return
			} else if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			shellcode, err := s.service.GetShellcodeByKey(session.Key)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			sessionResponse := "Active session\n" + "```" + "Key: " + shellcode.ShellKey + "\nEndpoint: " + shellcode.Endpoint + "```"
			s.sendReply(sessionResponse, request, response, true)
		},
	}
	s.slack.Command("session", cmd)
}

func (s *botService) cmdSetSession() {
	cmd := &slacker.CommandDefinition{
		Description: "Set shell session",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			key := request.Param("shellcode-key")
			if govalidator.IsNull(key) {
				s.sendError(errors.New("invalid shellcode key, format `sessionset <shellcode-key`"), request, response, true)
				return
			}

			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			shellcode, err := s.service.GetShellcodeByKey(key)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if govalidator.IsNull(shellcode.Endpoint) {
				s.sendError(errors.New("shellcode is not active, please set endpoint"), request, response, true)
				return
			}

			session, err := s.service.GetSessionByID(userInfo.ID)
			if err != nil && err.Error() == "key not found" {
				err := s.service.CreateSession(user.SlackID, shellcode.ShellKey)
				if err != nil {
					s.sendError(err, request, response, true)
					return
				}

				sessionResponse := "Session has been created successfully\nSession key: " + "`" + shellcode.ShellKey + "`"
				s.sendReply(sessionResponse, request, response, true)
			} else if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if session != nil {
				err := s.service.PutSession(session.ID, shellcode.ShellKey)
				if err != nil {
					s.sendError(err, request, response, true)
					return
				}

				sessionResponse := "Session has been created successfully\nSession key: " + "`" + shellcode.ShellKey + "`"
				s.sendReply(sessionResponse, request, response, true)
			}
		},
	}
	s.slack.Command("sessionset <shellcode-key>", cmd)
}

func (s *botService) cmdKillSession() {
	cmd := &slacker.CommandDefinition{
		Description: "Delete shell session",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			session, err := s.service.GetSessionByID(userInfo.ID)
			if err != nil && err.Error() == "key not found" {
				s.sendError(errors.New("no active session"), request, response, true)
				return
			} else if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			err = s.service.DeleteSession(session.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			sessionResponse := "Session has been killed successfully"
			s.sendReply(sessionResponse, request, response, true)
		},
	}
	s.slack.Command("sessionkill", cmd)
}
