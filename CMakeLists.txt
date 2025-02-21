cmake_minimum_required(VERSION 3.30)
project(lang)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

include(FetchContent)
FetchContent_Declare(
        googletest
        URL https://github.com/google/googletest/archive/6910c9d9165801d8827d628cb72eb7ea9dd538c5.zip
)

set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

add_library(lang STATIC library.cpp)

enable_testing()

add_executable(
        library_test
        library_test.cpp
)
target_link_libraries(
        library_test
        GTest::gtest_main
)

include(GoogleTest)
gtest_discover_tests(library_test)
