// Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/changaolee/skeleton.

package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"
)

const (
	qnameCharFmt     = "[A-Za-z0-9]"
	qnameExtCharFmt  = "[-A-Za-z0-9_.]"
	qualifiedNameFmt = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
)

const (
	qualifiedNameErrMsg    = "must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
	qualifiedNameMaxLength = 63
)

var qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

// IsQualifiedName 检查 value 是否是一个合法的名字.
func IsQualifiedName(value string) []string {
	var errs []string
	parts := strings.Split(value, "/")
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			errs = append(errs, "prefix part "+EmptyError())
		} else if msgs := IsDNS1123Subdomain(prefix); len(msgs) != 0 {
			errs = append(errs, prefixEach(msgs, "prefix part ")...)
		}
	default:
		return append(
			errs,
			"a qualified name "+RegexError(
				qualifiedNameErrMsg,
				qualifiedNameFmt,
				"MyName",
				"my.name",
				"123-abc",
			)+" with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')",
		)
	}

	if len(name) == 0 {
		errs = append(errs, "name part "+EmptyError())
	} else if len(name) > qualifiedNameMaxLength {
		errs = append(errs, "name part "+MaxLenError(qualifiedNameMaxLength))
	}
	if !qualifiedNameRegexp.MatchString(name) {
		errs = append(
			errs,
			"name part "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"),
		)
	}
	return errs
}

const labelValueFmt string = "(" + qualifiedNameFmt + ")?"

const labelValueErrMsg string = "a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"

// LabelValueMaxLength 是一个标签的最大长度.
const LabelValueMaxLength int = 63

var labelValueRegexp = regexp.MustCompile("^" + labelValueFmt + "$")

// IsValidLabelValue 检查 value 是否是一个合法的标签.
func IsValidLabelValue(value string) []string {
	var errs []string
	if len(value) > LabelValueMaxLength {
		errs = append(errs, MaxLenError(LabelValueMaxLength))
	}
	if !labelValueRegexp.MatchString(value) {
		errs = append(errs, RegexError(labelValueErrMsg, labelValueFmt, "MyValue", "my_value", "12345"))
	}
	return errs
}

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

const dns1123LabelErrMsg string = "a DNS-1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character"

// DNS1123LabelMaxLength 是 DNS（RFC 1123）中标签的最大长度.
const DNS1123LabelMaxLength int = 63

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

// IsDNS1123Label 检查 value 是否是符合 DNS（RFC 1123）中标签定义的字符串.
func IsDNS1123Label(value string) []string {
	var errs []string
	if len(value) > DNS1123LabelMaxLength {
		errs = append(errs, MaxLenError(DNS1123LabelMaxLength))
	}
	if !dns1123LabelRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123LabelErrMsg, dns1123LabelFmt, "my-name", "123-abc"))
	}
	return errs
}

const dns1123SubdomainFmt = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"

const dns1123SubdomainErrorMsg = "a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and   must start and end with an alphanumeric character"

// DNS1123SubdomainMaxLength 是 DNS（RFC 1123）中子域的最大长度.
const DNS1123SubdomainMaxLength int = 253

var dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")

// IsDNS1123Subdomain 检查 value 是否是符合 DNS（RFC 1123）中子域定义的字符串.
func IsDNS1123Subdomain(value string) []string {
	var errs []string
	if len(value) > DNS1123SubdomainMaxLength {
		errs = append(errs, MaxLenError(DNS1123SubdomainMaxLength))
	}
	if !dns1123SubdomainRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123SubdomainErrorMsg, dns1123SubdomainFmt, "example.com"))
	}
	return errs
}

// IsValidPortNum 检查 value 是否是非 0 端口号.
func IsValidPortNum(port int) []string {
	if 1 <= port && port <= 65535 {
		return nil
	}
	return []string{InclusiveRangeError(1, 65535)}
}

// IsInRange 检查 value 是否在指定范围内.
func IsInRange(value int, min int, max int) []string {
	if value >= min && value <= max {
		return nil
	}
	return []string{InclusiveRangeError(min, max)}
}

// IsValidIP 检查 value 是否是可用 IP 地址.
func IsValidIP(value string) []string {
	if net.ParseIP(value) == nil {
		return []string{"must be a valid IP address, (e.g. 10.9.8.7)"}
	}
	return nil
}

const (
	percentFmt    string = "[0-9]+%"
	percentErrMsg string = "a valid percent string must be a numeric string followed by an ending '%'"
)

var percentRegexp = regexp.MustCompile("^" + percentFmt + "$")

// IsValidPercent 检查字符串是否是百分比形式.
func IsValidPercent(percent string) []string {
	if !percentRegexp.MatchString(percent) {
		return []string{RegexError(percentErrMsg, percentFmt, "1%", "93%")}
	}
	return nil
}

// MaxLenError 字符串过长错误.
func MaxLenError(length int) string {
	return fmt.Sprintf("must be no more than %d characters", length)
}

// RegexError 正则不匹配错误.
func RegexError(msg string, fmt string, examples ...string) string {
	if len(examples) == 0 {
		return msg + " (regex used for validation is '" + fmt + "')"
	}
	msg += " (e.g. "
	for i := range examples {
		if i > 0 {
			msg += " or "
		}
		msg += "'" + examples[i] + "', "
	}
	msg += "regex used for validation is '" + fmt + "')"
	return msg
}

// EmptyError 空值错误.
func EmptyError() string {
	return "must be non-empty"
}

func prefixEach(msgs []string, prefix string) []string {
	for i := range msgs {
		msgs[i] = prefix + msgs[i]
	}
	return msgs
}

// InclusiveRangeError 范围错误.
func InclusiveRangeError(lo, hi int) string {
	return fmt.Sprintf(`must be between %d and %d, inclusive`, lo, hi)
}

const (
	minPassLength = 8
	maxPassLength = 16
)

// IsValidPassword 检查 password 是否是一个合法密码.
func IsValidPassword(password string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
			passLen++
		case unicode.IsUpper(ch):
			hasUpper = true
			passLen++
		case unicode.IsLower(ch):
			hasLower = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !hasLower {
		appendError("lowercase letter missing")
	}
	if !hasUpper {
		appendError("uppercase letter missing")
	}
	if !hasNumber {
		appendError("at least one numeric character required")
	}
	if !hasSpecial {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(
			fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength),
		)
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}

	return nil
}
