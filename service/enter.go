package service

type ServiceGroup struct {
	EsService
	BaseService
	JwtService
}

var ServiceGroupApp = new(ServiceGroup)
