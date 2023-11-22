cmake_minimum_required(VERSION 3.14...3.28)

project(jSDR VERSION 0.1
        LANGUAGES CXX)

include(FetchContent)

FetchContent_Declare(
    wxwidgets
    GIT_REPOSITORY  https://github.com/wxWidgets/wxWidgets.git
    GIT_TAG         master
)

if(NOT wxwidgets_POPULATED)
    FetchContent_Populate(wxwidgets)
    # Store old value of 'BUILD_SHARED_LIBS'
    set(BUILD_SHARED_LIBS_OLD ${BUILD_SHARED_LIBS})
    # Make subproject to use 'BUILD_SHARED_LIBS=OFF' setting
    set(BUILD_SHARED_LIBS OFF CACHE INTERNAL "Build Static Library")
    add_subdirectory(${wxwidgets_SOURCE_DIR} ${wxwidgets_BINARY_DIR})
    # Restore the old value of 'BUILD_SHARED_LIBS'
    set(BUILD_SHARED_LIBS ${BUILD_SHARED_LIBS_OLD})
endif()

add_executable(test main.cpp)