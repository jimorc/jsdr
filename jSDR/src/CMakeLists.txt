# To determine which wxWidgets libraries are required, see
# https://wxwidgets.org/3.2.4/wx_librarieslist.html
# Change 3.2.4 to match the current version.


set(CMAKE_CXX_FLAGS -Wall)


add_executable(
    jSDR
    jSDRApp.cpp
    jSDRMainFrame.cpp
    jSDRConfig.cpp
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