set(LIB "hashing")

set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} -std=c11 -D_GNU_SOURCE")

add_library(${LIB} src/hashing.c)

