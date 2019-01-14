// Copyright 2014 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.
//
//
// Copyright 2017,2018 Lei Ni (nilei81@gmail.com).
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

/*
Package utils contains helper functions and structs used by dragonboat's
examples.
*/
package utils

import (
	"sync"
)

// Stopper is a manager struct for managing worker goroutines.
type Stopper struct {
	shouldStopC chan struct{}
	wg          sync.WaitGroup
}

// NewStopper return a new Stopper instance.
func NewStopper() *Stopper {
	s := &Stopper{
		shouldStopC: make(chan struct{}),
	}

	return s
}

// RunWorker creates a new goroutine and invoke the f func in that new
// worker goroutine.
func (s *Stopper) RunWorker(f func()) {
	s.wg.Add(1)

	go func() {
		f()
		s.wg.Done()
	}()
}

// ShouldStop returns a chan struct{} used for indicating whether the
// Stop() function has been called on Stopper.
func (s *Stopper) ShouldStop() chan struct{} {
	return s.shouldStopC
}

// Wait waits on the internal sync.WaitGroup. It only return when all
// managed worker goroutines are ready to return and called
// sync.WaitGroup.Done() on the internal sync.WaitGroup.
func (s *Stopper) Wait() {
	s.wg.Wait()
}

// Stop signals all managed worker goroutines to stop and wait for them
// to actually stop.
func (s *Stopper) Stop() {
	close(s.shouldStopC)
	s.wg.Wait()
}

// Close closes the internal shouldStopc chan struct{} to signal all
// worker goroutines that they should stop.
func (s *Stopper) Close() {
	close(s.shouldStopC)
}
