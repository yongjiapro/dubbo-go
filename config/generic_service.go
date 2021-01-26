/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"context"
	"github.com/apache/dubbo-go/protocol"
)

// GenericService uses for generic invoke for service call
type GenericService struct {
	Invoke       func(ctx context.Context, req []interface{}) (interface{}, error) `dubbo:"$invoke"`
	referenceStr string
}

// NewGenericService returns a GenericService instance
func NewGenericService(referenceStr string) *GenericService {
	return &GenericService{referenceStr: referenceStr}
}

// Reference gets referenceStr from GenericService
func (u *GenericService) Reference() string {
	return u.referenceStr
}

// GenericResultService uses for generic invoke for service call, returns a raw protocol.Result
type GenericResultService struct {
	Invoke       func(ctx context.Context, args []interface{}) protocol.Result `dubbo:"$invoke"`
	referenceStr string
}

// NewGenericResultService returns a GenericService instance
func NewGenericResultService(referenceStr string) *GenericResultService {
	return &GenericResultService{referenceStr: referenceStr}
}

// Reference gets referenceStr from GenericService
func (u *GenericResultService) Reference() string {
	return u.referenceStr
}
