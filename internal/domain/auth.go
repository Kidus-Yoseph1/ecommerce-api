package domain

type AuthServiceInterface interface {
	Register(fullName, email, password string) error
	Login(email, password string) (string, error)
}
