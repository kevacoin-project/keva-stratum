cmake_minimum_required (VERSION 2.8.11)

project(keva-stratum)

add_subdirectory(cnutil)

add_custom_target(build ALL DEPENDS cnutil COMMAND go build)

