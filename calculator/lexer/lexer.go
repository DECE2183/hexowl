package lexer

import (
	"strings"

	"github.com/dece2183/hexowl/calculator/types"
)

const (
	_STRING_LITERALS   = "_@QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm"
	_NUM_SCI_LITERALS  = "0123456789.eE-+_"
	_NUM_DEC_LITERALS  = "0123456789._"
	_NUM_HEX_LITERALS  = "0123456789ABCDEFabcdef_"
	_NUM_BIN_LITERALS  = "01_"
	_CONTROL_LITERALS  = "()"
	_OPERATOR_LITERALS = ";#?:=-+*/%^!&|~<>,"
)

func Parse(str string) []types.Token {
	token := make([]types.Token, 0)
	tokenType := types.T_NUM_DEC
	tokenDone := false
	tokenBegin := -1

	for i, c := range str {
		if tokenBegin > -1 {
			switch tokenType {
			case types.T_UNIT:
				if !(strings.Contains(_STRING_LITERALS, string(c)) || strings.Contains(_NUM_DEC_LITERALS, string(c))) {
					tokenDone = true
				}
			case types.T_NUM_DEC, types.T_NUM_HEX, types.T_NUM_BIN, types.T_NUM_SCI:
				if (c == 'x' || c == 'b') && i-tokenBegin == 1 && tokenType == types.T_NUM_DEC {
					if c == 'x' {
						tokenType = types.T_NUM_HEX
					} else {
						tokenType = types.T_NUM_BIN
					}
					tokenBegin += 2
				} else if (c == 'e' || c == 'E') && tokenType == types.T_NUM_DEC {
					tokenType = types.T_NUM_SCI
				} else {
					switch tokenType {
					case types.T_NUM_SCI:
						if !strings.Contains(_NUM_SCI_LITERALS, string(c)) {
							tokenDone = true
						}
					case types.T_NUM_DEC:
						if !strings.Contains(_NUM_DEC_LITERALS, string(c)) {
							tokenDone = true
						}
					case types.T_NUM_HEX:
						if !strings.Contains(_NUM_HEX_LITERALS, string(c)) {
							tokenDone = true
						}
					case types.T_NUM_BIN:
						if !strings.Contains(_NUM_BIN_LITERALS, string(c)) {
							tokenDone = true
						}
					}
				}
			case types.T_STR:
				if c == '"' {
					tokenDone = true
					token = append(token, types.Token{Type: tokenType, Literal: str[tokenBegin:i]})
					tokenBegin = -1
					continue
				}
			case types.T_CTL:
				tokenDone = true
			case types.T_OP:
				if !strings.Contains(_OPERATOR_LITERALS, string(c)) {
					tokenDone = true
				}
			}

			if tokenDone && tokenType != types.T_NONE {
				token = append(token, types.Token{Type: tokenType, Literal: str[tokenBegin:i]})
				tokenBegin = -1
			}
		}

		if tokenBegin < 0 {
			tokenBegin = i
			tokenDone = false

			if c == '"' {
				tokenType = types.T_STR
				tokenBegin++
			} else if strings.Contains(_STRING_LITERALS, string(c)) {
				tokenType = types.T_UNIT
			} else if strings.Contains(_NUM_DEC_LITERALS, string(c)) {
				tokenType = types.T_NUM_DEC
			} else if strings.Contains(_CONTROL_LITERALS, string(c)) {
				tokenType = types.T_CTL
			} else if strings.Contains(_OPERATOR_LITERALS, string(c)) {
				tokenType = types.T_OP
			} else {
				tokenBegin = -1
				tokenType = types.T_NONE
			}
		}
	}

	if tokenBegin > -1 && tokenType != types.T_NONE {
		token = append(token, types.Token{Type: tokenType, Literal: str[tokenBegin:]})
	}

	return token
}
