package oapi

import (
	"net/http"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/wI2L/fizz"
)

type Info struct {
	fizz.OperationOption
}

func WrapInfos(infos ...Info) []fizz.OperationOption {
	var wrapped []fizz.OperationOption
	for _, info := range infos {
		wrapped = append(wrapped, info.OperationOption)
	}
	return wrapped
}

func ID(id string) Info {
	return Info{
		OperationOption: fizz.ID(id),
	}
}

func Summary(summary string) Info {
	return Info{
		OperationOption: fizz.Summary(summary),
	}
}

func Description(description string) Info {
	return Info{
		OperationOption: fizz.Description(description),
	}
}

func Header(name, description string, model interface{}) Info {
	return Info{
		OperationOption: fizz.Header(name, description, model),
	}
}

type ResponseBuilder struct {
	status   string
	desc     string
	model    interface{}
	examples []interface{}
}

type ResponseOption func(*ResponseBuilder)

func Response(status int, opts ...ResponseOption) Info {
	r := &ResponseBuilder{
		status: strconv.Itoa(status),
		desc:   http.StatusText(status),
	}
	for _, opt := range opts {
		opt(r)
	}
	if len(r.examples) == 0 {
		r.examples = append(r.examples, nil)
	}
	return Info{
		OperationOption: fizz.Response(r.status, r.desc, r.model, nil, r.examples[0]),
	}
}

func WithResponseDesc(description string) func(*ResponseBuilder) {
	return func(b *ResponseBuilder) {
		b.desc = description
	}
}

func WithResponseModel(model interface{}) func(*ResponseBuilder) {
	return func(b *ResponseBuilder) {
		b.model = model
		b.examples = append(b.examples, gofakeit.Struct(&model))
	}
}
