package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	runes := []rune(str)
	length := len(runes)
	const slash = '\\'

	if length == 0 {
		return "", nil
	}

	sb := strings.Builder{}

	isEscape := false
	for i := 0; i < length; i++ {
		// получение текущей руны и установка флагов что она может быть \ либо цифрой
		current := runes[i]

		isCurrentDigit := false
		_, err := strconv.Atoi(string(current))
		if err == nil {
			isCurrentDigit = true
		}

		isCurrentSlash := false
		if current == slash {
			isCurrentSlash = true
		}

		// если текущая руна число, а на предыдущей итерации не было экранирования - ошибка
		if isCurrentDigit && !isEscape {
			return "", ErrInvalidString
		}

		// получение следующей руны, если возможно, и установка флагов, что она может быть \ либо цифрой
		var next rune
		if i+1 < length {
			next = runes[i+1]
		}

		isNextDigit := false
		nextDigit, err := strconv.Atoi(string(next))
		if err == nil {
			isNextDigit = true
		}

		isNextSlash := false
		if next == slash {
			isNextSlash = true
		}

		// если текущая руна \ и она не экранирована на предыдущей итерации и следующая не цифра или слеш - ошибка
		if isCurrentSlash && !isEscape {
			if !isNextSlash && !isNextDigit {
				return "", ErrInvalidString
			}
			// ставим флаг, что происходит экранирование и переходим к следующей итерации
			isEscape = true
			continue
		}

		isEscape = false
		if isNextDigit {
			// Если следующая руна цифра и она больше 0 - добавляем к строке полученную подстроку
			// и передвигаемся на итерацию вперёд, чтобы "перепрыгнуть" обработанную цифру
			if nextDigit != 0 {
				sb.WriteString(strings.Repeat(string(current), nextDigit))
			}
			i++
		} else {
			sb.WriteRune(current)
		}
	}

	return sb.String(), nil
}
