package protoapi

/*******************************************************
* 全局服务组件Visitor
* 该步骤主要用于分离服务层与访问层,从而支持更灵活的组件部署模式
 *******************************************************/

var globalService []*component

type component struct {
	Registry  Registry
	Implement interface{}
	Aspects   []ServiceProcessor
}

func globalServiceVisitor(s *Server) {
	for _, p := range globalService {
		s.RegisterService(p.Registry, p.Implement, p.Aspects...)
	}
}

func RegisterService(registry Registry, implement interface{}, aspects ...ServiceProcessor) {
	globalService = append(globalService, &component{
		Registry:  registry,
		Implement: implement,
		Aspects:   aspects,
	})
}
