# Coding Standards

This document outlines the coding standards used in this project. I have developed these standards
from over 30 years of coding
in C++; they work for me and therefore are required for this project.

> *Note:* This is the first time that I have documented my preferred coding standards. Therefore,
this document is dynamic and is in the early stages of being written. Hence, you should expect that
changes will be
made to it.
> *Note:* clang-format is used to format the source code in this project. Therefore, this document
only covers the coding standards
that are not handled directly by clang-format.

## Line Length

The length of all lines (source code, documentation such as markdown files, etc.) should not exceed
132 characters. This length
is selected as a compromise because 80 columns leads to the wrapping of too many source lines, and
exceeding 100 characters may
cause problems on smaller displays.

> *Note:* The `.vscode/settings.json` file in this project has a setting that displays a vertical
line after 100 characters. Other
IDEs are not currently supported, but pull requests for setting the line length indicator on other
IDEs will be glady accepted.

:

> *Note:* clang-format modifies lines inC++ header and source files to not exceed 100 characters.
There is no check in other file types.

## File Naming Conventions

### Project, Application and Library Names

Project, application and library names shall match the following criteria:

- Lower case ASCII letters only, no numbers.
- Words in the names may be separated by single underscore characters (_). The use of dashes (-) is
discouraged.

### Directory and File Names

Directory and file names shall match the following criteria:

- Lower case ASCII only.
- Words in the names may be separated by single underscore (_) characters. The use of dashes (-) is
discouraged.
_ There shall be no whitespace in the names.

### Markdown File Mames

Markdown files that are specific to the Project shall:

- Be named in upper case ASCII only.
- Have words in the file names separated by an underscore (_) character.
- Have a file extension of ".md".
- Have a name that describes the file's content. For example, the Readme file is called `README.md`,
 and the file that outlines the
steps to build the project is called `PROJECT_BUILD_DIRECTIONS.md`.

Markdown files that document code interfaces and how to use the program built by this project shall:

- Be named in lower case ASCII only.
- Have words in the file names separated by an underscore (_) character.
- have a file extension of ".md".
- Have a name that describes the files's content. For example, the file that contains instructions
on starting the jsdr application
shall be called `starting_jsdr.md`.

### Text File Names

Text files shall:

- Be named in lower case ASCII only unless there is another standard that should be followed. For
example, a configuration file
that contains text might be called `config.txt`, but CMake configuration files shall be called
`CMakeLists.txt` because that is
the standard for CMake.
- Have words in the file names separated by an underscore (_) character.
- Have a name that describes the file's contents. For example, a text file that contains settings
for a program might be called `settings.txt` or `jsdr_settings.txt`.
- have a file extension of `.txt` unless there is a standard that specifies otherwise.

#### Exceptions

Projects that are generated on github.com contain a text file called `LICENSE` in the project's
top-level directory. It is not
necessary to rename this file to `license.txt`.

### Source and Header File Names

There are many different naming conventions for C++ header and source files. The following
convention shall be used for this
project:

- For C++ code files, the file names should indicate the names of the classes. For example, the
header file containing the
declaration for the class MyBaseClass should be called mybaseclass.h and the source file for the
class MyBaseClass should be called
mybaseclass.cpp.
- The file extension .h shall be used for C++ header files.
- The file extension .cpp shall be used for C++ source files.

There are more requirements for header and source files listed below.

## C++ Coding Standards

### Namespaces

TO BE WRITTEN

### Classes and Structs

#### Class and Struct Names

Classes and Structs shall have names that:

- Describe the purpose or content of the class or struct.
- Use PascalCase. That is, the first character of each word in the names shall be capitalized.
- Use ASCII letters (a-z and A-Z) in their names. Underscore and digit characters shall not used.
- Do not contain "Base" or "Derived". In almost every case, more descriptive names may be created.

#### Class and Struct Method Names

Methods shall have names that:

- describe what the method does.
- Contain ASCII letters (A-Z and a-z) only. Method names should not include underscore characters.
- Contain digits only in rare cases. For example, a method might be called `times2`, although a
better choice would be `timesTwo`.
- Use camelCase. That is, the first character of the method name shall be lower case, and each
subsequent word in the name shall
begin with an upper case letter. For example, a method might be called `calculateVolume`.
- Do not distinguish between public, protected, and private names. In other words, a private method
 that calculates the area of
a shape shall be called `calculateArea`, not `privateCalculateArea`.

#### Class Variables

Class variables are normally private members of a class. To document that, class variables shall
have names that:

- Contain ASCII letters (A-Z, a-z).
- Contain digits only in rare cases.
- Use camelCase. That is, the first character of the variable name, ignoring the next point, shall
be lower case and the start
of each subsequent word in the variable name shall begin with an upper case letter.
- Prepend each variable name with `m_`. For example, `m_address`, and `m_configurationFileName`.

#### Struct Variables

Struct variables are normally public members of a struct. To document that, struct variables shall
have names that:

- Contain ASCII letters (A-Z, a-z).
- Contain digits only in rare cases.
- Use camelCase. This is, the first character of the variable name shall be lower case and the
start of each subsequent word
in the variable name shall begin with an upper case letter.
- Not be prepended or postpended with any indicator to show that the variable has public access. =
There are sufficient other ways
built into the C++ language to illustrate that.

Here is a simple struct example:

```C++
struct Circle {
    double radius;
    Color displayColor;
};
```

### Header Files

Include guards shall not be used; instead, every header file shall contain the follwing at or near
the first line in the file:

```C++
#pragma once
```

While this directive is non-standard, it is supported by all compilers that can reasonably be
expected to be used to build this project (MSVC, gcc, clang, Apple Clang). It has the following
advantages:

- Less code.
- Avoids potential name clashes.
- Possible improvement in compilation speed.

After the `#pragma` directive, each header file shall contain a short copyright and license
statement
