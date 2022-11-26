// Copyright 2022 VMware, Inc.
// SPDX-License-Identifier: Apache-2.0

package set

import "strings"

type Set struct {
	M map[string]bool
}

func New() Set {
	s := Set{
		M: make(map[string]bool),
	}

	return s
}

func (s *Set) Add(value string) {
	s.M[value] = true
}

func (s *Set) AddAll(t string) {
	units := strings.Split(t, ",")

	for _, unit := range units {
		s.M[unit] = true
	}
}

func (s *Set) Remove(value string) {
	delete(s.M, value)
}

func (s *Set) Contains(value string) bool {
	_, c := s.M[value]

	return c
}
