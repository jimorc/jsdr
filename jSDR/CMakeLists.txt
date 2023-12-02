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

add_subdirectory(src)
add_subdirectory(tests)