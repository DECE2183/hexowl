package types

import "math/rand"

type Context struct {
	Random  *rand.Rand
	System  SystemInterface
	Builtin BuiltinContainer
	User    UserContainer
}

func NewEmptyContext(system SystemInterface) *Context {
	return &Context{
		System: system,
		Random: rand.New(rand.NewSource(system.GetRandomSeed())),
		Builtin: BuiltinContainer{
			functions: make(BuiltinFunctionMap),
			constants: make(BuiltinConstantMap),
		},
		User: UserContainer{
			functions: make(UserFunctionMap),
			variables: make(UserVariableMap),
		},
	}
}
