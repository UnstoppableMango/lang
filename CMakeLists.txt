cmake_minimum_required(VERSION 3.16)
project(MangoLang)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
FetchContent_Declare(
  gRPC
  GIT_REPOSITORY https://github.com/grpc/grpc
  GIT_TAG        v1.71.0
)
set(FETCHCONTENT_QUIET OFF)
FetchContent_MakeAvailable(gRPC)

add_executable(lang main.cpp)

target_link_libraries(lang grpc++)
