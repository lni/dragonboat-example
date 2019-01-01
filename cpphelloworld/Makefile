OS := $(shell uname)
DRAGONBOAT_BINDING_PATH ?= /usr/local
HELLOWORLD_SRC=helloworld.cpp
HELLOWORLD_OBJS=$(subst .cpp,.o,$(HELLOWORLD_SRC))
PLUGIN_SRC=statemachine.cpp plugin.cpp
PLUGIN_OBJS=$(subst .cpp,.o,$(PLUGIN_SRC))
HELLOWORLD_BIN=cpphelloworld
PLUGIN_BIN=dragonboat-cpp-plugin-cpphelloworld.so
CXXFLAGS=-std=c++11 -Wall -fPIC -O3 -I$(DRAGONBOAT_BINDING_PATH)/include

ifeq ($(OS),Darwin)
CPPDS_LDFLAGS=-bundle -undefined dynamic_lookup -Wl,-install_name,$(PLUGIN_BIN)
else ifeq ($(OS),Linux)
CPPDS_LDFLAGS=-shared -Wl,-soname,$(PLUGIN_BIN)
else
$(error OS type $(OS) not supported)
endif

all: $(HELLOWORLD_BIN) $(PLUGIN_BIN)

%.o:%.cpp
	$(CXX) $(CXXFLAGS) -c -o $@ $<

$(HELLOWORLD_BIN): $(HELLOWORLD_OBJS)
	$(CXX) -o $@ $(HELLOWORLD_OBJS) -L$(DRAGONBOAT_BINDING_PATH)/lib -ldragonboatcpp \
		-L$(DRAGONBOAT_BINDING_PATH)/lib -ldragonboat -lpthread -Wl,-rpath,$(DRAGONBOAT_BINDING_PATH)/lib

$(PLUGIN_BIN): $(PLUGIN_OBJS)
	$(CXX) $(CPPDS_LDFLAGS) -o $@ $(PLUGIN_OBJS)

clean:
	rm -f $(HELLOWORLD_BIN) $(HELLOWORLD_OBJS) $(PLUGIN_BIN) $(PLUGIN_OBJS)
