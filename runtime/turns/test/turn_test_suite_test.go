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

package test_test

import (
	"testing"

	"github.com/prolang/drydock/runtime/turns/async"
	"github.com/prolang/drydock/runtime/turns/test"
)

// TurnSuiteSuite is the test suite for Manager
type TurnSuiteSuite struct {
	test.Suite
}

// TestTurnSuiteSuite runs the test suite for Manager
func TestTurnSuiteSuite(t *testing.T) {
	test.RunSuite(t, new(TurnSuiteSuite))
}

func (t *TurnSuiteSuite) VerifySyncMethodsAllowed() {
}

func (t *TurnSuiteSuite) VerifyAsyncMethodsAllowed() async.R {
	return async.Done()
}

func (t *TurnSuiteSuite) verifyPrivateMethodsSkipped() {
	t.Fatalf("This should never run.")
}

func (t *TurnSuiteSuite) VerifyMethodsWrongSignatureSkipped1(notRunner int) async.R {
	t.Fatalf("This should never run.")
	return async.R{}
}

func (t *TurnSuiteSuite) VerifyMethodsWrongSignatureSkipped2() int {
	t.Fatalf("This should never run.")
	return 0
}
