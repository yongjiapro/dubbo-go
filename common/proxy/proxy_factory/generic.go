package proxy_factory

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/proxy"
	"github.com/apache/dubbo-go/protocol"
)

type GenericProxyFactory struct {
	*DefaultProxyFactory
}

func NewGenericProxyFactory(opts ...proxy.Option) proxy.ProxyFactory {
	return &GenericProxyFactory{
		DefaultProxyFactory: NewDefaultProxyFactory(opts...).(*DefaultProxyFactory),
	}
}

func (f *GenericProxyFactory) GetProxy(invoker protocol.Invoker, url *common.URL) *proxy.Proxy {
	return f.GetAsyncProxy(invoker, nil, url)
}

func (f *GenericProxyFactory) GetAsyncProxy(invoker protocol.Invoker, callback interface{}, url *common.URL) *proxy.Proxy {
	attrs := map[string]string{}
	attrs[constant.ASYNC_KEY] = url.GetParam(constant.ASYNC_KEY, "false")
	return proxy.NewProxyWith(
		invoker, callback, attrs,
		proxy.WithProxyImplementFunc(proxy.NewGenericResultServiceProxyImplFunc(attrs)),
	)
}

func (f *GenericProxyFactory) GetInvoker(url common.URL) protocol.Invoker {
	return f.DefaultProxyFactory.GetInvoker(url)
}
