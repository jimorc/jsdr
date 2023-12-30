# To determine which wxWidgets libraries are required, see
# https://wxwidgets.org/3.2.4/wx_librarieslist.html
# Change 3.2.4 to match the current version.


set(CMAKE_CXX_FLAGS -Wall)

set(SOURCES jsdr_app.cpp jsdr_mainframe.cpp jsdr_config.cpp)
set(INCLUDES ../include/jsdr_app.h ../include/jsdr_mainframe.h ../include/jsdr_config.h)

add_executable(
    jsdr
    ${SOURCES}
    ${INCLUDES}
)

include(Format)
Format(jsdr .)
Format(jsdr ./../include)

include(ClangTidy)
AddClangTidy(jsdr)

target_include_directories(
    jsdr
    PRIVATE ../include
    PRIVATE ${wxWidgets_INCLUDE_DIRS}
    PRIVATE ${jsoncpp_INCLUDE_DIR}
)

target_link_directories(
    jsdr
    PRIVATE ${jsoncpp_LIBRARIES}
)

target_link_libraries(
    jsdr
    ${wxWidgets_LIBRARIES}
    jsoncpp
)