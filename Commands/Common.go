package commands

import "strings"

// split to tokens, convert first token to lower case,
// ignore trailing whitespaces
func PrepareMessage(message string) ([]string) {
    tokens := strings.Split(strings.Trim(message, " \n\t"), " ")

    var result []string
    for _, t := range tokens {
        if t != "" {
            result = append(result, t)
        }
    }

    if len(result) == 0 {
        return result
    }

    result[0] = strings.ToLower(result[0])

    return result
}
