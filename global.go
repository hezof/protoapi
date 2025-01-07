package protoapi

/*******************************************************
* 全局服务组件Visitor
* 该步骤主要用于分离服务层与访问层,从而支持更灵活的组件部署模式
 *******************************************************/

var _globalRegister []*register

type register struct {
	Registry  ServiceRegistry
	Implement interface{}
	Aspects   []ServiceAspect
}

func RegisterService(registry ServiceRegistry, implement interface{}, aspects ...ServiceAspect) {
	_globalRegister = append(_globalRegister, &register{
		Registry:  registry,
		Implement: implement,
		Aspects:   aspects,
	})
}
