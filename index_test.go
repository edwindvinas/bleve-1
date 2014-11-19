//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package bleve

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIndexCreateNewOverExisting(t *testing.T) {
	defer os.RemoveAll("testidx")

	index, err := New("testidx", NewIndexMapping())
	if err != nil {
		t.Fatal(err)
	}
	index.Close()
	index, err = New("testidx", NewIndexMapping())
	if err != ErrorIndexPathExists {
		t.Fatalf("expected error index path exists, got %v", err)
	}
}

func TestIndexOpenNonExisting(t *testing.T) {
	_, err := Open("doesnotexist")
	if err != ErrorIndexPathDoesNotExist {
		t.Fatalf("expected error index path does not exist, got %v", err)
	}
}

func TestIndexOpenMetaMissingOrCorrupt(t *testing.T) {
	defer os.RemoveAll("testidx")

	index, err := New("testidx", NewIndexMapping())
	if err != nil {
		t.Fatal(err)
	}
	index.Close()

	// now intentionally corrupt the metadata
	ioutil.WriteFile("testidx/index_meta.json", []byte("corrupted"), 0666)

	index, err = Open("testidx")
	if err != ErrorIndexMetaCorrupt {
		t.Fatalf("expected error index metadata corrupted, got %v", err)
	}

	// no intentionally remove the metadata
	os.Remove("testidx/index_meta.json")

	index, err = Open("testidx")
	if err != ErrorIndexMetaMissing {
		t.Fatalf("expected error index metadata missing, got %v", err)
	}
}
