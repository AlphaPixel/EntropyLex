cmake_minimum_required(VERSION 3.23)

project(entropylex-kowalski LANGUAGES CXX)
include(FetchContent)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)


add_executable(entropylex)


target_sources(entropylex PRIVATE main.cc)
