cmake_minimum_required (VERSION 2.8.11)

project(keva-stratum)

add_subdirectory(cnutil)

add_custom_target(build ALL DEPENDS cnutil COMMAND go get -d ../. COMMAND go build ..)

if(MSYS)
    add_custom_target(copy ALL COMMAND cp cnutil/libcnutil.dll .)
endif(MSYS)


