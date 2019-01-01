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

#include <memory>
#include <cstring>
#include <iostream>
#include <sstream>
#include <atomic>
#include <thread>
#include "dragonboat/dragonboat.h"

const uint64_t defaultClusterID = 128;

const char addresses[3][16] =
{
  "localhost:63001",
  "localhost:63002",
  "localhost:63003",
};

int main(int argc, char **argv, char **env)
{
  uint64_t nodeID;
  if (argc == 2) {
    nodeID = std::stoi(argv[1]);
    if (nodeID > 3 || nodeID < 1)
    {
      std::cerr << "invalid node id" << std::endl;
      return -1;
    }
  } else {
    std::cerr << "Usage: cpphelloworld nodeID" << std::endl;
    return -1;
  }

  dragonboat::Config config(defaultClusterID, nodeID);
  config.ElectionRTT = 5;
  config.HeartbeatRTT = 1;
  config.CheckQuorum = true;
  config.SnapshotEntries = 10;
  config.CompactionOverhead = 5;
  dragonboat::Peers peers;
  for (uint64_t idx = 0; idx < 3; idx++)
  {
    peers.AddMember(addresses[idx], idx+1);
  }
  std::stringstream path;
  path << "example-data/helloworld-data/node" << nodeID;
  dragonboat::NodeHostConfig nhc(path.str(), path.str());
  nhc.RTTMillisecond = dragonboat::Milliseconds(200);
  nhc.RaftAddress = addresses[nodeID - 1];
  nhc.APIAddress = "";
  dragonboat::Status status;
  auto nh = std::unique_ptr<dragonboat::NodeHost>(
    new dragonboat::NodeHost(nhc));
  status = nh->StartCluster(peers,
    false, "dragonboat-cpp-plugin-cpphelloworld.so", config);
  if (!status.OK())
  {
    std::cerr << "failed to add cluster" << std::endl;
    return -1;
  }
  std::atomic<bool> readyToExit(false);
  auto timeout = dragonboat::Milliseconds(3000);
  auto readThread = std::thread([&nh, &readyToExit, &timeout]() {
    size_t count = 0;
    dragonboat::Buffer query(1);
    dragonboat::Buffer result(sizeof(int));
    while (true)
    {
      if (count == 100)
      {
        count = 0;
        dragonboat::Status rs =
          nh->SyncRead(defaultClusterID, query, &result, timeout);
        if (rs.OK())
        {
          dragonboat::Byte *c = const_cast<dragonboat::Byte *>(result.Data());
          std::cout << "count: " <<
            *(reinterpret_cast<int *>(c)) << std::endl;
        }
      } else {
        count++;
        std::this_thread::sleep_for(std::chrono::milliseconds(100));
      }
      if (readyToExit.load()) {
        return;
      }
    }
  });
  std::unique_ptr<dragonboat::Session> cs(nh->GetNoOPSession(defaultClusterID));
  for (std::string message; std::getline(std::cin, message);)
  {
    if (message == "exit")
    {
      readyToExit = true;
      break;
    }
    dragonboat::UpdateResult result;
    dragonboat::Buffer buf(
      (const dragonboat::Byte *)message.c_str(), message.size());
    status = nh->SyncPropose(cs.get(), buf, timeout, &result);
    if (!status.OK())
    {
      std::cerr << "make proposal failed, error code " <<
        status.Code() << std::endl;
    }
  }
  readThread.join();
  nh->Stop();
  return 0;
}
