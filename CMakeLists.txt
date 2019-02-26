cmake_minimum_required (VERSION 2.8.11)
project(pool)

add_subdirectory(hashing)

add_custom_command(
    OUTPUT build/bin
    COMMAND build/env.sh go get -v ./...
)

add_custom_target(build ALL DEPENDS hashing ${CMAKE_CURRENT_BINARY_DIR}/build/bin)
