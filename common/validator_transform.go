package common

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	playground "github.com/go-playground/validator/v10"
	localeTransId "github.com/go-playground/validator/v10/translations/id"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	validator  *playground.Validate
	translator ut.Translator
	initOnce   sync.Once
)

func initValidator() {
	initOnce.Do(func() {
		validator = playground.New()
		lang := id.New()
		uni := ut.New(lang, lang)
		translator, _ = uni.GetTranslator("id")
		if err := localeTransId.RegisterDefaultTranslations(validator, translator); err != nil {
			panic(fmt.Sprintf("Failed to register Indonesian validator translations: %v", err))
		}
	})
}

// ValidateRequest validates the struct and returns a gRPC error if invalid
func ValidateRequest(s interface{}) error {
	initValidator()

	err := validator.Struct(s)
	if err != nil {
		return scanner(err)
	}
	return nil
}

// scanner translates and formats validation errors into gRPC BadRequest error
func scanner(err error) error {
	var invalidValidationError *playground.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return status.Error(codes.Internal, "Invalid validation error")
	}

	verrs, ok := err.(playground.ValidationErrors)
	if !ok {
		return status.Error(codes.InvalidArgument, "Unknown validation error")
	}

	var violations []*errdetails.BadRequest_FieldViolation
	for _, verr := range verrs {
		field := strings.ToLower(verr.Field())
		description := verr.Translate(translator)
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: description,
		})
	}

	st := status.New(codes.InvalidArgument, "Invalid Argument")
	br := &errdetails.BadRequest{FieldViolations: violations}

	details, err := st.WithDetails(br)
	if err != nil {
		panic(fmt.Sprintf("Failed to attach validation details: %v", err))
	}

	return details.Err()
}
