include(FetchContent)

FetchContent_Declare(
    googletest
    GIT_REPOSITORY  https://github.com/google/googletest.git
    GIT_TAG         main
)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
FetchContent_MakeAvailable(googletest)

add_executable(tests
    configtests.cpp
    ../src/jSDRConfig.cpp
)

include(Format)
Format(tests .)


target_include_directories(tests 
    PRIVATE ../include
    PRIVATE ${wxWidgets_INCLUDE_DIRS}
    PRIVATE ${jsoncpp_INCLUDE_DIR}
)

target_link_directories(
    tests
    PRIVATE ${jsoncpp_LIBRARIES}
)

target_link_libraries(tests 
    PRIVATE gtest_main jsoncpp ${wxWidgets_LIBRARIES}
)

include(GoogleTest)
gtest_discover_tests(tests)