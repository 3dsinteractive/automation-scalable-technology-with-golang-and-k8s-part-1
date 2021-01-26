// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

var specialCharsRegex = regexp.MustCompile("[^a-zA-Z0-9]+")

func randString() string {
	i := rand.Int()
	return fmt.Sprintf("%d", i)
}

func escapeName(tokens ...string) string {
	// Any name rules
	// - Lowercase only (for consistency)
	// - . (dot), _ (underscore), - (minus) can be used
	// - Max length = 250
	var b bytes.Buffer

	// Name result must be token1-token2-token3-token4 without special characters
	for i, token := range tokens {
		if len(token) == 0 {
			continue
		}

		token = strings.ToLower(token)

		cleanToken := specialCharsRegex.ReplaceAllString(token, "-")
		if i != 0 {
			b.WriteString("-")
		}
		b.WriteString(cleanToken)
	}

	name := b.String()
	// - Cannot start with -, _, +
	for true {
		if len(name) == 0 || name[0] != '-' {
			break
		}
		name = name[1:]
	}

	// - Cannot be longer than 250 characters (max len)
	if len(name) > 250 {
		name = name[0:250]
	}

	return name
}
