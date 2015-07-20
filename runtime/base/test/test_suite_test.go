// Copyright 2015 The Drydock Authors.
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package test

import (
	"strings"
	"testing"
)

// SuiteSuite is the test suite for test.Suite.
type SuiteSuite struct {
	Suite
}

// TestTestSuite runs the test.Suite test suite.
func TestTestSuite(t *testing.T) {
	RunSuite(t, new(SuiteSuite))
}

// GetFileLine verify returns the correct file and line number information.
func (t *SuiteSuite) GetFileLine() {
	expected := "drydock/runtime/base/test/test_suite_test.go:39"
	wrapper := func() string {
		return t.getFileLine()
	}
	if s := wrapper(); !strings.HasSuffix(s, expected) {
		t.Errorf("Invalid file and line: Got: %s, Want: %s", s, expected)
	}
}

// TestInfof ensures that Infof has code coverage.
func (t *SuiteSuite) TestInfof() {
	t.Infof("This is a log statement produced by t.Infof")
}

// VerifyMethodsWrongSignatureSkipped1 ensures that public methods with the wrong signature
// (e.g. take arguments) are not executed as test methods.
func (t *SuiteSuite) VerifyMethodsWrongSignatureSkipped1(x int) {
	t.Fatalf("This should never run.")
}

// VerifyMethodsWrongSignatureSkipped2 ensures that public methods with the wrong signature
// (e.g. have return values) are not executed as test methods.
func (t *SuiteSuite) VerifyMethodsWrongSignatureSkipped2() int {
	t.Fatalf("This should never run.")
	return 0
}
