/*
|    Protect your secrets, protect your sensitive data.
:    Explore VMware Secrets Manager docs at https://vsecm.com/
</
<>/  keep your secrets... secret
>/
<>/' Copyright 2023-present VMware Secrets Manager contributors.
>/'  SPDX-License-Identifier: BSD-2-Clause
*/

package data

import (
	"fmt"
	"time"
)

// SecretStored represents a secret stored in VSecM Safe.
type SecretStored struct {
	// Name of the secret.
	Name string `json:"name"`

	// Raw value. A secret can have multiple values. Sentinel returns
	// a single value if there is a single value in this array. Sentinel
	// will return an array of values if there are multiple values in the array.
	Value string `json:"value"`

	// Transformed values. This value is the value that workloads see.
	//
	// Apply transformation (if needed) and then store the value in
	// one of the supported formats. If the format is json, ensure that
	// a valid JSON is stored here. If the format is yaml, ensure that
	// a valid YAML is stored here. If the format is none, then just
	// apply transformation (if needed) and do not do any validity check.
	ValueTransformed string `json:"valueTransformed"`

	// Additional information that helps format and store the secret.
	Meta SecretMeta `json:"meta"`
	// Timestamps
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	// Invalid before this time.
	NotBefore time.Time `json:"notBefore"`
	// Invalid after this time.
	ExpiresAfter time.Time `json:"expiresAfter"`
}

// ToMapForK8s returns a map that can be used to create a Kubernetes secret.
//
//  1. If there is no template, attempt to unmarshal the secret's value
//     into a map. If that fails, store the secret's value under the "VALUE" key.
//  2. If there is a template, attempt to parse it. If parsing is successful,
//     create a new map with the parsed data. If parsing fails, follow the same
//     logic as in case 1, attempting to unmarshal the secret's value into a map,
//     and if that fails, storing the secret's value under the "VALUE" key.
func (secret SecretStored) ToMapForK8s() map[string][]byte {
	// If there is no template, use the secret's value as is.
	if secret.Meta.Template == "" {
		return convertValueNoTemplate(secret.Value)
	}

	// Otherwise, apply the template.
	newData, err := parseForK8sSecret(secret)
	if err == nil {
		return convertMapToStringBytes(newData)
	}

	// If the template fails, use the secret's value as is.
	return convertValueToMap(secret.Value)
}

// ToMap converts the SecretStored struct to a map[string]any.
// The resulting map contains the following key-value pairs:
//
//	"Name": the Name field of the SecretStored struct
//	"Values": the Values field of the SecretStored struct
//	"Created": the Created field of the SecretStored struct
//	"Updated": the Updated field of the SecretStored struct
func (secret SecretStored) ToMap() map[string]any {
	return map[string]any{
		"Name":    secret.Name,
		"Value":   secret.Value,
		"Created": secret.Created,
		"Updated": secret.Updated,
	}
}

// Parse takes a data.SecretStored type as input and returns the parsed
// string or an error.
//
// It parses all the `.Values` of the secret, and for each value tries to apply
// a template transformation.
//
// Here is how the template transformation is applied:
//
//  1. Compute parsedString:
//     If the Meta.Template field is empty, then parsedString is the original
//     value. Otherwise, parsedString is the result of applying the template
//     transformation to the original value.
//
// 2.	Compute the output string:
//   - If the Meta.Format field is Json, then the output string is parsedString
//     if parsedString is a valid JSON, otherwise it's the original value.
//   - If the Meta.Format field is Yaml, then the output string is the result of
//     transforming parsedString into Yaml if parsedString is a valid JSON,
//     otherwise it's parsedString.
//   - If the Meta.Format field is Raw, then the output string is simply the
//     parsedString, without any specific format checks or transformations.
func (secret SecretStored) Parse() (string, error) {
	if len(secret.Value) == 0 {
		return "", fmt.Errorf("no value found for secret %s", secret.Name)
	}
	transformed, err := transform(secret.Value, secret.Meta.Template,
		secret.Meta.Format)

	if transformed == "" {
		return secret.Value, fmt.Errorf("failed to parse secret %s", secret.Name)
	}

	return transformed, err
}
