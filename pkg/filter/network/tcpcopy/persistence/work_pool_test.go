/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package persistence

import (
	"fmt"
	"mosn.io/layotto/pkg/filter/network/tcpcopy/model"
	_type "mosn.io/layotto/pkg/filter/network/tcpcopy/type"
	"mosn.io/mosn/pkg/log"
	"testing"
	"time"
)

func TestNewDefaultWorkPool(t *testing.T) {

	GetTcpcopyLogger().SetLogLevel(log.DEBUG)
	GetStaticConfLogger().SetLogLevel(log.DEBUG)

	workPool := NewDefaultWorkPool(10)

	go func() {
		for i := 0; i < 5; i++ {
			model_1 := model.NewDumpUploadDynamicConfig("uuid_1", _type.RPC, "12200", nil, "1s")
			model_2 := model.NewDumpUploadDynamicConfig("uuid_2", _type.RPC, "12200", nil, "2s")
			model_3 := model.NewDumpUploadDynamicConfig("uuid_3", _type.RPC, "12200", nil, "3s")
			go func() {
				workPool.Schedule(model_1)
				workPool.Schedule(model_2)
				workPool.Schedule(model_3)
			}()
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			model_4 := model.NewDumpUploadDynamicConfig("uuid_4", _type.RPC, "12200", nil, "4s")
			model_5 := model.NewDumpUploadDynamicConfig("uuid_5", _type.RPC, "12200", nil, "5s")
			model_6 := model.NewDumpUploadDynamicConfig("uuid_6", _type.RPC, "12200", nil, "6s")
			go func() {
				workPool.Schedule(model_4)
				workPool.Schedule(model_5)
				workPool.Schedule(model_6)
			}()
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			model_7 := model.NewDumpUploadDynamicConfig("uuid_7", _type.CONFIGURATION, "12200", nil, "7s")
			model_8 := model.NewDumpUploadDynamicConfig("uuid_8", _type.CONFIGURATION, "12200", nil, "8s")
			model_9 := model.NewDumpUploadDynamicConfig("uuid_9", _type.CONFIGURATION, "12200", nil, "9s")
			go func() {
				workPool.Schedule(model_7)
				workPool.Schedule(model_8)
				workPool.Schedule(model_9)
			}()
		}
	}()

	time.Sleep(3 * time.Second)

	totalTasksCount := 0
	workPool.workers.Range(func(key, value interface{}) bool {
		if value, ok := workPool.workers.Load(key); ok {
			worker := value.(*WorkGoroutine)
			worker.tasks.Range(func(key, value interface{}) bool {
				totalTasksCount++
				return true
			})
		}
		return true
	})

	if totalTasksCount != 0 {
		fmt.Println(totalTasksCount)
		t.Errorf("Test_WorkPool Failed")
	}
}
