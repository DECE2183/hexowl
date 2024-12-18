package runtime

import (
	"github.com/dece2183/hexowl/v2/types"
)

type actionHandler func(rn *Runtime, opLeft, opRight types.Value) (interface{}, error)

var actionHandlerMap map[types.OperatorType]actionHandler

func init() {
	actionHandlerMap = map[types.OperatorType]actionHandler{
		types.O_NONE:         implNONE,
		types.O_DECLFUNC:     implDECLFUNC,
		types.O_SEQUENCE:     implSEQUENCE,
		types.O_ASSIGN:       implASSIGN,
		types.O_ASSIGNLOCAL:  implASSIGNLOCAL,
		types.O_ASSIGNMINUS:  implASSIGNMINUS,
		types.O_ASSIGNPLUS:   implASSIGNPLUS,
		types.O_ASSIGNMUL:    implASSIGNMUL,
		types.O_ASSIGNDIV:    implASSIGNDIV,
		types.O_ASSIGNBITAND: implASSIGNBITAND,
		types.O_ASSIGNBITOR:  implASSIGNBITOR,
		types.O_ENUMERATE:    implENUMERATE,
		types.O_LOGICOR:      implLOGICOR,
		types.O_LOGICAND:     implLOGICAND,
		types.O_EQUALITY:     implEQUALITY,
		types.O_NOTEQ:        implNOTEQ,
		types.O_MORE:         implMORE,
		types.O_LESS:         implLESS,
		types.O_MOREEQ:       implMOREEQ,
		types.O_LESSEQ:       implLESSEQ,
		types.O_PLUS:         implPLUS,
		types.O_MINUS:        implMINUS,
		types.O_MULTIPLY:     implMULTIPLY,
		types.O_DIVIDE:       implDIVIDE,
		types.O_MODULO:       implMODULO,
		types.O_POWER:        implPOWER,
		types.O_BITOR:        implBITOR,
		types.O_BITAND:       implBITAND,
		types.O_BITXOR:       implBITXOR,
		types.O_BITCLEAR:     implBITCLEAR,
		types.O_BITINVERSE:   implBITINVERSE,
		types.O_LEFTSHIFT:    implLEFTSHIFT,
		types.O_RIGHTSHIFT:   implRIGHTSHIFT,
		types.O_LOGICNOT:     implLOGICNOT,
		types.O_POPCNT:       implPOPCNT,
		types.O_CALLFUNC:     implCALLFUNC,
	}
}
