// Copyright 2015 Google Inc. All Rights Reserved.
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

package gcs

import "io"

// A request to create an object, accepted by Bucket.CreateObject.
type CreateObjectRequest struct {
	// The name with which to create the object. This field must be set.
	//
	// Object names must:
	//
	// *  be non-empty.
	// *  be no longer than 1024 bytes.
	// *  be valid UTF-8.
	// *  not contain the code point U+000A (line feed).
	// *  not contain the code point U+000D (carriage return).
	//
	// See here for authoritative documentation:
	//     https://cloud.google.com/storage/docs/bucket-naming#objectnames
	Name string

	// Optional information with which to create the object. See here for more
	// information:
	//
	//     https://cloud.google.com/storage/docs/json_api/v1/objects#resource
	//
	ContentType     string
	ContentLanguage string
	ContentEncoding string
	CacheControl    string
	Metadata        map[string]string

	// A reader from which to obtain the contents of the object. Must be non-nil.
	Contents io.Reader

	// If non-nil, the object will be created/overwritten only if the current
	// generation for the object name is equal to the given value. Zero means the
	// object does not exist.
	GenerationPrecondition *int64
}

// A request to read the contents of an object at a particular generation.
type ReadObjectRequest struct {
	// The name of the object to read.
	Name string

	// The generation of the object to read. Zero means the latest generation.
	Generation int64
}

type StatObjectRequest struct {
	// The name of the object in question.
	Name string
}

type ListObjectsRequest struct {
	// List only objects whose names begin with this prefix.
	Prefix string

	// Collapse results based on a delimiter.
	//
	// If non-empty, enable the following behavior. For each run of one or more
	// objects whose names are of the form:
	//
	//     <Prefix><S><Delimiter><...>
	//
	// where <S> is a string that doesn't itself contain Delimiter and <...> is
	// anything, return a single Collaped entry in the listing consisting of
	//
	//     <Prefix><S><Delimiter>
	//
	// instead of one Object record per object. If a collapsed entry consists of
	// a large number of objects, this may be more efficient.
	Delimiter string

	// Used to continue a listing where a previous one left off. See
	// Listing.ContinuationToken for more information.
	ContinuationToken string

	// The maximum number of objects and collapsed runs to return. Fewer than
	// this number may actually be returned. If this is zero, a sensible default
	// is used.
	MaxResults int
}

// A set of objects and delimter-based collapsed runs returned by a call to
// ListObjects. See also ListObjectsRequest.
type Listing struct {
	// Records for objects matching the listing criteria.
	//
	// Guaranteed to be strictly increasing under a lexicographical comparison on
	// (name, generation) pairs.
	Objects []*Object

	// Collapsed entries for runs of names sharing a prefix followed by a
	// delimiter. See notes on ListObjectsRequest.Delimiter.
	//
	// Guaranteed to be strictly increasing.
	CollapsedRuns []string

	// A continuation token, for fetching more results.
	//
	// If non-empty, this listing does not represent the full set of matching
	// objects in the bucket. Call ListObjects again with the request's
	// ContinuationToken field set to this value to continue where you left off.
	//
	// Guarantees, for replies R1 and R2, with R2 continuing from R1:
	//
	//  *  All of R1's object names are strictly less than all object names and
	//     collapsed runs in R2.
	//
	//  *  All of R1's collapsed runs are strictly less than all object names and
	//     prefixes in R2.
	//
	// (Cf. Google-internal bug 19286144)
	//
	// Note that there is no guarantee of atomicity of listings. Objects written
	// and deleted concurrently with a single or multiple listing requests may or
	// may not be returned.
	ContinuationToken string
}

// A request to update the metadata of an object, accepted by
// Bucket.UpdateObject.
type UpdateObjectRequest struct {
	// The name of the object to update. Must be specified.
	Name string

	// String fields in the object to update (or not). The semantics are as
	// follows, for a given field F:
	//
	//  *  If F is set to nil, the corresponding GCS object field is untouched.
	//
	//  *  If *F is the empty string, then the corresponding GCS object field is
	//     removed.
	//
	//  *  Otherwise, the corresponding GCS object field is set to *F.
	//
	//  *  There is no facility for setting a GCS object field to the empty
	//     string, since many of the fields do not actually allow that as a legal
	//     value.
	//
	// Note that the GCS object's content type field cannot be removed.
	ContentType     *string
	ContentEncoding *string
	ContentLanguage *string
	CacheControl    *string

	// User-provided metadata updates. Keys that are not mentioned are untouched.
	// Keys whose values are nil are deleted, and others are updated to the
	// supplied string. There is no facility for completely removing user
	// metadata.
	Metadata map[string]*string
}