package parser

// SyntaxError - error thrown when parser was expecting different token from lexer
type SyntaxError struct {
	expecting []string
	got       string
}

func (m *SyntaxError) Error() string {
	var expectingText string

	if len(m.expecting) == 1 {
		expectingText = m.expecting[0]
	} else {
		for i, expected := range m.expecting {
			expectingText += expected
			if i != len(m.expecting)-1 {
				expectingText += ", "
			}
		}
	}

	return "syntax error, expecting: {" + expectingText + "}, got: {" + m.got + "}"
}

// SyntaxCommandExpectedError - error thrown when there was command that logically should only
// appear after certain different command, but it wasn't found
type SyntaxCommandExpectedError struct {
	command        string
	neededCommands []string
}

func (m *SyntaxCommandExpectedError) Error() string {
	var neededCommandsText string

	if len(neededCommandsText) == 1 {
		neededCommandsText = m.neededCommands[0] + " command"
	} else if len(neededCommandsText) == 2 {
		neededCommandsText = m.neededCommands[0] + " or " + m.neededCommands[1] + " commands"
	} else {
		for i, command := range m.neededCommands {
			if i == len(m.neededCommands)-1 {
				neededCommandsText += " or "
			}

			neededCommandsText += command

			if i != len(m.neededCommands)-1 || i != len(m.neededCommands)-2 {
				neededCommandsText += ", "
			}
		}
		neededCommandsText += " commands"
	}

	return "syntax error, {" + m.command + "} command needs {" + neededCommandsText + "} before"
}

// SyntaxInvalidCommandError - error thrown when invalid (non-existing) type of command has been
// found
type SyntaxInvalidCommandError struct {
	invalidCommand string
}

func (m *SyntaxInvalidCommandError) Error() string {
	return "syntax error, invalid command found: {" + m.invalidCommand + "}"
}

// LogicalExpressionParsingError - error thrown when logical expression inside WHERE statement
// couldn't be parsed correctly
type LogicalExpressionParsingError struct {
	afterToken *string
}

func (m *LogicalExpressionParsingError) Error() string {
	errorMsg := "syntax error, logical expression within WHERE command couldn't be parsed correctly"
	if m.afterToken != nil {
		return errorMsg + ", after {" + *m.afterToken + "} character"
	}
	return errorMsg
}

// ArithmeticLessThanZeroParserError - error thrown when parser found integer value that shouldn't
// be less than 0, but it is
type ArithmeticLessThanZeroParserError struct {
	variable string
}

func (m *ArithmeticLessThanZeroParserError) Error() string {
	return "syntax error, {" + m.variable + "} value should be more than 0"
}

// NoPredecessorParserError - error thrown when parser found integer value that shouldn't
// be less than 0, but it is
type NoPredecessorParserError struct {
	command string
}

func (m *NoPredecessorParserError) Error() string {
	return "syntax error, {" + m.command + "} command can't be used without predecessor"
}
