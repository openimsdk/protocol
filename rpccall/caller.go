package rpccall

import (
	"context"
	"fmt"
	"reflect"

	"github.com/openimsdk/protocol/sdkws"
	"google.golang.org/grpc"
)

func NewRpcCaller[Req, Resp any](name string) RpcCaller[Req, Resp] {
	return RpcCaller[Req, Resp]{
		methodName: name,
	}
}

type RpcCaller[Req, Resp any] struct {
	conn       *grpc.ClientConn
	methodName string
}

func (r RpcCaller[Req, Resp]) SetConn(conn *grpc.ClientConn) {
	r.conn = conn
}

func (r RpcCaller[Req, Resp]) Invoke(ctx context.Context, in *Req, opts ...grpc.CallOption) (*Resp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Resp)
	if err := r.conn.Invoke(ctx, r.methodName, in, out, cOpts...); err != nil {
		return nil, err
	}
	return out, nil
}

func (r RpcCaller[Req, Resp]) Execute(ctx context.Context, req *Req, opts ...grpc.CallOption) error {
	_, err := r.Invoke(ctx, req, opts...)
	return err
}

func (r RpcCaller[Req, Resp]) MethodName() string {
	return r.methodName
}

// ExtractField is a generic function that extracts a field from the response of a given function.
func ExtractField[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C, opts ...grpc.CallOption) (C, error) {
	resp, err := fn(ctx, req, opts...)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

type pagination interface {
	GetPagination() *sdkws.RequestPagination
}

func Page[Req pagination, Resp any, Elem any](ctx context.Context, req Req, api func(ctx context.Context, req Req, opts ...grpc.CallOption) (*Resp, error), fn func(*Resp) []Elem, opts ...grpc.CallOption) ([]Elem, error) {
	if req.GetPagination() == nil {
		vof := reflect.ValueOf(req)
		for {
			if vof.Kind() == reflect.Ptr {
				vof = vof.Elem()
			} else {
				break
			}
		}
		if vof.Kind() != reflect.Struct {
			return nil, fmt.Errorf("request is not a struct")
		}
		fof := vof.FieldByName("Pagination")
		if !fof.IsValid() {
			return nil, fmt.Errorf("request is not valid Pagination field")
		}
		fof.Set(reflect.ValueOf(&sdkws.RequestPagination{}))
	}
	if req.GetPagination().PageNumber < 0 {
		req.GetPagination().PageNumber = 0
	}
	if req.GetPagination().ShowNumber <= 0 {
		req.GetPagination().ShowNumber = 200
	}
	var result []Elem
	for i := int32(0); ; i++ {
		req.GetPagination().PageNumber = i + 1
		resp, err := api(ctx, req, opts...)
		if err != nil {
			return nil, err
		}
		elems := fn(resp)
		result = append(result, elems...)
		if len(elems) < int(req.GetPagination().ShowNumber) {
			break
		}
	}
	return result, nil
}
