package core

import (
	"errors"
	"github.com/codnect/goo"
	"sync"
)

type TypeConverter interface {
	Support(sourceTyp goo.Type, targetTyp goo.Type) bool
	Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (interface{}, error)
}

type StringToNumberConverter struct {
}

func NewStringToNumberConverter() StringToNumberConverter {
	return StringToNumberConverter{}
}

func (converter StringToNumberConverter) Support(sourceTyp goo.Type, targetTyp goo.Type) bool {
	if sourceTyp.IsString() && targetTyp.IsNumber() && goo.ComplexType == targetTyp.(goo.Number).GetNumberType() {
		return true
	}
	return false
}

func (converter StringToNumberConverter) Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (interface{}, error) {
	if sourceTyp.IsString() && targetTyp.IsNumber() && goo.ComplexType == targetTyp.(goo.Number).GetNumberType() {
		number := targetTyp.(goo.Number)
		return sourceTyp.(goo.String).ToNumber(source.(string), number)
	}
	return nil, errors.New("unsupported type")
}

type NumberToStringConverter struct {
}

func NewNumberToStringConverter() NumberToStringConverter {
	return NumberToStringConverter{}
}

func (converter NumberToStringConverter) Support(sourceTyp goo.Type, targetTyp goo.Type) bool {
	if targetTyp.IsString() && sourceTyp.IsNumber() && goo.ComplexType == sourceTyp.(goo.Number).GetNumberType() {
		return true
	}
	return false
}

func (converter NumberToStringConverter) Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (interface{}, error) {
	if targetTyp.IsString() && sourceTyp.IsNumber() && goo.ComplexType == sourceTyp.(goo.Number).GetNumberType() {
		return targetTyp.(goo.Number).ToString(source), nil
	}
	return nil, errors.New("unsupported type")
}

type StringToBooleanConverter struct {
}

func NewStringToBooleanConverter() StringToBooleanConverter {
	return StringToBooleanConverter{}
}

func (converter StringToBooleanConverter) Support(sourceTyp goo.Type, targetTyp goo.Type) bool {
	if sourceTyp.IsString() && targetTyp.IsBoolean() {
		return true
	}
	return false
}

func (converter StringToBooleanConverter) Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (result interface{}, err error) {
	if sourceTyp.IsString() && targetTyp.IsBoolean() {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(r.(string))
			}
		}()
		return targetTyp.(goo.Boolean).ToBoolean(source.(string)), nil
	}
	return nil, errors.New("unsupported type")
}

type BooleanToStringConverter struct {
}

func NewBooleanToStringConverter() BooleanToStringConverter {
	return BooleanToStringConverter{}
}

func (converter BooleanToStringConverter) Support(sourceTyp goo.Type, targetTyp goo.Type) bool {
	if targetTyp.IsString() && sourceTyp.IsBoolean() {
		return true
	}
	return false
}

func (converter BooleanToStringConverter) Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (interface{}, error) {
	if targetTyp.IsString() && sourceTyp.IsBoolean() {
		return sourceTyp.(goo.Boolean).ToString(source.(bool)), nil
	}
	return nil, errors.New("unsupported type")
}

type TypeConverterRegistry interface {
	RegisterConverter(converter TypeConverter)
}

type TypeConverterService interface {
	TypeConverterRegistry
	CanConvert(sourceTyp goo.Type, targetTyp goo.Type) bool
	Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (interface{}, error)
}

type DefaultTypeConverterService struct {
	converters map[goo.Type]TypeConverter
	mu         sync.RWMutex
}

func NewDefaultTypeConverterService() *DefaultTypeConverterService {
	converterService := &DefaultTypeConverterService{
		converters: make(map[goo.Type]TypeConverter, 0),
	}
	converterService.registerDefaultConverters()
	return converterService
}

func (cs *DefaultTypeConverterService) registerDefaultConverters() {
	/* number to string and string to number */
	cs.RegisterConverter(NewNumberToStringConverter())
	cs.RegisterConverter(NewStringToNumberConverter())
	/* bool to string and string to bool */
	cs.RegisterConverter(NewBooleanToStringConverter())
	cs.RegisterConverter(NewStringToBooleanConverter())
}

func (cs *DefaultTypeConverterService) CanConvert(sourceTyp goo.Type, targetTyp goo.Type) bool {
	var result bool
	cs.mu.Lock()
	for _, converter := range cs.converters {
		if converter.Support(sourceTyp, targetTyp) {
			result = true
			break
		}
	}
	cs.mu.Unlock()
	return result
}

func (cs *DefaultTypeConverterService) Convert(source interface{}, sourceTyp goo.Type, targetTyp goo.Type) (result interface{}, err error) {
	var typConverter TypeConverter
	cs.mu.Lock()
	for _, converter := range cs.converters {
		if converter.Support(sourceTyp, targetTyp) {
			typConverter = converter
		}
	}
	cs.mu.Unlock()
	if typConverter != nil {
		defer func() {
			err = errors.New("converting error has just occurred")
		}()
		result, err = typConverter.Convert(source, sourceTyp, targetTyp)
	}
	return
}

func (cs *DefaultTypeConverterService) RegisterConverter(converter TypeConverter) {
	if converter == nil {
		panic("converter must not be nil")
	}
	cs.mu.Lock()
	cs.converters[goo.GetType(converter)] = converter
	cs.mu.Unlock()
}
