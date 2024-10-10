package dto

type Validator interface {
	Validate() map[string]string
}
