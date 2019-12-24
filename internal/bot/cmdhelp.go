package bot

import (
	"github.com/herwonowr/slacker"
)

func (s *botService) cmdHelp() {
	cmd := &slacker.CommandDefinition{
		Description: "Show help command",
		Handler: func(request slacker.Request, response slacker.ResponseWriter) {
			res := "`help` - Show help commands\n`cmd <command>` - Execute command on the target client\n`generate <type>` - Generate shellcode, valid type `php` and `asp`\n`shellcodes` - List generated shellcodes\n`shellupdate <key> <endpoint>` - Update shellcode endpoint\n`shelldel <key>` - Delete shellcode\n`user <user>` - Get authorized user information\n`users` - List authorized users\n`useradd <user> <role>` - Add authorized user, valid role `admin` and `agent`\n`userupdate <user> <role>` - Update authorized user, valid role `admin` and `agent`\n`userdel <user>` - Delete authorized user\n`session` - Get active session\n`sessionset <shellcode-key>` - Set active shell session\n`sessionkill` - Kill active shell session"
			s.sendReply(res, request, response, true)
		},
	}

	s.slack.Help(cmd)
}
