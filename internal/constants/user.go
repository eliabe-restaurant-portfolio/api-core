package constants

import "time"

const MAX_FAILED_LOGIN_ATTEMPTS = 5
const WAIT_UNTIL_NEXT_ATTEMPT = 5 * time.Second

type UserStatus string

var (
	UserActive   UserStatus = "active"
	UserBloqued  UserStatus = "blocked"
	UserInactive UserStatus = "inactive"
)
