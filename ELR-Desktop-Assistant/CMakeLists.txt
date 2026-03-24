cmake_minimum_required(VERSION 3.16)

project(ELRDesktopAssistant VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Qt6 REQUIRED COMPONENTS Widgets Network)

set(PROJECT_SOURCES
    main.cpp
    ELRDesktopAssistant.cpp
    ELRDesktopAssistant.h
    TrayIcon.cpp
    TrayIcon.h
    DesktopWidget.cpp
    DesktopWidget.h
    ELRClient.cpp
    ELRClient.h
)

qt_add_executable(ELRDesktopAssistant
    ${PROJECT_SOURCES}
)

target_link_libraries(ELRDesktopAssistant PRIVATE Qt6::Widgets Qt6::Network)

set_target_properties(ELRDesktopAssistant PROPERTIES
    WIN32_EXECUTABLE TRUE
)
