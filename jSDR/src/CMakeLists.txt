# To determine which wxWidgets libraries are required, see
# https://wxwidgets.org/3.2.4/wx_librarieslist.html
# Change 3.2.4 to match the current version.

find_package(wxWidgets COMPONENTS base core REQUIRED)
include(${wxWidgets_USE_FILE})

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
)

target_link_libraries(
    jSDR
    ${wxWidgets_LIBRARIES}
)