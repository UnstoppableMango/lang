cmake_minimum_required(VERSION 3.10)

project(MangoLang)

find_package(llvm CONFIG REQUIRED)

add_executable(Kaleidoscope kaleidoscope.cpp)

target_link_libraries(MangoLang PRIVATE llvm::llvm)
