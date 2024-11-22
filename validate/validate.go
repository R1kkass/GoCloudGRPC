package validate

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidateStruct struct {
	Rule  string
	Value any
}

type ValidType map[string]ValidateStruct

func Valid(options ValidType) error {
	for key, valueStruct := range options {
		validateRules := strings.Split(valueStruct.Rule, "|")
		value := valueStruct.Value
		for _, validRule := range validateRules {
			if err := validRequired(value, validRule, key); err != nil {
				return err
			}
			if err := validType(value, validRule, key); err != nil {
				return err
			}
			if err := validMax(value, validRule, key); err != nil {
				return err
			}
            if err := validMin(value, validRule, key); err != nil {
				return err
			}
		}
	}
	return nil
}

func validType(value any, validRule string, key string) error {
	r, _ := regexp.Compile(`string|int|int8|int16|int32|int64|float64|float32|uint|uint8|uint16|uint32|uin64`)
	if matched := r.MatchString(validRule); matched {
		if reflect.TypeOf(value).String() != validRule {
			return errors.New(key + " non valid: " + validRule)
		}
	}
	return nil
}

func validRequired(value any, validRule string, key string) error {
	if validRule == "required" {
		typeValue := reflect.TypeOf(value).String()
		rInt, _ := regexp.Compile(`int|int8|int16|int32|int64|float64|float32|uint|uint8|uint16|uint32|uin64`)
		if rInt.MatchString(typeValue) && value == 0 {
			return errors.New("Value non valid: " + key + " must required")
		}
		if typeValue == "string" && value == "" {
			return errors.New("Value non valid: " + key + " must required")
		}
	}
	return nil
}

func validMax(value any, validRule string, key string) error {
	r, _ := regexp.Compile(`max:`)
	if matched := r.MatchString(validRule); matched {
		typeValue := reflect.TypeOf(value).String()
		maxCount, err := strconv.Atoi(strings.Split(validRule, ":")[1])
		if err != nil {
			return errors.New("Value non valid: " + key + " must required")
		}

		rInt, _ := regexp.Compile(`int|int8|int16|int32|int64|float64|float32|uint|uint8|uint16|uint32|uin64`)

		if matched := rInt.MatchString(typeValue); matched && maxCount > value.(int) {
			return errors.New(key + " non valid: " + strconv.Itoa(maxCount) + " max value")
		}
		if typeValue == "string" && maxCount > utf8.RuneCountInString(value.(string)) {
			return errors.New(key + " non valid: " + strconv.Itoa(maxCount) + " max value")
		}
	}
	return nil
}

func validMin(value any, validRule string, key string) error {
	r, _ := regexp.Compile(`min:`)
	if matched := r.MatchString(validRule); matched {
		typeValue := reflect.TypeOf(value).String()
		minCount, err := strconv.Atoi(strings.Split(validRule, ":")[1])
		if err != nil {
			return errors.New("Value non valid: " + key + " must required")
		}

		rInt, _ := regexp.Compile(`int|int8|int16|int32|int64|float64|float32|uint|uint8|uint16|uint32|uin64`)

		if matched := rInt.MatchString(typeValue); matched && minCount < value.(int) {
			return errors.New(key + " non valid: " + strconv.Itoa(minCount) + " max value")
		}
		if typeValue == "string" && minCount < utf8.RuneCountInString(value.(string)) {
			return errors.New(key + " non valid: " + strconv.Itoa(minCount) + " max value")
		}
	}
	return nil
}