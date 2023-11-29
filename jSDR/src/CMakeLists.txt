# To determine which wxWidgets libraries are required, see
# https://wxwidgets.org/3.2.4/wx_librarieslist.html
# Change 3.2.4 to match the current version.

find_package(wxWidgets COMPONENTS base core REQUIRED)
include(${wxWidgets_USE_FILE})

include(FetchContent)
FetchContent_Declare(
    jsoncpp
    GIT_REPOSITORY  https://github.com/open-source-parsers/jsoncpp.git
    GIT_TAG         origin/master
)

FetchContent_MakeAvailable(jsoncpp)
set(jsoncpp_INCLUDE_DIR ${CMAKE_BINARY_DIR}/_deps/jsoncpp-src/include)
set(jsoncpp_LIBRARIES ${CMAKE_BINARY_DIR}/_deps/jsoncpp-build/src/lib_json)

set(CMAKE_CXX_FLAGS -Wall)


add_executable(
    jSDR
    jSDRApp.cpp
    jSDRMainFrame.cpp
)

target_include_directories(
    jSDR
    PRIVATE ../include
    PRIVATE ${wxWidgets_INCLUDE_DIRS}
    PRIVATE ${jsoncpp_INCLUDE_DIR}
)

target_link_directories(
    jSDR
    PRIVATE ${jsoncpp_LIBRARIES}
)

target_link_libraries(
    jSDR
    ${wxWidgets_LIBRARIES}
    jsoncpp
)