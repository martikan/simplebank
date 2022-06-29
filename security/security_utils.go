package security

var (
	SecurityUtils securityUtilsInterface = &securityUtils{}
)

type securityUtils struct{}

type securityUtilsInterface interface {
}
