cmake_minimum_required(VERSION 3.20)

project(jSDR VERSION 0.1
        LANGUAGES CXX
)

list(APPEND CMAKE_MODULE_PATH "${CMAKE_SOURCE_DIR}/cmake")

set(CMAKE_CXX_STANDARD 20)

add_subdirectory(jSDR)