package bacnet

import (
	"fmt"
	"strings"
)

type PropertyName string

const (
	PropPresentValue PropertyName = "present-value"
	PropStatusFlags  PropertyName = "status-flags"
	PropUnits        PropertyName = "units"
	PropObjectName   PropertyName = "object-name"
	PropDescription  PropertyName = "description"
	PropOutOfService PropertyName = "out-of-service"
)

func ParseProperty(name string) (PropertyName, error) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	switch PropertyName(normalized) {
	case PropPresentValue, PropStatusFlags, PropUnits, PropObjectName, PropDescription, PropOutOfService:
		return PropertyName(normalized), nil
	default:
		return "", fmt.Errorf("unsupported property %s", name)
	}
}
func Writable(property PropertyName) bool {
	return property == PropPresentValue || property == PropOutOfService
}
func ReadOnly(property PropertyName) bool { return !Writable(property) }

type ServiceError struct {
	Class string
	Code  string
}

func (e ServiceError) Error() string            { return e.Class + ":" + e.Code }
func NewServiceError(class, code string) error  { return ServiceError{Class: class, Code: code} }
func ValidateObjectType(objectType uint16) bool { return objectType < 1024 }
func ValidateInstance(instance uint32) bool     { return instance <= 0x3fffff }
