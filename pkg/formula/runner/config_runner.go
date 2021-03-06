/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package runner

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/ZupIT/ritchie-cli/pkg/formula"
)

var _ formula.ConfigRunner = ConfigManager{}

const FileName = "default-formula-runner"

var ErrConfigNotFound = errors.New("you must configure your default formula execution method, run \"rit set formula-runner\" to set up")

type ConfigManager struct {
	filePath string
}

func NewConfigManager(ritHome string) ConfigManager {
	return ConfigManager{
		filePath: filepath.Join(ritHome, FileName),
	}
}

func (c ConfigManager) Create(runType formula.RunnerType) error {
	data, err := json.Marshal(runType)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(c.filePath, data, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (c ConfigManager) Find() (formula.RunnerType, error) {
	data, err := ioutil.ReadFile(c.filePath)
	if err != nil {
		return formula.DefaultRun, ErrConfigNotFound
	}

	runType, err := strconv.Atoi(string(data))
	if err != nil {
		return formula.DefaultRun, err
	}

	return formula.RunnerType(runType), nil
}
