package proxy

import (
	"context"
	"fmt"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol"
	"github.com/apache/dubbo-go/protocol/invocation"
	"reflect"
)

var (
	typeOfInterface = reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem()).Type()
	typeOfResult    = reflect.Zero(reflect.TypeOf((*protocol.Result)(nil)).Elem()).Type()
)

// Implement
func NewGenericResultServiceProxyImplFunc(attachments map[string]string) ImplementFunc {
	return func(p *Proxy, rpc common.RPCService) {
		serviceValue := reflect.ValueOf(rpc)
		serviceElem := serviceValue.Elem()
		serviceType := serviceElem.Type()
		// Enforce GenericResultService
		if _, ok := rpc.(*config.GenericResultService); !ok {
			logger.Errorf("Generic proxy requires <*config.GenericResultService> to implement, was: %s",
				serviceType.String())
			return
		}
		numField := serviceElem.NumField()
		// Impl proxy functions
		for i := 0; i < numField; i++ {
			f := serviceElem.Field(i)
			if !(f.Kind() == reflect.Func && f.IsValid() && f.CanSet()) {
				continue
			}
			funcField := serviceType.Field(i)
			funcName := funcField.Tag.Get("dubbo")
			if funcName == "" {
				funcName = funcField.Name
			}
			// Only generic method: Invoke/$invoke
			if funcName != "Invoke" && funcName != "$invoke" {
				continue
			}
			// Enforce param types：
			// Invoke(Context, []{Method, []{Arguments}}) protocol.Result
			inNum := funcField.Type.NumIn()
			if inNum != 2 {
				logger.Errorf("Generic func requires 2 in-arg type, func: %s(%s), was: %d",
					funcField.Name, funcField.Type.String(), inNum)
				continue
			}
			outNum := funcField.Type.NumOut()
			if outNum != 1 || funcField.Type.Out(0) != typeOfResult {
				logger.Errorf("Generic func requires 1 out-result type, func: %s(%s), was: %d",
					funcField.Name, funcField.Type.String(), outNum)
				continue
			}
			f.Set(reflect.MakeFunc(f.Type(), makeGenericProxyImpl(funcName, p, attachments)))
		}
	}
}

func makeGenericProxyImpl(funcName string, proxy *Proxy, usrAttr map[string]string) func([]reflect.Value) []reflect.Value {
	// Invoke(goctx, []interface{}{service.Method, types, values}) protocol.Result;
	return func(funArgs []reflect.Value) []reflect.Value {
		// Context
		invCtx := funArgs[0].Interface().(context.Context)
		invReply := reflect.New(typeOfInterface)
		invArgs := funArgs[1].Interface().([]interface{})
		invValues := funArgs[1]
		inv := invocation.NewRPCInvocationWithOptions(
			invocation.WithMethodName(funcName),
			invocation.WithCallBack(proxy.GetCallback()),
			invocation.WithArguments(invArgs),
			invocation.WithParameterValues([]reflect.Value{invValues}),
			invocation.WithReply(invReply.Interface()),
		)
		if invCtx == nil {
			invCtx = context.Background()
		}
		// Build-in attachments
		for k, value := range usrAttr {
			inv.SetAttachments(k, value)
		}
		// User context attachments
		attachments := invCtx.Value(constant.AttachmentKey)
		if ssmap, ok := attachments.(map[string]string); ok {
			for k, v := range ssmap {
				inv.SetAttachments(k, v)
			}
		} else if simap, ok := attachments.(map[string]interface{}); ok {
			for k, v := range simap {
				inv.SetAttachments(k, fmt.Sprintf("%s", v))
			}
		} else {
			logger.Errorf("Attachments requires map[string]string OR map[string]interface{}, was: %T",
				attachments)
		}
		// Invoke and unwrap reply value
		result := proxy.GetInvoker().Invoke(invCtx, inv)
		if nil == result.Error() {
			result.SetResult(invReply.Elem().Interface())
		}
		return []reflect.Value{reflect.ValueOf(&result).Elem()}
	}
}
