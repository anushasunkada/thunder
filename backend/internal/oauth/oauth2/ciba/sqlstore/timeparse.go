/*
 * Copyright (c) 2026, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package sqlstore

import (
	"fmt"
	"strings"
	"time"
)

// parseTimeField parses a time field from the database result. It preserves any timezone offset
// present in the stored string so a value written in one zone reads back as the same instant.
func parseTimeField(field interface{}, fieldName string) (time.Time, error) {
	switch v := field.(type) {
	case string:
		date, offset := splitTimeAndOffset(v)
		if offset != "" {
			if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999 -0700", date+" "+offset); err == nil {
				return parsedTime, nil
			}
		}
		// No offset present: treat the wall-clock value as UTC (write side normalizes to UTC).
		if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999", date); err == nil {
			return parsedTime.UTC(), nil
		}
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z07:00", v)
		if err != nil {
			return time.Time{}, fmt.Errorf("error parsing %s: %w", fieldName, err)
		}
		return parsedTime, nil
	case time.Time:
		return v, nil
	default:
		return time.Time{}, fmt.Errorf("unexpected type for %s", fieldName)
	}
}

// splitTimeAndOffset splits a database time string into its "date time" portion and timezone
// offset token (e.g. "+0530"), if present. Go's time.Time.String() renders values such as
// "2026-06-02 21:57:49.157215 +0530 +0530 m=+595..."; the third space-separated token is the
// numeric offset that must be retained to read the value back as the same instant.
func splitTimeAndOffset(timeStr string) (date, offset string) {
	parts := strings.SplitN(timeStr, " ", 4)
	if len(parts) < 2 {
		return timeStr, ""
	}
	date = parts[0] + " " + parts[1]
	if len(parts) >= 3 && isNumericOffset(parts[2]) {
		offset = parts[2]
	}
	return date, offset
}

// isNumericOffset reports whether the token is a numeric timezone offset like "+0530" or "-0700".
func isNumericOffset(token string) bool {
	if len(token) != 5 || (token[0] != '+' && token[0] != '-') {
		return false
	}
	for _, c := range token[1:] {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
