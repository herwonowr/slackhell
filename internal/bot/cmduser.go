package bot

import (
	"errors"
	"strings"

	"github.com/herwonowr/slackhell/internal/repository"
	"github.com/shomali11/slacker"
)

func (s *botService) cmdUser() {
	cmd := &slacker.CommandDefinition{
		Description: "Get authorized user information",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, true)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			userParam := request.Param("user")
			trimUserPrefix := strings.TrimPrefix(userParam, "<@")
			userID := strings.TrimSuffix(trimUserPrefix, ">")
			userData, err := s.slack.GetUserInfo(userID)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			user, err := s.service.GetAccountBySlackID(userData.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			role := "Agent"
			if user.Role == repository.AdminRole {
				role = "Admin"
			}

			if user.Role == repository.AgentRole {
				role = "Agent"
			}

			userInfo := "```" + "Slack ID: " + user.SlackID + "\nReal Name: " + user.SlackRealName + "\nRole: " + role + "```"
			s.sendReply(userInfo, request, response, true)
		},
	}
	s.slack.Command("user <user>", cmd)
}

func (s *botService) cmdUsers() {
	cmd := &slacker.CommandDefinition{
		Description: "Get authorized users",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, true)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			users, err := s.service.GetAccounts()
			if err != nil {
				s.sendError(errors.New("get authorized users"), request, response, true)
				return
			}

			//TODO: add get users take and skip
			s.mentionAgent(request, response)
			for _, user := range users {
				userRole := "Agent"
				if user.Role == repository.AdminRole {
					userRole = "Admin"
				}

				userInfo := "```" + "Slack ID: " + user.SlackID + "\nReal Name: " + user.SlackRealName + "\nRole: " + userRole + "```"
				s.sendReply(userInfo, request, response, false)
			}
		},
	}
	s.slack.Command("users", cmd)
}

func (s *botService) cmdAddUser() {
	cmd := &slacker.CommandDefinition{
		Description: "Add authorized user",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, true)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			userParam := request.Param("user")
			trimUserPrefix := strings.TrimPrefix(userParam, "<@")
			userID := strings.TrimSuffix(trimUserPrefix, ">")
			user, err := s.slack.GetUserInfo(userID)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			roleParam := request.Param("role")
			if !(roleParam == "admin" || roleParam == "agent") {
				s.sendError(errors.New("invalid user role, valid role `admin` and `agent`"), request, response, true)
				return
			}

			role := repository.AgentRole
			if roleParam == "admin" {
				role = repository.AdminRole
			}

			if roleParam == "agent" {
				role = repository.AgentRole
			}

			err = s.service.CreateAccount(role, user.ID, user.RealName)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			userInfo := "User has been added successfully\n" + "```" + "Slack ID: " + user.ID + "\nReal Name: " + user.RealName + "\nRole: " + strings.Title(strings.ToLower(roleParam)) + "```"
			s.sendReply(userInfo, request, response, true)
		},
	}
	s.slack.Command("useradd <user> <role>", cmd)
}

func (s *botService) cmdUpdateUser() {
	cmd := &slacker.CommandDefinition{
		Description: "Update authorized user",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, true)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			userParam := request.Param("user")
			trimUserPrefix := strings.TrimPrefix(userParam, "<@")
			userID := strings.TrimSuffix(trimUserPrefix, ">")
			userInfo, err := s.slack.GetUserInfo(userID)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			roleParam := request.Param("role")
			if !(roleParam == "admin" || roleParam == "agent") {
				s.sendError(errors.New("invalid user role, valid role `admin` and `agent`"), request, response, true)
				return
			}

			role := repository.AgentRole
			if roleParam == "admin" {
				role = repository.AdminRole
			}

			if roleParam == "agent" {
				role = repository.AgentRole
			}

			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if agent == user.SlackID {
				s.sendError(errors.New("cannot update your self, please contact other administrator"), request, response, true)
				return
			}

			if user.Role == role {
				s.sendError(errors.New("user already have "+strings.ToLower(roleParam)+" permission"), request, response, true)
				return
			}

			err = s.service.PutAccount(user.ID, role, user.SlackID, user.SlackRealName)
			if err != nil {
				s.sendError(errors.New("update authorize user, format `userupdate <user>`"), request, response, true)
				return
			}

			userResponse := "User @" + user.SlackRealName + " has been updated successfully"
			s.sendReply(userResponse, request, response, true)
		},
	}
	s.slack.Command("userupdate <user> <role>", cmd)
}

func (s *botService) cmdDeleteUser() {
	cmd := &slacker.CommandDefinition{
		Description: "Delete authorized user",
		AuthorizationFunc: func(request slacker.Request) bool {
			return s.authorized(request.Event().User, true)
		},
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			agent := request.Event().User
			userParam := request.Param("user")
			trimUserPrefix := strings.TrimPrefix(userParam, "<@")
			userID := strings.TrimSuffix(trimUserPrefix, ">")

			userInfo, err := s.slack.GetUserInfo(userID)
			if err != nil {
				s.sendError(errors.New("get user information"), request, response, true)
				return
			}

			user, err := s.service.GetAccountBySlackID(userInfo.ID)
			if err != nil {
				s.sendError(err, request, response, true)
				return
			}

			if agent == user.SlackID {
				s.sendError(errors.New("cannot delete your self, please contact other administrator"), request, response, true)
				return
			}

			err = s.service.DeleteAccount(user.ID)
			if err != nil {
				s.sendError(errors.New("delete authorize user, format `userdel <user>`"), request, response, true)
				return
			}

			userResponse := "User @" + user.SlackRealName + " has been successfully deleted"
			s.sendReply(userResponse, request, response, true)
		},
	}
	s.slack.Command("userdel <user>", cmd)
}
