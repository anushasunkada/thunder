/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
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

package executor

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BuiltInExecutorCatalogTestSuite struct {
	suite.Suite
}

func TestBuiltInExecutorCatalogSuite(t *testing.T) {
	suite.Run(t, new(BuiltInExecutorCatalogTestSuite))
}

func (suite *BuiltInExecutorCatalogTestSuite) executorNameConstantsFromFile() []string {
	suite.T().Helper()

	_, thisFile, _, ok := runtime.Caller(0)
	require.True(suite.T(), ok)

	constantsPath := filepath.Join(filepath.Dir(thisFile), "constants.go")
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, constantsPath, nil, 0)
	require.NoError(suite.T(), err)

	var names []string
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.CONST {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, ident := range valueSpec.Names {
				if ident.Name == "" || ident.Name[0] == '_' {
					continue
				}
				if len(ident.Name) < len("ExecutorName") || ident.Name[:len("ExecutorName")] != "ExecutorName" {
					continue
				}
				if i >= len(valueSpec.Values) {
					continue
				}
				lit, ok := valueSpec.Values[i].(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				name, err := strconv.Unquote(lit.Value)
				require.NoError(suite.T(), err)
				names = append(names, name)
			}
		}
	}
	return names
}

func (suite *BuiltInExecutorCatalogTestSuite) TestBuiltInExecutorCatalog_MatchesExecutorNameConstants() {
	constants := suite.executorNameConstantsFromFile()
	assert.ElementsMatch(suite.T(), constants, builtInExecutorNames)
	assert.ElementsMatch(suite.T(), constants, defaultBuiltInExecutorNames())
}

func (suite *BuiltInExecutorCatalogTestSuite) TestValidateBuiltInExecutorCatalog() {
	require.NoError(suite.T(), validateBuiltInExecutorCatalog())
}
