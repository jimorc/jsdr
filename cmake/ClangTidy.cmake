find_program(CLANG_TIDY_PATH clang-tidy REQUIRED)
set(CLANG_TIDY_CONFIG "${CMAKE_SOURCE_DIR}/.clang-tidy")
set(CMAKE_CXX_CLANG_TIDY ${CLANG_TIDY_PATH} 
    --config-file=${CLANG_TIDY_CONFIG};--header-filter=../include/jsdr*;--warnings-as-errors=*;)