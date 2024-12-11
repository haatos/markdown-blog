package data

import (
	"fmt"
	"strings"
	"testing"
)

func TestUpdateQuery(t *testing.T) {
	tests := []struct {
		table      string
		attributes []string
		want       string
	}{
		{table: "users", attributes: []string{"name"}, want: "update users set name = $1"},
		{table: "users", attributes: []string{"name", "age"}, want: "update users set name = $1, age = $2"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("table:%s,attributes:%s", tt.table, strings.Join(tt.attributes, ";"))
		t.Run(testname, func(t *testing.T) {
			q := UpdateQuery(tt.table, tt.attributes...)
			if q != tt.want {
				t.Errorf("got %s, want %s", q, tt.want)
			}
		})
	}
}

func TestUpdateIDQuery(t *testing.T) {
	tests := []struct {
		table      string
		attributes []string
		want       string
	}{
		{table: "users", attributes: []string{"name"}, want: "update users set name = $1 where id = $2"},
		{table: "users", attributes: []string{"name", "age"}, want: "update users set name = $1, age = $2 where id = $3"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("table:%s,attributes:%s", tt.table, strings.Join(tt.attributes, ";"))
		t.Run(testname, func(t *testing.T) {
			q := UpdateIDQuery(tt.table, tt.attributes...)
			if q != tt.want {
				t.Errorf("got %s, want %s", q, tt.want)
			}
		})
	}
}
