cmake_minimum_required(VERSION 3.10)

project(MangoLang)

find_package(LLVM CONFIG REQUIRED)

add_executable(lang main.cpp)

list(APPEND CMAKE_MODULE_PATH "${LLVM_CMAKE_DIR}")
include(HandleLLVMOptions)
add_definitions(${LLVM_DEFINITIONS})

target_include_directories(lang PRIVATE ${LLVM_INCLUDE_DIRS})

# Find the libraries that correspond to the LLVM components that we wish to use
llvm_map_components_to_libnames(llvm_libs Support Core IRReader)

# Link against LLVM libraries
target_link_libraries(lang PRIVATE ${llvm_libs})
