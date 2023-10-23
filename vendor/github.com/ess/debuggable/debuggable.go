// package debuggable is a quick-and-dirty shortcut for taking actions in an
// app if the DEBUG environment variable is present.
package debuggable

import (
	"os"
)

const (
	envVar = "DEBUG"
)

// Enabled returns true if debugging is enabled and false otherwise.
func Enabled() bool {
	_, present := os.LookupEnv(envVar)

	return present
}

// Copyright © 2019 Dennis Walters
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
