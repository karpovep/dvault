package constants

type AppContextKey string

const (
	Ctx            AppContextKey = "ctx"
	AppConfig      AppContextKey = "config"
	UserRepository AppContextKey = "userRepository"
	UserService    AppContextKey = "userService"
)
