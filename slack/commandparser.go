package slack

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseSlashCommand(commandArguments string) []string {
	tokens := []string{}
	curTokenStart := 0
	inQuotedToken := false
	inUnquotedToken := false

	if commandArguments == "" {
		return tokens
	}

	for i, c := range commandArguments {
		if c == '"' {
			if !inQuotedToken && !inUnquotedToken {
				curTokenStart = i + 1
				inQuotedToken = true
			} else if inQuotedToken {
				tokens = append(tokens, commandArguments[curTokenStart:i])
				inQuotedToken = false
			}
		} else if c == ' ' {
			if inUnquotedToken {
				tokens = append(tokens, commandArguments[curTokenStart:i])
				inUnquotedToken = false
			}
		} else {
			if !inQuotedToken && !inUnquotedToken {
				curTokenStart = i
				inUnquotedToken = true
			}
		}
	}
	if (inQuotedToken || inUnquotedToken) && commandArguments[curTokenStart:] != "" {
		tokens = append(tokens, commandArguments[curTokenStart:])
	}
	return tokens
}
