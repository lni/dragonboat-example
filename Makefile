# Copyright 2017-2020 Lei Ni (nilei81@gmail.com).
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GOCMD=go build -v

all: helloworld multigroup ondisk optimistic-write-lock

helloworld:
	$(GOCMD) -o example-helloworld github.com/lni/dragonboat-example/v3/helloworld

multigroup:
	$(GOCMD) -o example-multigroup github.com/lni/dragonboat-example/v3/multigroup

ondisk:
	$(GOCMD) -o example-ondisk github.com/lni/dragonboat-example/v3/ondisk

optimistic-write-lock:
	$(GOCMD) -o example-optimistic-write-lock github.com/lni/dragonboat-example/v3/optimistic-write-lock

clean:
	@rm -f example-helloworld example-multigroup example-ondisk example-optimistic-write-lock

.PHONY: helloworld multigroup ondisk optimistic-write-lock clean
