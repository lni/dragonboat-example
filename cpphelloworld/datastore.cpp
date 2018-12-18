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

#include <cstdio>
#include <cstddef>
#include <cstring>
#include <unistd.h>
#include <iostream>
#include "datastore.h"

HelloWorldDataStore::HelloWorldDataStore(uint64_t clusterID,
  uint64_t nodeID) noexcept
  : dragonboat::DataStore(clusterID, nodeID), update_count_(0)
{
}

HelloWorldDataStore::~HelloWorldDataStore()
{
}

uint64_t HelloWorldDataStore::update(const dragonboat::Byte *data,
  size_t sz) noexcept
{
  char *c = reinterpret_cast<char*>(const_cast<dragonboat::Byte*>(data));
  std::cout << "message: " << std::string(c, sz) << std::endl;
  update_count_++;
  return update_count_;
}

LookupResult HelloWorldDataStore::lookup(const dragonboat::Byte *data,
  size_t sz) const noexcept
{
  // return the update_count_ value
  LookupResult r;
  r.result = new char[sizeof(int)];
  r.size = sizeof(int);
  std::memcpy(r.result, &update_count_, sizeof(int));
  return r;
}

uint64_t HelloWorldDataStore::getHash() const noexcept
{
  return (uint64_t)update_count_;
}

SnapshotResult HelloWorldDataStore::saveSnapshot(
  dragonboat::SnapshotWriter *writer,
  dragonboat::SnapshotFileCollection *collection,
  const dragonboat::DoneChan &done) const noexcept
{
  SnapshotResult r;
  dragonboat::IOResult ret;
  r.error = SNAPSHOT_OK;
  r.size = 0;
  ret = writer->Write((dragonboat::Byte *)&update_count_, sizeof(int));
  if (ret.size != sizeof(int)) {
    r.error = FAILED_TO_SAVE_SNAPSHOT;
    return r;
  }
  r.size = sizeof(int);
  return r;
}

int HelloWorldDataStore::recoverFromSnapshot(dragonboat::SnapshotReader *reader,
  const std::vector<dragonboat::SnapshotFile> &files,
  const dragonboat::DoneChan &done) noexcept
{
  dragonboat::IOResult ret;
  dragonboat::Byte data[sizeof(int)];
  ret = reader->Read(data, sizeof(int));
  if (ret.size != sizeof(int)) {
    return FAILED_TO_RECOVER_FROM_SNAPSHOT;
  }
  ::memcpy(&update_count_, data, sizeof(int));
  return SNAPSHOT_OK;
}

void HelloWorldDataStore::freeLookupResult(LookupResult r) noexcept
{
  delete[] r.result;
}
