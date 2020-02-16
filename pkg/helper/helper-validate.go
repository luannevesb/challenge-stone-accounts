package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/luannevesb/challenge-stone-accounts/internal/types"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"regexp"
	"strconv"
	"unicode"
)

//Adiciona a Regra Customizada de CPF no goValidator
func InitCustomRule() {
	govalidator.AddCustomRule("cpf", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if !isCPF(val) {

			if message == "" {
				message = "The CPF is invalid."
			}

			return fmt.Errorf(message)
		}
		return nil
	})
}

//Valida a request com o BodyJSON
func ValidateJsonRequest(w http.ResponseWriter, r *http.Request, rules map[string][]string, messages govalidator.MapData) *types.Account {
	var account types.Account

	opts := govalidator.Options{
		Request:         r,
		Data:            &account,
		Rules:           rules,
		Messages:        messages,
		RequiredDefault: true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.ErrorResponse{Error: e})
		return nil
	}

	return &account
}

//Função auxiliar da regra customizada de CPF
func isCPF(doc string) bool {

	const (
		sizeWithoutDigits = 9
		position          = 10
	)

	return isCPFOrCNPJ(
		doc,
		validateCPFFormat,
		sizeWithoutDigits,
		position,
	)
}

//Função auxiliar da regra customizada de CPF
func validateCPFFormat(doc string) bool {

	const (
		pattern = `^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`
	)

	return regexp.MustCompile(pattern).MatchString(doc)
}

//Função auxiliar da regra customizada de CPF
func isCPFOrCNPJ(doc string, validate func(string) bool, size int, position int) bool {

	if !validate(doc) {
		return false
	}

	// Removes special characters.
	clean(&doc)

	// Invalidates documents with all
	// digits equal.
	if allEq(doc) {
		return false
	}

	// Calculates the first digit.
	d := doc[:size]
	digit := calculateDigit(d, position)

	// Calculates the second digit.
	d = d + digit
	digit = calculateDigit(d, position+1)

	return doc == d+digit
}

//Função auxiliar da regra customizada de CPF
func clean(doc *string) {

	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}

//Função auxiliar da regra customizada de CPF
func allEq(doc string) bool {

	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

//Função auxiliar da regra customizada de CPF
func calculateDigit(doc string, position int) string {

	var sum int
	for _, r := range doc {

		sum += int(r-'0') * position
		position--

		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}

	return strconv.Itoa(11 - sum)
}
