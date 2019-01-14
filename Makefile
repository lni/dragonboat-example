# Copyright 2017 Lei Ni (nilei81@gmail.com).
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

OS := $(shell uname)
ROCKSDB_MAJOR_VER=5
ifeq ($(OS),Darwin)
ROCKSDB_SO_FILE=librocksdb.$(ROCKSDB_MAJOR_VER).dylib
else ifeq ($(OS),Linux)
ROCKSDB_SO_FILE=librocksdb.so.$(ROCKSDB_MAJOR_VER)
else
$(error OS type $(OS) not supported)
endif

ROCKSDB_INC_PATH ?=
ROCKSDB_LIB_PATH ?=
# in /usr/local/lib?
ifeq ($(ROCKSDB_LIB_PATH),)
ifeq ($(ROCKSDB_INC_PATH),)
ifneq ($(wildcard /usr/local/lib/$(ROCKSDB_SO_FILE)),)
ifneq ($(wildcard /usr/local/include/rocksdb/c.h),)
$(info rocksdb lib found at /usr/local/lib/$(ROCKSDB_SO_FILE))
ROCKSDB_LIB_PATH=/usr/local/lib
endif
endif
endif
endif

ifeq ($(ROCKSDB_LIB_PATH),)
CDEPS_LDFLAGS=-lrocksdb
else
CDEPS_LDFLAGS=-L$(ROCKSDB_LIB_PATH) -lrocksdb
endif
ifneq ($(ROCKSDB_INC_PATH),)
CGO_CXXFLAGS=CGO_CFLAGS="-I$(ROCKSDB_INC_PATH)"
endif
CGO_LDFLAGS=CGO_LDFLAGS="$(CDEPS_LDFLAGS)"
GOCMD=$(CGO_LDFLAGS) $(CGO_CXXFLAGS) go build -v

all: helloworld

helloworld:
	$(GOCMD) -o example-helloworld github.com/lni/dragonboat-example/helloworld

multigroup:
	$(GOCMD) -o example-multigroup github.com/lni/dragonboat-example/multigroup

cpphelloworld:
	make -C cpphelloworld

clean:
	@make -C cpphelloworld clean
	@rm -f example-helloworld example-multigroup

.PHONY: helloworld multigroup cpphelloworld clean
