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

const (
	AttributeID        = "id"
	AttrinuteName      = "name"
	AttrinuteCPF       = "cpf"
	AttributeBallance  = "ballance"
	AttributeCreatedAt = "created_at"
	ValidationRequired = "validation.required"
	ValidationFormat   = "validation.format_invalid"
	ValidationAlpha    = "validation.alpha_invalid"
	ValidationFloat    = "validation.float_invalid"
	ValidationNumeric  = "validation.numeric_invalid"
	ValidationDate     = "validation.date_invalid"
)

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

func ValidateJsonRequest(w http.ResponseWriter, r *http.Request, rules map[string][]string, messages govalidator.MapData) *types.Account {
	var account types.Account
	/*opts := govalidator.Options{
		Request:         r,        // request object
		Rules:           rules,    // rules map
		Messages:        messages, // custom message map (Optional)
		RequiredDefault: true,     // all the field to be pass the rules
	}*/

	opts := govalidator.Options{
		Request: r,
		Data:    &account,
		Rules:   rules,
		Messages:messages,
		RequiredDefault:true,
	}

	v := govalidator.New(opts)
	e := v.ValidateJSON()

	if len(e) > 0 {
		//_ := map[string]interface{}{"validationError": e}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(types.ErrorResponse{Error:e})
		return nil
	}

	return &account
}

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

// ValidateCPFFormat verifies if the CPF has a
// valid format.
func validateCPFFormat(doc string) bool {

	const (
		pattern = `^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`
	)

	return regexp.MustCompile(pattern).MatchString(doc)
}

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

// clean removes every rune that is not a digit.
func clean(doc *string) {

	buf := bytes.NewBufferString("")
	for _, r := range *doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	*doc = buf.String()
}

// allEq checks if every rune in a given string
// is equal.
func allEq(doc string) bool {

	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

// calculateDigit calculates the next digit for
// the given document.
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