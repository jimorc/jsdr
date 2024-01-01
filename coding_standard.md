# Coding Standards

This document outlines the coding standards used in this project. I have developed these standards
from over 30 years of coding
in C++; they work for me and therefore are required for this project.

This document is based on the structure and topics in the 
[Google C++ style guide](https://google.github.io/styleguide/cppguide.html). That means that the
sections, headings and so forth are as outlined in that style guide; while some of the coding
standards reported in that document are copied to this standards document, there are many
differences.

> *Note:* The Google C++ style guide does not include any license. Because of that, I have assumed
that any use of the contents of the document is fair use. For example, in some cases, I have
copied the text of the document verbatim without noting each occasion where I did so. For every
standard in this document that has a corresponding section in the Google style guide, you should
assume that at least some text has been copied directly.

:

> *Note:* This is the first time that I have documented my preferred coding standards. Therefore,
this document is dynamic and is in the early stages of being written. Hence, you should expect that
changes will be made to it.

:

> *Note:* clang-format and clang-tidy are both applied to the source code in this project.
clang-format will automatically reformat the source code, but clang-tidy only reports errors.
Therefore, the code will not compile if you do not follow the required standards.

## C++

### Compilers

The code in this project shall be written to be compiled by at least the following compilers:
MSVC, gcc, clang, and Apple Clang.

#### Why These Compilers

The outputs of this project, the applications and libraries, are intended for use on 
the latest versions of at least the
following operating systems: Windows, MacOS, and Linux. The code must be compilable by each of the
listed compilers to satisfy that requirement. Additionally, if the project can be built using those
compilers, it can also probably be ported to lesser used operating systems as well.

#### How These Compilers Are Enforced

There is no direct way of enforcing this standard other than ensuring that the project compiles
and runs on each of the required operating systems using the most common compilers used on those
operating systems.

### C++ Version

This project targets C++23.

Do not use non-standard extensions.

#### Why This Version

The latest approved standard as of December 2023 is C++23. At that date, each of the
commonly used C++
compilers (MSVC, gcc, clang, Apple clang) support most or all C++23 features, but they each
support a different subset of C++26 features.
It would be difficult to port this project from one compiler to another if C++26 features are used.
Non-standard exensions are just that: non-standard. They are supported by one or a couple of
the common compilers, but not all. Using non-standard extensions makes porting code difficult if
not impossible.

The most important feature of C++23 that is specifically mentioned in this document is
`std::format` for which there is no prior implementation. 

#### How This Version Is Enforced

The required C++ standard is declared in the project's top-level CMakeLists.txt file. If there
are libraries that are included in this project that must be built from source code and that
require an earlier C++ version, then the CMakeLists.txt file must change the C++ version prior
to building the library, and then change the C++ version back to this required standard immediately
after the external library is built.

## Header Files

In general, every .cpp source file should have an associated .h header file. There are some common
exceptions,
such as unit tests and possibly small .cpp files that just contain a main() function.

### Why There Should Be A Matching Header File For Each Source File

Correct use of header files can make a huge difference to the readability, size and performance of
the code.

### How This Header Files Standard is Enforced

I know of no automatic way to enforce this standard. Therefore, the only way appears to be via
code review.

### Self-contained Headers

Header files should be self-contained (compile on their own) and end in .h. Non-header files
should end in .inc and be used sparingly.

When a header declares inline functions or templates that clients of the header will instantiate,
the inline functions and templates must also have definitions in the header, either directly or
in files it includes. Do not move these definitions to separately included header (-inl.h) files.
When all instantiations of a template occur in one .cpp file, either because they are explicit or
because the definition is accessible only in the .cpp file, the template definition should be
kept in that file.

There are rare occasions where a file designed to be included is not self-contained. These are
typically intended to be included at unusual locations, such as the middle of another file. They
might not use include guards, and might not include their prerequisites. Name such files with the
.inc extension. Use sparingly, and prefer self-contained header files when possible.

#### Why Header Files Should Be Self-Contained

Users and refactoring tools should not have to adhere to special conditions to include a header.
Specifically, a header should have include guards and include all other headers that it needs.

#### How Self-Contained Header Files Are Enforced

The compiler will find some cases where header files are not self-contained.

clang-tidy misc-header-include-cycle

### Include Guards

To prevent header files from being included more than once in a compile unit, they should contain
include guards or `#pragma once` directives. Include guards are preferred over `#pragma once`.

To guarantee uniqueness, the macro should be based on the full path in the project's source tree.
For example, the file `foo/src/bar/baz.h` should contain the following include guard:

```C++
#ifndef FOO_BAR_BAZ_H
#define FOO_BAR_BAZ_H

...

#endif      // FOO_BAR_BAZ_H
```

The comment following the `#endif` directive shall contain the macro name.

#### Why Include Guards Should Be Used

For the reasons outlined in the post
[Include Guard vs #pragma once in C and C++](https://computingonplains.wordpress.com/2023/12/05/include-guard-vs-pragma-once-in-c-and-c/),
include guards are preferred over `#pragma once`.

#### How Include Guards Are Enforced

clang-tidy has the rule `llvm-header-guard` which checks that the macro in the `#define` matches the
macro in the `#ifndef`.

The existence and equality of the comment following the `#endif` directive
is not checked. There are no checks for the existence of an include guard, and there are
no checks for the existence of `#pragma once`. These can only be checked via code review.

### Include What You Use

If a source or header file refers to a symbol defined elsewhere, the file should directly include
a header file which provides a declaration or definition of that symbol. It should not include
header files for any other reason.

Do not rely on transitive inclusions. The file `foo.cpp` should include `bar.h` even if `foo.h`
includes `bar.h`.

#### Why You Should Include What You Use

Transitive inclusions allow people to remove no-longer-needed `#include` directives from their
headers without breaking clients.

#### How Include What You Use is Enforced

The clang-tidy rule `misc-include-cleaner` checks for both the inclusion of header files that are
not referenced direcly and the non-inclusion of header files that are referenced directly.

## Avoid Forward Declarations

Avoid forward declarations where possible. Instead, 
[include the header files you need](#include-what-you-use).

A forward declaration is a declaration of an entity without an associated definition.

You may find that a forward declaration is needed to prevent circular `#include`s. If that is the
case, you should consider a design change instead.

### Why to Avoid Forward Declarations

* Forward declarations can hide a dependency, allowing user code to skip necessary recompilation
when headers change.
* A forward declaration as opposed to an `#include` makes it difficult for automatic tooling to
discover the module defining the symbol.
* A forward declaration may be broken by subsequent changes to the library. Forward declarations
of functions and templates can prevent the header owners from making otherwise-compatible changes
to their APIs, such as widening a parameter type, adding a template parameter with a default value,
or migrating to a new namespace.
* Forward declaring symbols from namespace `std::` yields undefined behavior.
* It can be difficult to determine whether a forward declaration or a full `#include` is needed.
Replacing an `#include` with a forward declaration can silently change the meaning of code.
* Forward declaring multiple symbols from a header can be more verbose than simply `#inlude`ing
the header.
* Structuring the code to enable forward declarations (e.g. using pointer members instead of
object members) can make the code slower and more complex.

#### How Avoid Forward Declarations Is Enforced

There is no check available in clang-tidy to flag forward declarations, so the only way to check
for the presence of forward declarations is via code review.

Note, however, that clang-tidy has the following rule related to forward references:

* `bugprone-forward-declaration-namespace`
  
### Inline Functions

Define functions inline only when they are small, say, 10 lines or fewer.

Function types that should be considered for inlining include accessors, mutators, and other short, 
performance critical functions.

Functions are not always inlined even if they are declared as such. For example, virtual and
recursive functions are not normally inlined. The main reason for making a virtual function
inline is to place its definition in the class, either for convenience or to document its
behavior, e.g. for accessors and mutators.

#### Why You Should Inline Functions

Inlining a function can generate more efficient object code, as long as the inlined function is
small.

#### Why You Should Not Inline Functions

Overuse of inlining can actually make programs slower. Depending on a function's size, inlining
can cause the code size to increase or decrease. Inlining a very small accessor function will
usually decrease code size while inlining a very large function can dramatically increase code
size. On modern processors, smaller code usually runs faster due to better use of the instruction
cache.

Destructors should not be inlined because they are often longer than they appear because of
implicit member and base-destructor calls.

It is usually not cost-effective to inline functions with loops or switch statements unless in the
common case, the loop or switch statement is never executed.

#### How Inlining Is Enforced

As specifying functions for inlining is optional, clang-tidy does not provide any rules for checking
for the use or abuse of inlining. Functions that are marked for inlining may not be inlined by
modern compilers. Code review is the only way to check that inlining is not being abused.

clang-tidy does have one rule related to inlining: `llvmlibc-inline-function-decl`. That rule is
meant only for use on the `libc` library.

### Include The C++ Header Rather Than The Equivalent C Header

The C++ standard provides C++ header files for C the standard library. Therefore, the C++ header 
files should be included rather than the corresponding
C library headers. For example, include `#include <cstdlib>` instead of `#include <stdlib.h`.

#### Why Include The C++ Header Rather Than The Equivalent C Header

The C library header files will be deleted in a future C++ standard. Using the equivalent
C++ header files will ensure that the header file will be available in future compiler releases.

#### How Include the C++ Header Rather Than The Equivalent C Header Is Enforced

Deprecated C library header files are detected by the clang-tidy rule
`modernize-deprecated-headers`.

### Names and Order of Includes

Header files shall be included and blocked in the following order:

1. Related header
2. A blank line
3. C system headers (headers in angle brackets with the .h extension); e.g. <unistd.h>, <stdlib.h>
4. A blank line
5. C## standard library headers (without file extension); e.g. <algorithn>, <cstddef>
6. A blank line
7. Other libraries' .h files
8. A blank line
9. This project's .h files

Separate each non-empty group with one blank line.

#### Why Use This Order Of Include

With this ordering, if the related header omits any necessary includes, the build of the source
(.cpp) file will break. Thus, this rule ensures that the build breaks show up first for the
people working on these files, not for innocent people in other packages.

#### How Names and Order of Includes Is Enforced

The setting of the clang-format style option `IncludeBlocks`, in conjunction with the style option
`IncludeCategories` determines the order in which `#include` directives are sorted.

The `IncludeBlocks` style option has been set to `Regroup`. No `IncludeCategories` are currently
defined, but should be to force the include order.

## Scoping

### Namespaces

#### Use Namespaces

With few exceptions, code should be placed in a namespace. Namespaces should have unique names
based on the project name and possibly its path.

Code such as macros that reference unscoped libraries cannot be placed inside a namespace. For
example, the wxWidgets library is not inside a namespace because wxWidgets predates
C++ namespaces. Attempting the following results in a compiler error:

```C++
namespace foo {
    wxIMPLEMENT_APP(foo::fooApp);
    ...
}   // namespace foo
```

Instead, you must code this as:

```C++
wxIMPLEMENT_APP(foo::fooApp);

namespace foo {
    ...
}   // namespace foo
```

##### Why Use Namespaces

Namespaces divide the global scope into distinct, named scopes, and so is useful for preventing
name collisions in the global scope.

For example, if tow different projects have a class `Foo` in the global scope, these symbols may
collide at compile time or at runtime. If each project places their code in a namepsace,
`project1::Foo` and `project2::Foo` are now distinct symbols that do not collide, and code within
each project's namespace can continue to refer to `Foo` without the prefix.

##### How Use Namespaces Is Enforced

There are no clang-tidy rules to enforce the use of namespaces. The only method of enforcement
is code review.

#### Base Namespace Names On The Project Name

Namespaces should have unique names based on the project name, and possibly its path.

##### Why Base Namespace Names On The Project Name

Basing the top-level namespace name on the project name is more likely to prevent namespace name
collisions.

##### How Base Namespace Names On The Project Name Is Enforced

This standard can only be enforced by code review.

#### Inline Namespaces

Inline namespaces automatically place their names in the enclosing scope. For example:

```C++
namespace outer {
inline namespace inner {
    void foo();
}   // namespace inner
}   // namespace outer
```

The expressions `outer::inner::foo` and `outer::foo` are interchangeable.

##### Why To Use Inline Namespaces

Inline namespaces are primarily intended for ABI compatibility across versions.

##### Why Not To Use Inline Namespaces

Inline namespaces can be confusing because names aren't actually restricted to the namespace where
they are declared. They are only useful as part of some larger versioning policy.

##### How Inline Namespaces Is Enforced

The only way to enforce the use or non-use of inline namespaces is by code review.

#### Terminate Multiline Namespaces With Namespace Name Comment

Terminate each namespace with a comment specifying the namespace name. The example in
[Inline Namespaces](#inline-namespaces) shows this:

```C++
}   // namespace inner
}   // namespace outer
```

##### Why Terminate Multiline Namespaces With Namespace Name Comment

Doing so helps to delineate the end of the namespace declarations or definitions. Not doing so
may result in extending the time required to determine the causes of compiler errors.

##### How Terminate Multiline Namespaces With Namespace Name Comment Is Enforced

Code review is the currently the only way this standard can be enforced.

#### Do Not Use Using-Directive Indiscrimiately

The following using directive imports all `std::` names into the current namespace:

```C++
using namespace std;
```

This is almost certainly what you do not want to do. If you will be referencing a name in a
namespace seldom then include the namespace when referencing the name. For example:

```C++
std::cout << "Print something\n";
```

Alternatively, for times when you might use a namespace name frequently, then use, for example:

```C++
using std::cout;
cout << "Print first line.\n"
cout << "Print next line.\n"
cout << "Print last line.\n"
```

##### Why Not To Use Using-Declaration Indiscriminately

Each time you use the `using` declaration, you pollute the enclosing namespace. If you use the
`using` declaration to import all of a namespace, you are typically adding many new names to the
enclosing namespace. With each added the likelihood of a name collision increases.

Never, ever use `using std;`. The number of names in `std` is huge.

##### How Not To Use Using-Declaration Is Enforced

The clang-tidy rule `google-build-using-namespace` flags all using-declarations that make all
names from a namespace available.

#### Do Not Use Namespace Aliases At Namespace Scope In Header Files

Do not place a declaration such as:

```C++
namespace baz = ::foo::bar::baz;
```

inside a header file except in explicitly marked internal-only namespaces because anything imported
into a namespace in a header file becomes part of the public API exported by that file. Note that
including that statement in a .cpp file is acceptable.

Instead, use:

```C++
namespace foo {
namespace impl {    // Internal, not part of API
namespace bar = ::diagnostics::bar;
}   // namespace impl

void someFunction() {
    // namespace alias local to a function or method
    namespace baz = ::foo::bar::baz;
    ...
}
}   // namespace foo
```

##### Why Not To Use Namespace Aliases At Namespace Scope In Header Files

As mentioned above, including a using declaration in a header file, you pollute the enclosing
namespace. Doing so makes all names in the added namespace part of the API specified by the
header file.

##### How Do Not Use Namespace Aliases At Namespace Scope In Header Files

Initially, this can only be caught by code review. Unfortunately, the problem usually becomes 
obvious only much later, typically showing up as strange and mysterious compile or runtime errors.

#### Internal Linkage

When definitions in a .cpp file do not need to be referenced outside that file, give them
internal linkage by placing them in an unnamed namespace or by declaring them `static`. Do not
use either of these constructs in .h files.

All declarations can be given internal linkage by placing them in unnamed namespaces. Functions
and variables can also be given internal linkage by declaring them `static`.

Format unnamed namespaces like named namespaces. In the
[terminating comment](#terminate-multiline-namespaces-with-namespace-name-comment),
leave the namespace name empty:

```C++
namespace {
    ...
}   // namespace
```

##### Why Use Internal Linkage

Anything you declare with internal linkage cannot be accessed from another file. If a differ###ent
file declares something with the same name, the two entities are completely independent.

##### How Internal Linkage is Enforced

Because compilers cannot determine if you want anything outside of named namespaces to be internally
or externally linked, the only way to enforce internal linkage is via code review.

#### Nonmember, Static Member, and Global Functions

Prefer placing nonmember functions in a namespace; use completely global functions rarely.

Do not use a class simply to group static members.

Nonmember and static member functions may make more sense as members of a new class, especially
if they access external resources or have significant dependencies.

Sometimes it is useful to define a function that is not bound by a class instance. Such a function
can be either a static member or a nonmember function. Nonmember functions should not depend on
external variables, and should nearly always exist in a namespace. Do not create classes only to
group static members; this is no different than just giving the names a common prefix, and such
grouping is usually unnecessary anyway.

If you define a nonmember function and it is only needed in its .cpp file, use
[internal linkage](#internal-linkage) to limit its scope.

##### Why Place Nonmember Functions in a Namespace

Putting nonmember functions in a namespace avoids polluting the global namespace.

##### How Nonmember and Static Member Placement is Enforced.

Compilers cannot determine why nonmember and static member functions are placed where they are.
Therefore, the only way to enforce placement is through code review.

#### Local Variables

Place a function's variables in the narrowest scope possible, and initialize variables in the
declaration.

C++ allows you to declare variables anywhere in a function. You are encouraged to declare them in
a scope as local as possible, and as close to the first use as possible.

Initialization should be used instead of declaration and assignment.

Variables needed for `if`, `while`, and `for` statements should normally be declared within those
statements so that these variables are confined to those scopes. There is one caveat: if the
variable is an object, its constructor is invoked every time it enters scope and is created, and
its destructor is invoked every time it goes out of scope. In that case, it may be more efficient
to declare such a variable outside the loop.

##### Why Declare Local Variables In As Local Scope and As Close To First Use As Possible

This makes it easier for the reader to find the declaration, see what type the variable is, and
what it is initialized to.

##### How Local Variable Declaration Is Enforced.

The clang-tidy rule `cppcoreguidelines-init-variables` checks that local variables are declared
with an initial value.

For other local variable issues with the placement of local variable declartion, code review 
appears to be the only choice.

#### Static and Global Variables

Objects with [[[static storage duration]]] are forbidden unless they are
[[[trivially destructible]]]. Informally, this means that the destructor does not do anything,
even taking member and base destructors into account. More formally, it means that the type
has no user-defined or virtual destructor and that all bases and non-static members are trivially
destructible. Static function-local variables may use dynamic allocation. Use of dynamic
initialization for static class member variables or variables at namespace scope is discouraged, but
allowed in limited circumstances.

As a rule of thumb, a global variable satisfies these requirements if its declaration, considered
in isolation, could be `constexpr`.


##### When To Use Static and Global Variables

Static and global variables are useful for a large number of applications: named constants,
auxilliary data structures internal to some translation unit, command-line flags, logging,
registration mechanism, background infrastructure, and so forth.

##### When Not To Use Static and Global Variables

Static and global variables that use dynamic initialization or have non-trivial destructors create
complexity that can easily lead to hard-to-find bugs. Dynamic initialization is not ordered across
translation units, and neither is destruction, except that destruction happens in reverse order
of initialization. When one initialization refers to another variable with static storage duration,
it is possible that this causes an object to be accessed before its lifetime has begun or after
its lifetime has ended. Moreover, when a program starts threads that are not joined at exit, those
threads may attempt to access objects after their lifetime has ended if their destructor has
already run.

##### How To Enforce Static and Global Variable Initialization and Destruction

The only way to ensure that static and global variables do not use dynamic initialization or have
non-trival destruction is through code review.

##### Static and Global Variable Examples

###### Global Strings

If you require a global or static string, consider using a `constexpr` variable of `string_view`,
character array, or character pointer that points to a string literal. String literals have
static storage duration already and are usually sufficient.

###### Maps, Sets and Other Dynamic Containers

If you require a static, fixed collection, such as a set to search against or a lookup table, you
cannot use the dynamic containers from the standard library as a static variable since they have
non-trivial destructors. Instead, consider a simple array of trivial types, e.g. an array of arrays
of ints (for a map from int to int), or an array of pairs, e.g. pairs of int and const char*. For
small collections, linear search is entirely sufficient, and efficient due to memory locality. If
necessary, keep the collection in sorted order and use a binary search algorithm. If you do
really prefer a dynamic container from the standard library, consider using a function-local static
pointer, as described in the next section.

###### Smart Pointers

Smart pointers (std::unique_ptr, std::shared_ptr) execute cleanup during destruction and are
therefore forbidden. Consider whether you use case fits into one of the other patterns described
in this set of examples. One simple solution is to use a plain pointer to a dynamically allocated
object and never delete it. See
[Create An Object Dynamically And Never Delete It](#create-an-object-dynamically-and-never-delete-it).

###### Static Variables of Custom Types

If you require static, constant data of a type that you need to define yourself, give the type a
trivaial destructor and a `constexpr` constructor.

###### Create An Object Dynamically And Never Delete It

If all else fails, you can create an object dynamically and never delete it by using a
function-local static pointer or reference (e.g. `static const auto& impl = *new T(args...);`)

#### thread_local Variables

`thread_local` variables that aren't declared inside a function must be initialized with a
compile-time constant, i.e. they must have no dynamic initialization.

Prefer `thread_local` over other ways of defining thread-local data.

Variables can be declared with the `thread_local` specifier:

```C++
thread_local Foo foo = ...;
```

Such a variable is actually a collection of variables, so that when different threads access it,
they are actually accessing different objects. `thread_local` variables are much like
[static storage duration variables](#static-and-global-variables) in many respects. Fir instance,
they can be declared at namespace scope, inside functions, or as static class members, but not
as ordinary class members.

`thread_local` variable instances are initialized much like static variables, except that they
must be initialized separately for each thread, rather than once at program startup. This means
that `thread_local` variables declared within a function are safe, but other `thread_local`
variables are subject to the same initialization-order issues as static variables, and more
besides.

`thread_local` variables have a subtle destruction-order issue: during thread shutdown,
`thread_local` variables will be destroyed in the opposite order of their initialization, as is
generally true in C++. If code triggered by the destructor of any `thread_local` variable refers
to any already-destroyed `thread_local` variable on that thread, we will get a particularly hard
to diagnose use-after-free.

##### Why Use `thread_local` Variables

Thread-local data is inherently safe from races because only one thread can ordinarily access it,
which makes `thread_local` useful for concurrent programming. `thread_local` is the only
standard supported way of creating thread-local data.

##### Why Not Use `thread_local` Variables

Accessing a `thread_local` variable may trigger execution of an unpredicable and uncontrollable
amount of other code during thread-start or first use on a given thread.

`thread_local` variables are effectively global variables, and have all of the drawbacks other
than lack of thread-safety.

The memory consumed by a `thread_local` variable scales with the number of running threads in the
worst case, which can be quite large for a program.

Data members cannot be `thread_local` unless they are also `static`.

We may suffer use-after-free bugs if `thread_local` variables have complex destructors. In
particular, the destructor of any such variable must not call any code transitively, that refers
to any potentially-destroyed `thread_local`. This property is hard to enforce.

Approaches for avoiding use-after-free in global/static contexts do not work for `thread_local`s.
Specifically, skipped destructors for globals and static variables is allowable because their
lifetimes end at program shutdown. This any "leak" is managed immediately by the operating system
cleaning up our memory and other resources. By contrast, skipping destructors for `thread_local`
variables leads to resource leaks proportional to the total number of threads that terrminate
during the lifetime of the program.

##### How To Enforce `thread_local` Initialization

As stated above, it is very difficult to spot when `thread_local` variables are initialized with
dynamic data. Code review _may_ spot some problems.

### Classes

Classes are the fundamental unit of code in C++, so we use them extensively. This section lists
the main do's and don'ts you should follow when writing a class.

#### Doing Work In Constructors

It is possible to perform arbitrary initialization in the body of a constructor.
Performing work in constructors has a number of advantages and disadvantages which are covered in
subsections below. There are some tasks that cannot be preformed in a constructor, so an alternative
way of performing those tasks, such as an Init() method, must be used.

* Do not call virtual functions from within a constructor.
* Throwing an exception, or terminating the program are the only ways to handle errors in a
constructor.
* Avoid Init() methods on objects with no other states that affect which public methods may be
called. Semi-constructed objects of this form are particularly hard to work with correctly.

##### Why To Do Work In Constructors

There is no need to worry about whether the class has been initialized or not. If the object exists,
it has been initialized.

Objects that have been fully initialized by a constructor call can be `const` and may also be
easier to use with standard containers or algorithms.

##### Why Not To Do Work In Constructors

If the constructor attempts to call a virtual function, the call will not be dispatched to the
subclass implementation. That is because the subclass does not yet exist when the parent class's
constructor is executed. A future modification to the class can quietly introduce this
problem even if the class is not currently subclassed, causing much confusion.

There is no easy way to signal errors, short of crashing the program or using exceptions.

If the work fails, we now have an object whose initialization code has failed, so it may be in an
unusual state requiring a `bool IsValid()` state checking mechanism which is easy to forget to call.

You cannot take the address of a constructor, so whatever work is done in the constructor cannot
easily be handed off to, for example, another thread.

##### How To Enforce Work In Constructors

Code review is the only semi-reliable method.

#### Implicit Conversions

Whenever possible, use `explicit` conversion operators and single-argument constructors. This, of
course, does not apply to copy and move constructors.

Implicit conversions allow an object of one type to be used where a different type is expected,
such as passing an `int` argument to a function that takes a `double` parameter.

In addition to the implicit conversions defined by the language, users can define their own, by
adding appropriate members to the class definition of the source or destination type. An implicit
conversion in the source type is defined by a type conversion operator named after the destination
type (e.g. `operator bool()`). An implicit conversion in the destination type is defined by a
constructor that can take the source type as its only argument, or only argument with no
default value.

The `explicit` keyword can be applied to a constructor or a conversion operator to ensure that it
can only be used when the destination type is explicit at the point of use, e.g., with a cast.
This applies not only to implicit conversions, but to list initialization syntax. See this example:

```C++
class Foo {
    explicit Foo(int x, double y);
    ...
};

void Func(Foo f);

Func({42, 3.14159});  // This is an error
```

This code isn't technically an implicit conversion, but the language treats it as one as far as
`explicit` is concerned.

##### Why To Use Implict Conversions

* Implicit conversions can make a type more usable and expressive by eliminating the need to
explicitly name a type when it is obvious.
* Implicit conversions can be a simpler alternative to overloading, such as when a single function
with a `string_view` parameter takes the place of separate overloads for `std::string` and
`const char*`.

##### Why To Not Use Implicit Conversions

* Implicit conversions can hide type-mismatch bugs, where the destination type does not match the
user's expectation, or the user is unaware that any conversion will take place.
* Implicit conversions can make code harder to read, particularly in the presence of overloading,
by making it less obvious what code is actually being called.
* Constructors that take a single argument may accidentally be usable as implicit type conversions,
even if they are not intended to do so.
* When a single-argument constructor is not marked `explicit` there's no reliable way to tell
whether it is intended to define an implicit conversion, or the author simply forgot to mark it.
* Implicit conversions can lead to call-site ambiguities, especially when there are bidirectional
implicit conversions. This can be caused either by having two types that both provide an implicit
conversion, or by a single type that has both an implicit constructor and an implicit type
conversion operator.
* List initialization can suffer from the same problems if the destination type is implicit,
particularly if the list has only a single argument.

##### Enforcing `explict` Type Conversions

clang-tidy has the rule `google-explicit-constructor` to ensure that constructors callable with a
single argument and conversion operators are marked explicit to avoid the risk or unintentional
implicit conversions.

#### Copyable and Movable Types

A class's API must make clear whether the class is copyable, move-only, or neither copyable nor
movable. Support copying and/or moving if these operations are clear and meaningful in the class.

A movable type is one that can be initialized and assigned from temporaries. A copyable type is
one that can be initialized or assigned from any other object of the same type, so it is also
movable by definition, with the stipulation that the value of the source does not change.
`std:unique_ptr` is an example of a movable but not copyable type. `int` and `std::string` are
examples of movable types that are also copyable. For `int`, the move and copy operations are the
same; for `std::string`, there exists a move operation that is less expensive that a copy.

For user-defined types, the copy behavior is defined by the copy constructor and the
copy-assignment operator. Move behavior is defined by the move constructor and the move-assignment
operator, if they exist, or by the copy constructor and the copy-assignment operator otherwise.

The copy/move constructors can be implicitly invoked by the compiler in some situations.

Every class's public interface must make clear which copy and move operations the class supports.
This usually takes the form of explicitly declaring and/or deleting the appropriate operations in
the `public` section of the declaration.

Specifically, a copyable class should explicitly declare the copy operations, a move-only class
should explicitly declare the move operations, and a non-copyable but moveable class should
explicitly delete the copy operations. A copyable class must also declare move operations in order
to support efficient moves. Explicitly declaring or deleting all four copy/move operations is
permitted but not required. If you provide a copy or move assignment operator, you must also
provide the corresponding constructor.

Declarations/deletions can be omitted only if they are obvious:

* If the class has no `private` secions, like a [`struct`](#structs-vs-classes) or an
interface-only base class, then the copyability/movability can be determined by the
copyability/movability of any public data members.
* If the base class clearly isn't copyable or movable, then derived classes won't be either. An
interface-only base class that leaves these operations implicit is not sufficient to make
concrete subclasses clear.
* If you explicitly declare or delete either the constructor or assignment operation for copy, the
other copy operation is not obvious and must be declared or deleted. Likewise for move operations.

##### Why Specify Copy/Move Types

* Objects of copyable and movable types can be passed and returned by value, which makes APIs
simpler, safer, and more general. Unlike when passing objects by pointer or reference, there's no
rish of confusion over ownership, lifetime, mutability, and similar issues, and no need to specify
them in the contract. It also prevents non-local interactions between the client and the
implementation, which makes them easier to understand, maintain, and optimize by the compiler.
Further, such objects can be used with generic APIs that require pass-by-value, such as most
containers, and they allow for additional flexibility, in e.g. type conversion.
* Copy/move constructors and assignment operators are usually easier to define correctly than
alternatives like `clone()`, `copyFrom()`, or `swap()`, because they can be generated by the
compiler, either implicitly or with `=default`. They are concise, and ensure that all data members
are copied. Copy and move constructors are also generally more efficient, because they don't
require heap allocation or separate initialization and assignment steps, and they are eligible for
optimizations such as [[[copy elision]]].
* Move operations allow the implicit and efficient transfer fo resources out of `rvalue` obejts.
This allows a plainer coding stype in some cases. 

##### Why Not Specify Copy/Move Types

* Some types do not need to be copyable, and providing copy operations for such types can be
confusing, nonsensical, or outright incorrect. Types representing singleton objects (`Registerer`),
objects tied to a specific scope (`Cleanup`), or closely coupled to object identity (`Mutex`)
cannot be copied meaningfully. Copy operations for base class types that are used polymorphically
are hazardous, because use of them can lead to [[[object slicing]]]. Defaulted or
carelessly-implemented copy operations can be incorrect, and the resulting bugs can be confusing
and difficult to diagnose.
* Copy constructors are invoked implicitly, which makes the invocation easy to miss. This may
cause confusion for programmers used to languages where pass-by-reference is conventional or
mandatory. It may also encourage excessive copying, which can cause performance problems.

##### How Copyable/Movable Types Are Enforced

There are currently no clang-tidy rules for this standard, presumably because different classes
may be copyable and/or movable, so code review is the only way to check that copy/move operations
are properly specified.

#### Structs vs. Classes

Use a `struct` only for passive objects that carry data; everything else is a `class`.

The `struct` and `class` keywords behave almost identically in C++. We add our own semantic
meanings to each keyword, so we should use the appropriate keyword for the data types we define.

`struct`s should be used for passive objects that carry data, and may have associated constants.
All fields must be public. The struct must not have invariants that imply relationships between
different fields, since direct user access to those fields would break those invariants.
Constructors, destructors, and helper methods may be present; however, these methods must not
require or enforce andy invariants.

If more functionality or if invariants are required, or the struct has wide visibility and is
expected to evolve, then a `class` is more appropriate. If in doubt, make it a `class`.

For consistency with STL, you can use `struct` instead of `class` for stateless types, such as
traits, [[[template metafunctions]]], and some functors.

Note that member variables in structs and classes have [[[differnt naming rules]]].

##### How Structs vs. Classes Is Enforced

Whether a user-defined object is a struct or class is up to the programmer. Code review helps
to catch the obvious wrong choices, such as private data or invariants in a struct.

#### Structs vs. Pairs And Tuples

Prefer to use a `struct` instead of a pair or a tuple whenevere the elements can have meaningful
names.

Pairs and tuples may be appropriate in generic code where there are not specific meanings for the
elements of the pair or tuple. Their use may also be required in ordeer to interoperate with any
existing code or APIs.

##### Why Use Structs vs. Pairs And Tuples

* While using pairs and tuples can avoid the need to define a custom type, potentially saving work
when _writing_ code, a meaningful field name will almost always be much clearer when _reading_
code than `.first` and `.second`, or `std::get<X>`. While C++14'a introduction of `std::get<Type>`
to access the tupe element by type rather than index when type is unique can sometimes partially
mitigate this, a field name is usually substantially clearer and more informative than a type.

##### How Structs vs. Pairs And Tuples Is Enforced

As with all standards that present a choice, the only way to enforce correct use is code review.

#### Inheritance

Compostition is often more appropriate than inheritance.

All inheritance should be public. If you want to use private inheritance, you should be including
an instance of the base class as a member instead. You can use `final` on classes when you don't
intend to support using them as base classes.

Do not overuse _implementation inheritance_. Composition is often more appropriate. Try to
restrict the use of inheritance to the "is-a" relationship.

Limit the use of `protected` to those member functions that might need to be accessed from
subclasses.

[[[Data members should be private]]].

Explicitly annotate overrides of virtual functions or virtual destructors with exactly one of
the `override`, or less frequently, the `final` specifier. Do not use `virtual` when declaring an
override. Rationale: A function or destructor marked `override` or `final` that is not an override
of a base class virtual function will not compile, and this helps to catch common errors. The
specifiers serve as documentation; if no specifier is present, the reader has to check all
ancestors of the class in question to determine if the function or destructor is virtual or not.

Limit the use of `protected` to those member functions that might need to be accessed from
subclasses. 

Multiple inheritace is permitted, but multiple _implementation_ inheritance is strongly
discouraged.

##### Why Use Inheritance

* _Implementation inheritance_ reduces code size by reusing the base class code as it specializes
an existing type. Because inheritance is a compile-time declaration, you and the compiler can
understand the operation and detect errors. _Interface inheritance_ can be used to programmatically
enforce that a class expose a particular API Again, the compiler can detect errors, in this case,
when a class does not define a necessary method of the API.

##### Why Not Use Inheritance

* For _implementation inheritance_, because the code implementing a subclass is spread between the
base and the subclass, it can be difficult to understand an implementation. The subclass cannot
override functions that are not virtual, so the subclass cannot change implementation.
* Multiple inheritance is especially problematic, because it often imposes a higher performance
overload. In fact, the performance drop from single inheritance to multiple inheritance can often
be greater than the performance drop from ordinary to virtual dispatch.
* There is the risk that multiple inheritance can lead to "diamond" inheritance patterns, which are
prone to ambiguity, confusion, and bugs.

##### How Inheritance vs. Composition Is Enforced

The use of inheritance vs. composition is a design decision. Once choice may be made with
consideration for future enhancements. Because of this, most of enforcement must be made using
code review.

#### Operator Overloading

C++ permits user code to declare overloaded versions of the built-in operators using the `operator`
keyword, so long as one of the parameters is a user-defined type. The `operator` keyword also
permits user code to define new kinds of literals using `operator""`, and to define 
type-conversion functions such as `operator bool()`.

Defgine overloaded operators only if their meaning is obvious, unsurprising, and consistent with
built-in operators. Foe example, use `|` as a bitwise- or logical-or, not as a shell-style pipe.

Define operators only on your own types. More precisely, define them in the same headers, cpp
files, and namespaces as the types they operate on. That way, the operators are available
whenever the type is, minimizing the risk of multiple definitions. If possible, avoid defining
operators as templates, because they must satisfy this rule for any possible template arguments.
If you define an operator, also define Any related operators that make sense, and make sure they
are defined consistently.

Prefer to define non-modifying binary operators as non-member functions. If a binary operator is
defined as a class member, implicit conversions will apply  to the right-hand argument, but not
the left-hand one. It will confuse your users if `a + b` compiles, but `b + a` does not.

For a type T whose values can be compared for equality, define a non-member `operator ==` and
document when two values of type T are considered equal. If there is a single obvious notion of
when a value `t1` of type T is less than another value such as `t2`, then you should also define
`operator <=>` which should also be consistent with `operator ==`. Prefer not to overload the other
comparison and ordering operators.

Do not go out of your way to avoid defining operator overloads. For example, prefer to define
`==`, `=`, and `<<` rather than `equals()`, `copyFrom()`, and `printTo()`. Conversely do not define
operator overloads just because other libraries expect them. For example, if your type doesn't have
a natural ordering, but you want to store it in a `std::set`, use a custom comparator rather than
overloading `<`.

Do not overload `&&`, `||`, `,` (comma), or unary `&`. Do not introduce user-defined literals.
Do not use any such literals provided by others (including the standard library).

Type conversion operators are covered in the section on
[implicit conversions](#implicit-conversions). The `=` operator is covered in the section on
[copy constructors](#copyable-and-movable-types). Overloading `<<` for use with streams is covered
in the section on [[[streams]]]. See also the rules on [[[function overloading]]], which apply to
operator overloading as well.

##### Why Use Operator Overloading

* Operator overloading can make code more concise and intuitive by enabling user-defined types to
behave the same as built-in types. Overloaded operators are the idomatic names for certain
operations (e.g. ==, <, +, and <<), and adhering to those conventions can make user-defined types
more readable and enable them to interoperate with libraries that expect those names.
* Usser-defined literals are a very concise notation for creating objects of user-defined types.

##### Why Not To Use Operator Overloading

* Providing a correct, consistent, and unsurprising set of operator overloads requires some care,
and failure to do so can lead to confusion and bugs.
* Overuse of operators can lead to obfuscated code, particularly if the overloaded operator's
semantics don't follow convention.
* The hazards of function overloading apply just as much to operator overloading, if not more so.
* Operator overloads can fool our intuition into thinking that expensive operators are cheap,
built-in operations.
* Finding the call sites for overloaded operators may require a search tool that is aware of C++
syntax, rather than e.g. grep.
* If you get the argument type of an overloaded operator wrong, you may get a different overload
rather than a compiler error. For example: `foo < bar` may do one thing, while `&foo < &bar` does
something totally different.
* Certain operator overloads are inherently dangers. Overloading unary `&` can cause the same code
to have different meanings depending on whether the overload definition is visible. Overloads of
`&&`, `||`, and `'` (comma) cannot match evaluation-order semantics of the built-in operators.
* Operators are often defined outside the class, so ther is a risk of different files introducing
different definitions of the same operator. If both definitions are linked into the same binary,
this results in undefined behavior, which can manifest as subtle runtime bugs.
* User-defined literals allow the creation of new syntactic forms that are unfamiliar even to
experienced C++ programmers, such as `"Hello World"sv` as a shorthand for
`std::string_view("Hello World)`. Existing notations are clearer, though less terse.
* Because they cannot be namespace qualified, uses of user-defined literals also require use of
either using-directives 
[which we ban](#do-not-use-namespace-aliases-at-namespace-scope-in-header-files) or 
using-declarations which [[[we ban in header files - point to Aliases]]] except when the imported
names are part of the interface exposed by the header in question. Given that header files would
have to avoid user-defined literal suffixes, wew prefer to avoid having conventions for literals
differ between header files and source files.

##### How Operator Overloading Rules are Enforced

The only reasonable way to check that these rules are applied is to do a code review.

#### Access Control

Make class data members `private` unless they are [[[constants p- poiint to constant names]]]. 
This simplifies reasoning about
invariants, at the cost of some easy boilerplate in the form of accessors (usually `const`)
if necessary.

For technical reasons, data members of a Google Test test fixture class defined in a .cpp file to
be protected. If a test fixture is defined outside of the .cpp file that it is used in, for example
in a .h file, make the data members `private`.

##### Why Class Data Members Should Be Private

Class data members are part of the private implementation of a class and therefore should not
be visible outside the class.

##### How Access Control is Enforced

clang-tidy provides the rule `misc-non-private-member-variables-in-classes` to check that all class
member variables are `private`.

#### Declaration Order

Group similar declarations together, placing `public` parts at the top, followed by `protected`, and
finally `private` parts of a class declaration. Omit empty sections.

Within each section, prefer grouping similar kinds of declarations together, and prefer the
following order:

1. Types and type aliases (typedefs, using, enum, nested structs and classes, and `friend` types)
2. Optionally and for structs only, non-`static` data members
3. Static constants
4. Factory functions
5. Constructors and assignment operators
6. Destructor
7. All other functions (`static` and non-`static` member functions, and `friend` functions)
8. All other data members (`static` and non-`static`)

Do not put large method definitions inline in the class definition. Usually, only trivial or
performance-critical, and very short, methods may be defined inline. See
[Inline Functions](#inline-functions) for more details.

##### Why Declaration Order is Important

* When similar members are grouped together, it is easier to see similar member types.

##### How Member Order Is Enforced

This is yet another standard that can only currently be enforced through code review.

### Functions

#### Inputs And Outputs

The output of a C++ function is naturally provided via a return value and sometimes via output
parameters or in/out parameters.

Prefer using return values over output parameters: they improve readability, and often provide the
same or better performance.

Prefer to return by value, or failing that, return by reference. Avoid returning a pointer unless
it can be null.

Parameters are either inputs to the function, outputs from the function, or both. Non-optional
input parameters should usually be values or `const` references, while non-optional output
parameters should usually be references which cannot be null. Generally, use `std::optional` to
represent optional by-value inputs, ans use a `const` pointer when the non-optional form would
have used a reference. Use non-const pointers to represent optional outputs and optional
input/output parameters.

Avoid defining functions that require a `const` reference parameter to outlive the call, because
`const` reference parameters bind to temporaries. Instead, find a way to eliminate the lifetime
requirement, for example, by copying the parameter, or pass it by `const` pointer and document
the lifetime and non-null requirements.

When ordering function parameters, put all input-only parameters before any output parameters. In
particular, do not add new parameters to the end of the function just because they are new; place
new input-only parameters before the output parameters. This is not a hard-and-fast rule.
Parameters that are both input and output muddy the waters, and, as always, consistency with
related functions may require you to bend the rule. Variadic functions may also require unusual
parameter ordering.

##### How Input And Output Rules Are Enforced

Code review must be performed to check that the input and output rules are enforced.

#### Write Short Functions

Prefer small and focused functions.

Sometimes long functions are appropriate, so no hard limit is placed on the length of functions. If
a function exceeds about 40 lines, think about whether it could be broken up without harming the
structure of the program.

##### Why Write Short Functions

* Even if your long function works perfectly now, someone modifying it in the future may add new
behavior. This could result in bugs that are hard to find. Keeping functions short and simple
makes it easier for other people to read and modify the code. 
* Small functions are easier to test.

##### How Write Short Functions Is Enforced

Because long functions are occasionally required, and it is difficult to specify the number of lines
that separate a short function from a long function, code review must be used to enforce short
functions.

#### Function Overloading

Use overloaded functions, including constructors, only if a reader looking at a call site can get
a good idea of what is happening without first having to figure out exactly which overload is
being called.

You can overload a function when there are no semantic differences between variants. These overloads
may vary in types, qualifiers, or argument count. However, a reader of such a call must not need to
know which member of the overload set is chose, only that something from the set is being called. If
you can document all entries in the overload set with a single comment in the header, that is a good
sign that it is a well-defined overload set.

##### Why Use Function Overloading

* Overloading an make code more intuitive by allowsing an identically-named function to take
different arguments. It may be necessary for templatized code, and it can be inconvenient for
Visitors.
* Overloading based on `const` or ref qualification may make code more usable, more efficient, or
both.

##### Why Not To Use Function Overloading

* If a function is overloaded by argument types alone, a reader may have to understand C++'s
complex matching rules in order to tell what's going on. 
* Many people are confused by the
semantics of inheritance if a derived class overrides only some of the variants of a function.

##### How Function Overloading Rules Are Enforced

Code review is the only option.

#### Default Arguments

Default arguments are allowed on non-virtual functions when the default is guaranteed to always
have the same value. Follow the same restrictions as for 
[function overloading](#function-overloading) and prefer overloaded functions if the readability
gained with default arguments does not outweigh  the downsides listed below.

Default arguments are banned on virtual functions, where they don't work properly, and in cases
where the specified default might not evaluate to the same value depending on when it is
evaluated. For example, don't write `void f(int n = counter++);`

In some other cases, default arguments can improve the readability of their function declarations
enough to overcome the downside listed below. When in doubt, use overloads.

##### Why Use Default Arguments

* Ofter you have a function that uses default values, but occasionally you want to override the
defaults. Default parameters allow an easy way to do this without having to define many functions
for the rare exceptions. Compared to overloading the function, default arguments have a cleaner
syntax, with less boilerplate and a cleaner distinction between "required" and "optional"
arguments.

##### Why Not To Use Default Arguments

* Default arguments are another way to achieve the semantics of overloaded functions, so all the
[reasons not to overload functions](#function-overloading) apply.
* The defaults for arguments in a virtual function call are determined by the static type of the
target object, and there's no guarantee that all overrides of a given function declare the same
defaults.
* Default parameters are re-evaluated at each call site, which can bloat the generated code. Readers
may also expect the default value to be fixed at the declaration instead of varying at each call.
* Function pointers are confusing in the presence of default arguments, since the function
signature often does not match the call signature. Adding function overloads avoids these problems.

##### How Default Argument Rules Are Enforced

The only way that default argument rules can be checked is via code review.

#### Trailing Return Type Syntax

Use trailing return types.

 Here is an example of the leading syntax:

```C++
int foo(int x);
```

The newer form uses the `auto` keyword before the function name and a trailing return type after
the argument list:

```C++
auto foo(int x) -> int;
```

The trailing return type is in the function's scope. This does not make a difference for a simple
case like `int` but it matters for more complicated cases, like types declared in class scope or
types written in terms of function parameters.

##### Why Use Trailing Return Type Syntax

* Trailing return type syntax is the only way to explicitly specify the return type of a 
[lambda expression](#lambda-expressions). In some cases, the compiler is able to deduce a lambda's 
return type, but
not in all cases. Even when the compiler can deduce it automatically, sometimes specifying it
explicitly would be clearer for readers.
* Sometimes it is easier and more readable to specify a return type after the function's
parameter list has already appeared. This is particularly tye when the return type depends on
template parameters.

##### Why Not To Use Trailing Return Type Syntax

* Trailing return type syntax is relatively new and it has no analog in C++-like languages like
C and Java, so readers may find it unfamiliar.
* Existing code bases have an enormous number of function declarations that are not going to get
changed to the new syntax. Using a uniformity of style may be better for readability.
* clang-tidy is used to check for trailing return type. Unfortunately, clang-tidy only addresses
the contents of .cpp files and not .h files. Consequently, using clang-tidy to check that
trailing return type is used only checks in the .cpp files and ignores .h files.

##### How Trailing Return Type Syntax Is Enforced

The clang-tidy rule `modern-use-trailing-return-type` flags all occurrences of not using trailing
return type syntax. Because clang-tidy does not process header files, .h files must be checked
using code review.

## Other C++ Features

### Ownership and Smart Pointers

If dynamic ownership is necessary, prefer to keep ownership with the code that allocated it. If
other code needs access to the object, consider passing it a copy, or passing a pointer or
reference without transferring ownership. Prefer to use `std::unique_ptr` to make ownership
transfer explicit.

Do not design your code to use shared ownership without a good reason. One such reason is to
avoid expensive copy operations, but you should only do this if the performance benefits are
significant, and the underlying object is immutable. For example: `std::shared_ptr<const Foo>`.
If you do use shared ownership, prefer to use `std::shared_ptr`.

Never use `std::auto_ptr`. Use `std::unique_ptr` instead.

Note that code which predates C++11 did not have `std::unique_ptr` and `std::shared_ptr` If you
must call such code, hopefully the documentation states who owns the underlying object.

#### Why Use Smart Pointers

* It is virtually impossible to manage dynamically allocated memory without some sort of ownership
logic.
* Transferring ownership of an object can be cheaper than copying it, if copying is even possible.
* Transferring ownership can be simpler than "borrowing" a pointer or reference, because it
reduces the need to coordinate the lifetime of the object between the two users.
* Smart pointers can improve readability by making ownership logic explicit, self-documenting, and
unambiguous.
* Smart pointers can eliminate manual ownership bookkeeping, simplifying the code and ruling out
large classes of errors.
* For `const` objects, shared ownership can be a simple and efficient alternative to deep copying.

#### Why Not To Use Smart Pointers

* Ownership must be represented and transferred via pointers, whether smart or plain. Pointer
semantics are more complicated than value semantics, especially in APIs; you have to worry not
just about ownership. but also aliasing, lifetime, and mutability, among other issues.
* The performance costs of value semantics is often overestimated, so the performance benefits of
ownership transfer might not justify the readability and complexity costs.
* APIs that transfer ownership force their clients into a single memory management model.
* Code using smart pointers is less explicit about where the resource releases take place.
* `std::unique_ptr` expresses ownership transfer using move semantics, which are relatively new
and may confuse some programmers.
* Shared ownership can be a tempting alternative to careful ownership design, obfuscating the
design of a system.
* Shared ownership requires explicit bookkeeping at runtime, which can be costly.
* In some cases, such as cyclic references, objects with shared ownership may never be deleted.
* Smart pointers are not perfect substitutes for plain pointers.

#### How Ownership and Smart Pointers Are Enforced

The bookkeeping required of smart pointers can be ver complicated. The only way to possibly enforce
the requirements and rules is through code review.

### Rvalue References

Use rvalue references only as follows:

* To define move constructors and move assignment operators as defined in
[Copyable and Movable Types](#copyable-and-movable-types).
* To define `&&`-qualified methods that logically "consume" `*this`, leaving it in an unusable or
empty state. This applies only to method qualifiers which come after the closing parenthesis of a
function signature; if you want to "consume" an ordinary function parameter, prefer to pass it by
value.
* To use forwarding references in conjunction with `std::forward`, to support perfect forwarding.
* To define pairs of overloads, such as one taking `Foo&&` and the other taking `const Foo&&`.
Usually the preferred solution is just to pass by value, but an overloaded pair of functions
sometimes yields better performance and is sometimes necessary in generic code that needs to
support a wide variety of types. If you are writiing more complicated code for the sake of
performance, make sure you have evidence that it actually helps.

#### Why Use Rvalue References

* Defining a move constructor makes it possible to move a value instead of copying it. If `v1` is a
`std::vector<std::string>`, for example, then `auto v2(std::move(v1))` will probably just result in
some pointer manipulation instead of copying a large amount of data. In many cases, this can result
in a major performance improvement.
* Rvalue references make it possible to implement types that are movable but not copyable, which can
be useful for types that have no sensible definition of copying but where you might still want to
pass them as function arguments, put them in containers, etc.
*`std::move` is necessary to make effective use of some standard library types, such as
`std::unique_ptr`.
* Forwarding references which use the rvalue reference token make it possible to write a generic
function wrapper that forwards its arguments to another function, and works whether or not its
arguments are temporary objects and/or `const`. This is called "perfect forwarding".

#### Why Not To Use Rvalue References

* Rvalue references are not yet widely understood. Rules like reference collapsing and the special
deduction rule for forwarding references are somewhat obscure.
* Rvalue references are often misused. Using rvalue references is counter-intuitive in signatures
where the argument is expected to have a valid specified state after the function call, or where
no move operation is performed.

#### How Rvalue References Rules Are Enforced

Aside from the enforcement methods for [movable types](#copyable-and-movable-types), the only method
that can be used to enforce rvalue references is code review.

### Friends

`friend` classes and functions may be used within reason.

Friends should be defined in the same file so that the reader does not have to look into another
file to file uses of the private members of a class.

#### Why Use Friends

* A common use of `friend` is to have a `FooBuilder` class be a friend of `Foo` so that it can
construct the inner state of `Foo` correctly without exposing this state to the world. In some
cases, it may be useful to make a unit test class a friend of the class it tests.
* Friends extend but do not break the encapsulation boundard of a class. In some cases this is
better than making a member `public` when you want to give only one other class access to it.
However, most classes should interact with other classes solely through their public members.

#### How The Friends Standard Can Be Enforces

The only way to enforce the `friend`s standard is through code review.

### Exceptions and Error Handling

There are a number of different ways of handling error conditions. For example, see
[Old School Ways of Handling Status and Error Conditions](https://computingonplains.wordpress.com/old-school-ways-of-handling-status-and-error-conditions)
and
[Modern Ways of Handling Status and Error Conditions](https://computingonplains.wordpress.com/modern-ways-of-handling-status-and-error-conditions/).

Error conditions should be handled as close to their source as possible, with exceptions limited
to states that require either program termination or lengthy recovery procedures, or
where exceptions are thrown in external libraries that are used in the project.

The benefits of using exceptions, especially in new projects, appear to outweigh the costs.
However, for existing code without exceptions, the introduction of exceptions has implications on
all dependent code. If exceptions can be propagated beyond a new project, it also becomes
problematic to integrate the new project into existing exception-free code.

The alternatives to exceptions, such as error codes and assertions, do not appear to introduce a
severe burden and therefore should be considered instead.

#### Why Use Exceptions

* Exceptions allow higher levels of an application to decide how to handle "can't happen" failures
in deeply nested functions, without the obscuring and error-prone bookkeeping of error codes.
* Exceptions are used by most other modern languages. Using them in C++ would make it more
consistent with Python, Java, C#, and the C++-like languages that others are familiar with.
* Some third-party libraries use exceptions, and turning them off internally makes it harder to
integrate with those libraries.
* Exceptions are the only way for a constructor to fail. We can simulate this with a factory
function or an `Init()` method, but these require heap allocation or a new "invalid" state,
respectively.
* Exceptions are really handy in test frameworks.

#### Why To Not Use Exceptions

* When you add a `throw` statement to an existing function, you must examine all of its
transitive callers. Either they must make at least the basic exception safety guarantee, or they
must never catch the exception and be happy with the program terminating as a result. For instance,
if `f()` calls `g()` calls `h()`, and `h` throws an exception that `f` catches, `g` has to be
careful or it may not clean up properly.
* More generally, exceptions make the control flow of programs difficult to evaluate by looking at
code: functions may return in places you don't expect. This causes maintainability and debugging
difficulties. You can minimize this cost via some rules on how and where exceptions can be used,
but at the cost of more that the developer needs to know and understand. See, for example:
[C++ Exceptions and wxWidgets](https://computingonplains.wordpress.com/cpp-exceptions-and-wxwidgets/).
* Exception safety requires both RAII and different coding practices. Lots of supporting
machinery is needed to make writing correct exception-safe code easy. Further, to avoid requiring
readers to understand the entire call graph, exception-safe code must isolate logic that writes
to persistent state into a "commit" phase. This will have both benefits and costs, perhaps when
you are forced to obfuscate code to isolate the commit. Allowing exceptions forces us to always pay
those costs even when they are not worth it.
* Turning on exceptions adds data to each binary produced, increasing compile time, probably
slightly, and possibly increasing address space pressure.
* The availability of exceptions may encourage developers to throw them when they are not
appropriate or recover from them when it is not safe to do so. For example, invalid user input
should not cause exceptions to be thrown. We would need to make the coding standard even longer to
document these restrictions.

#### How Error and Exception Handling Is Enforced

The proper use of error codes and exception handling can only be enforced through the use of
code reviews.

### `noexcept`

Specify `noexcept` when it is useful for performance if it accurately reflects the intended
semantics of the function, i.e. that if an exception is somehow thrown from within the function
body then it represents a fatal error. You can assume that `noexcept` on move constructors has a
meaningful performance benefit.

Prefer unconditional `noexcept` if exceptions are completely disabled. Otherwise, while using
conditional `noexcept` with simple conditions may be appropriate in a few cases, it may be simpler
and easier for a reader to understand if the `noexcept` specifier is simply not used. Note that in
many cases, the only possible cause for an exception is allocation failure, and it is almost
certainly most appropriate to terminate the application rather than attempt to recover from this
exception.

Instead of writing a complicated `noexcept` clause that depends on whether a hash function can
throw, for example, simply document that the component doesn't support hash functions throwing
and make it unconditionally `noexcept`.

#### Why Use `noexcept`

* Specifying move constructors as `noexcept` improves performance in some cases, e.g.
`std::vector<T>::resize()` moves rather than copies the objects if `T`'s move constructor is
`noexcept`.
* Specifying `noexcept` on a function can trigger compiler optimizations in environments where
exceptions are enabled, e.g. compiler does not have to generate extra code for stack-unwinding, if
it knows that no exceptions can be thrown due to a `noexcept` specifier.

#### Why Not To Use `noexcept`

* In projects following this coding standards document that have exceptions disabled, it is hard
to ensure that `noexcept` specifiers are correct, and hard to define what correctness even means.
* It is hard, if not impossible, to undo `noexcept` because it eliminates a guarantee that callers
may be relying on, in ways that are hard to detect.

#### How Use Of `noexcept` Is Enforced

Correct use of `noexcept` can be checked using code review. Otherwise, incorrect use may be seen
from unexpected program termination.

### Run-Time Type Information (RTTI)

RTTI allows a programmer to querey the C++ class of an object at run-time. This is done by use of
`tyepid` or `dynamic_cast`.

RTTI has legitimate uses but is prone to abuse, so you must be careful when using it. It may be
used freely in unittests, but aviod it when possible in other code. If you find yourself needing
to write code that behaves differently based on the class of an object, consider one of the
following alternatives to query the type:

* Virtual methods are the preferred way of executing different code paths depending on a specific
subclass type. This puts the work within the object itself.
* If the work belongs outside the object and instead in some processing code, consider a
double-dispatch solution, such as the Visitor design pattern. This allows a facility outside the
object itself to determine the type of the class using the built-in type system.

When the logic of a program guarantees that a given instance of a base class is in fact an
instance of a particular derived class, then a `dynamic_cast` may be used freely on the object.
Usually one can use a `static_cast` as an alternative in such situations.

Decision trees based on type are a strong indication that your code is on the wrong track.

```C++
if (typeid(*data) == typeid(D1)) {
    ...
} else if (typeid(*data) == typeid(D2)) {
    ...
} else if (typeid(*data) == typeid(D3)) {
    ...
}
```

Code such as this usually breaks when additional subclasses are added to the class hierarchy.
Moreover, when properties of a subclass change, it is difficult to find and modify all the
affected code segments.

Do not hand-implement an RTTI-like workaround. The arguments against RTTI apply just as much to
workarounds like class hierarchies with type tags. Moreover, workarounds disguise your true intent.

#### Why Use RTTI

* The standard alternatives to RTTI require modification or redesign of the class hierarchy in
question. Sometimes such modifications are infeasible or undesirable, particularly in
widely-used or mature code.
* RTTI can be used in some unit tests. For example, it is useful in tests of factory classes where
the test has to verify that a newly created object has the expected dynamic type. It is also
useful in managing the relationship between objects and their mocks.
* RTTI is useful when considering multiple abstract objexts. For example:

```C++
bool Base::Equal(Base* other) = 0;
bool Derived::Equal(Base* other) {
    Derived* that = dynamic_cast<Derived*>(other);
    if (that == nullptr) {
        return false;
        ...
    }
}
```
### Why Not To Use RTTI

* Querying the type of an object at run-time frequently means a design problem. Needing to know
the type of an object at run-time is often an indication that the design of your class hierarchy
is flawed.
* Undisciplined use of RTTI makes code hard to maintain. It can lead to type-based decision trees
or switch statements scattered throughout the code, all of which must be examined when making 
further changes.

#### How Use of RTTI Is Enforced

At this time, there are no checks in clang-tidy for the presence of RTTI, so the only way to
enforce proper use of RTTI is through code review.

### Casting

In general, do not use C-style casts. Instead use these C++-style casts when explicit type
conversion is necessary:

* Use brace initialization to convert arithmetic types (e.g. int64_t{x}). This is the safest
approach because code will not compile if conversion can result in information loss. This syntax
is also concise.
* Use `static_cast` as the equivalent of a C-stype cast that does value conversion, when you need to
explicitly up-cast a pointer froma superclass to a subclass. In this last case, you must be sure
that your object is actually an instance of the subclass.
* Use `const_cast` to remove the `const` qualifier. See [[[const]]].
* Use `reinterpret_cast` to do unsafe conversions of pointer types to and from integer and other
pointer types, including `void*`. Use this only if you know what you are doing and you understand
the aliasing issues.

See the [RTTI section](#run-time-type-information-rtti) for guidance on the use of `dynamic_cast`.

#### Why Use C++-Style Casting

* The problem with C-style casts is the ambiguity of the operation; sometimes you are doing a
conversion (e.g. (int)3.5), and sometimes you are doing a cast (e.g. (int)"hello"). Brace
initialization and C++ casts can often help avoid this ambiguity. Additionally, C++ casts are more
visible when searching for them.

#### Why Not To Use C++-Style Casting

The C++-style cast syntax is verbose and cumbersome.

#### How C++-Style Casting Is Enforced

The only way to enforce C++-style casting is through the use of code reviews.

### Streams

Streams are the standard I/O abstraction in C++, as exemplified by the standard header `<iostream>`.

Use streams when they are the best tool for the job. Be consistent with the code around you, and
with the codebase as a whole; if there is an established tool for your problem, use that tool
instead. In particular, logging libraries are usually better than `std::cerr` and `std::clog` for
diagnostic output.

When using streams, avoid the stateful parts of the streams API other than error state, for example:
`imbue()`, `xalloc()`, and `register_callback()`.

Overload `<<` as a streaming operator for your type only if your type represents a value, and `<<`
writes out a human-readable string representation of that value. Avoid exposing implementation
details in the output of `<<`; if you need to print object internals for debugging, using named
functions such as `DebugString()` is the most common convention.

#### Why Use Streams

* The `<<` and ``>>` stream operators provide an API for formatted I/O that is easily learned,
portable, reusable, and extensible. `printf`, by contrast, doesn't support `std::string`, to say
nothing of user-defined types, and is difficult to use portably. `printf` also obliges you to
choose among the numerous slightly different versions of that function, and navigate the dozens
of conversion specifiers.

### Why Not To Use Streams

* Stream formatting can be configured by mutating the state of the stream. Such mutation are
persistent, so the behavior of your code can be affected by the entire history of the stream,
unless you go out of your way to restore it to a known state every time other code might have
touched it. User code can not only modify the built-in state, it can add new state variables and
behavior through a registration system.
* It is difficult to precisely control stream output, due to the above issues, the way code and
data are mixed in streaming code, and the use of operator overloading which may select a different
overload than you expect.
* The practice of building up output through chains of `<<` operators interferes with
internationalization, because it bakes word order into the code, and streams' support for
localization is flawed. 
* Internationalization and locale support can add complications to stream output. For example:
  * Windows, MacOS, and Linux all use different methods for substituting strings of one language 
  for another.
  * The order of values to be printed may vary from one language to another based on the translated
  string.
  * Each locale has a different set of monetary, date/time, numeric, and parsing requirements. See
  the [C++ localization library](https://en.cppreference.com/w/cpp/locale) 
  for an idea of how complicated this can be.
  * Some languages are written from right to left rather than left to right.

#### How Stream Handling Is Enforced

The only sure way to enforce proper handling of streams is through code review.

### Preincrement and Predecrement

Use prefix increment and decrement unless the code explicitly needs the result of the postfix
increment or decrement exression.

#### Why Use Preincrement and Predecrement

* A postfix increment/decrement expression evaluates to the value as it was before it was modified.
This can result in code that is more compact but harder to read. The prefix format is generally
more readable, is never less efficient, and can be more efficient because it doesn't need to make
a copy of the values as it was before the operation.

#### Why Not To Use Preincrement and Predecrement

* The tradition developed, in C, of using postincrement, even when the expression is not used,
especially in `for` loops.

#### How Preincrement and Predecrement Are Enforced

As there may be a reason for using postincrement or postdecrement, the only way to ensure that
postincrement and postdecrement are used only when needed is through code review.

Note that the use of preincrement and predecrement is a convention and may or may not have any
effect on efficiency.

### Use of Const

The use of `const` in APIs (e.g. on function parameters, methods, and non-local variables)
is strongly recommended whenever it is meaningful and accurate. This provides consistent, mostly
compiler-verified documenation of what objects an operation can mutate. Having a consistent and
reliable way to distinguish reads from writes is critical to writing thread-safe code, and is
useful in many other contexts as well. In particular:

* If a function guarantees that it will not modify an argument passed by reference or pointer, the
corresponding function parameter should be a reference-to-const (`const T&`) or pointer-to-const
(`const T*`), respectively.
* For a function parameter passed by value `const` has no effect on the caller, thus it is not
recommended in function declarations.
* Decalre methods to be const unless they alter the logical state of the object or enable the user
to modify that state, e.g. by returning a non-`const` reference, but that is rare, or they cannot
be safely invoked concurrently.

Using `const` on local variables is neither encouraged nor discouraged.

All of a class's `const` operations should be safe to invoke concurrently with each other. If that
is not feasible, the class must be clearly documented as "thread-unsafe".

#### Where To Put The `const`

Some people favor the form `int const *foo` to `const int* foo`. They argue that this is more
readable because it is more consistent: it keeps the rule that const always follows the obejct it is
describing. However, this consistency argument does not apply in codebases with few deeply-nesteds
pointer expressions since most `const` expressions have only one `const`, and it applies to the
underlying value. In such cases, there is not consistency to maintain. Putting the `const` first
is arguably more readable, since it follows English in putting the "adjective" (`const`) before
the noun (`int)`).

That said, while we encourage putting `const` first, we do not require it. Just be consistent with
the code around you.

#### Why Use `const`

* It is easier for people to understand how variables are being used.
* Allows the compiler to do better type checking, and conceivably, to generate better code.
*It helps people know shat functions are safe to use without locks in multithreaded programs.

#### Why Not To Use `const`

`const` is viral: if you pass a `const` variable to a function, that function must have `const` in
its prototype or the variable will need a `const_cast`. This can be a particular problem when
calling library functions.

#### How Use Of `const` Is Enforced

The following clang-tidy checks apply to the use of `const`:

* cppcoreguidelines-avoid-const-or-ref-data-members
* cppcoreguidelines-avoid-non-const-global-variables
* cppcoreguidelines-pro-type-const-cast

Additional checks must be performed using code review.

### Use Of `constexpr`, `constinit`, And `consteval`

`constexpr` definitions enable a more robust specification of the constant parts of an interface.
Use `constexpr` to specify true constants and the functions that support their definitions.

`consteval` may be used with code that must not be invoked at run-time.

Avoid complexifying function definitions to enble their use with `constexpr`.

Use `constinit` to ensure constant initialization for non-constant variables.

Do not use `constexpr` or `consteval` to force inlining.

#### Why Use `constexpr`

* Use of `constexpr` enables definition of constants with floating-point expressions rather than
just literals, definition of constants of user-defined types, and definition of constants with
function calls.

#### Why Not To Use `constexpr`

* Prematurely making something as `constexpr` may cause migration problems if later on it has to
be downgraded.
* Current restrictions on what is allowed in `constexpr` functions and constructors may invite
obscure workarounds in these definitions.

### Integer Types

C++ does not specify exact sizes for the integer types like `int`. Common sizes on contemporary
architectures are 16 bits for `short`, 32 bits for `int`, 32 or 64 bits for `long`, and 64 bits for
`long long`, but different platforms make different choices, in particular for `long`.

The standard library header `<cstdint>` defines types like `int16_t`, `uinit32_t`, `int64_t` and so
forth. You should always use these in preference to `short`, `unsigned long long`, and the like
when you need a guarantee on the size of an integer. Omit the `std::` preefix for these types as
the extra 5 characterrs do not merit the added clutter. Of the built-in types, only `int` should be
used. When appropriate, you may use the standard type aliases like `size_t` and `ptrdiff_t`.

Most code uses `int` very often for integers that we know are not going to be too big, such as for
loop counters. You can assume that `int` is at least 32 bits, but do not assume that it has more
than 32 bits. If you need a 64 bit integer type, use `int64_t` or `uint64_t`.

Do not use unsigned integer types such as `uint32_t` unless there is a valid reasojn for
representing a bit pattern rather than a number, or you need defined overflow modulo 2^32. In
particular, do not use unsigned types to say a number is never negative. Instead, use assertions
for this.

If your code is a container that returns a size, be sure to use a type that will accommodate any
possible usage of your container. The use of `size_t` is recommended.

Use care when converting integer types. Integer conversions and promotions can cause undefined
behavior, leading to security bugs and other problems.

#### When To Use Unsigned Integers

* Unsigned integers are good for representing bitfields and modular arithmetic. Because of
historical accident, the C++ standard also uses unsiogned integers to represent the size of
containers, but it is impossible to fix at this point.
* The fact that unsigned arithmetic does not model the behavior of a simple integer, but is instead
defined by the standard to model modular arithmetic, wrapping around on overflow and underflow,
means that a significant class of bugs cannot be diagnosed by the compiler. In other cases, the
defined behavior impedes optimization.
* Mixing signedness of integer types is responsible for an equally large class of problems. The
best advice is to:
  * Use iterators and containers rather than pointers and sizes.
  * Try not to mix signedness.
  * Avoid unsigned types except for representing bitfields or modular arithmetic.
  * Do not use an unsigned type merely to assert that a variable is non-negative.

#### Why Use Types Defined in `<cstdint>`

* Uniformity of declaration.

#### How To Enforce `<cstdint>` Types Use

clang-tidy does not currently have any rules related to integer type definition. Therefore the
only way to enforce integer type definitions is through the use of code review.

### 64-bit Portability

Code should be developed and run on 64 bit systems including Windows, MacOS, and Linux. Support
for other operating systems is encouraged but not required. 32-bit-specific code should not be
added to the code base. If you must share data between 32-bit and 64-bit systems, keep the following
points in mind:

* Use `std::ostream` for writing data to a file, never `printf`. Correct portable `printf()`
conversion specifiers for some integral typedefs rely on macro expansions that are unpleasant to
use and impractical to require.
* `sizeof(void *) != sizeof(int)`. Use `intptr_t` instead if you want a pointer-sized integer.
* Be careful of structure alignments, especially those stored on disk. Any class/structure with a
`int64_t/uint64_t` member will by default end up being 8-bit aligned on a 64-bit system. If you
have such structures being shared on disk between 32-bit and 64-bit code, you need to ensure that
they are packed the same on both architectures. Most compilers offer a way to alter structure
alignment.
* Use [brace-initialization](#casting) as needed to create 64-bit constants. For example:

```C++
int64_t aValue{0X123456789};
uint64_t aMask{uint64_t{3} << 48};
```

#### Why Support Only 64-Bit Systems

* 32-bit systems are becoming very rare. The main operating systems that the software for which
this coding standard applies are increasingly removing support for 32-bit systems.
* The software projects that this coding standard apply to are "greenfield" projects, so no
sharing of data between 32-bit and 64-bit systems is envisaged.

#### How 64-Bit Portability Is Enforced

There is no specific mechanism to enforce 64-bit code. In some cases, 32-bit code may not compile,
or there may not be 32-bit libraries on the build system. Even if 32-bit software is built, it may
not function properly. It is up to the developer to determine if support for such problems will be
provided.

### Preprocessor Macros

Avoid defining macros, especially in header files. Prefer inline functions, `enum`s, and `const` or
`constexpr` variables.

If you must define a macro, name it with a project-specific prefix.

Do not use macros to define pieces of a C++ API.

The following usage pattern will avoid many problems with macros; if you must use macros, follow
these whenever possible:

* Do not define macros in a .h file.
* `#define` macros immediately before you use them, and #undef then immediately after.
* Do not `#undef` an existing macro before replacing it with your own; instead, pick a name
that is likely to be unique.
* Try not to use macros that expand to unbalanced C++ constructs, or at least document that
behavior well.
* Prefer not using `##` to generate function/class/variable names.

#### Why Use Preprocesor Macros

* System-defined macros define the operating system and architecture for the system being built.
These macros may be required to separate operating system or architecture specific code. This is
strongly discouraged but may be necessary in some instances.
* Macros can do things other techniques cannot. Some of their special features, like stringifying,
concatenation, and so forth are not otherwise available through the C++ language. Before using
a macro, consider carefully whether there is a non-macro way to achieve the same result, and use
that.
* Macros are used heavily in unittest frameworks.

#### Why Not To Use Preprocessor Macros

* When macros are used, the code you see is not the code the compiler sees. This can introduce
unexpected behavior, especially since macros have global scope.
* The problems introduced by macros are especially severe when they are used to define pieces of a
C++ API, and even more so for public APIs. Every error message from the compiler when developers
incorrectly use that interface now must explain how the macros formed the interface. Refactoring
and analysis tools have a dramatically harder time updating the interface.
* Testing is much more difficult.
* Modern C++ provides a number of features that replace the need for macros:
  * Instead of using a macro to inline performance critical code, use an inline function.
  * Instead of using a macro to store a constant, use a `const` or `constexpr` variable.
  * Instead of using a macro to abbreviate a long variable name, use a reference.
  * Instead of groups of adjacent macros, use an `enum`. For new code, use an `enum class`. In
  existing code, replace the macros with an unscoped anonymous `enum`.
  * Instead of using a macro to conditionally compile code..., well don't do that, except of course,
  for the [include guards](#22-include-guards) to prevent multiple inclusion of header files.

### How Preprocessor Macro Use Is Enforced

clang-tidy has the following rules related to macros:

* modernize-macro-to-enum
* modernize-replace-disallow-copy-and-assign-macro
* bugprone-macro-parentheses
* bugprone-macro-repeated-side-effects
* bugprone-multiple-statement-macro

Code review is needed to ensure that macros are otherwise defined and used according to the
requirements listed above.

### '\0' And `nullptr/NULL`

Use `nullptr` for pointers, and '\0' for chars. Do not use the `0` literal.

#### Why Use '\0' And `nullptr`

* Using `nullptr` for pointers provides type-safety.
* Using '\0' for the null character makes the code more readable.

#### How Use Of '\0' and `nullptr` Is Enforced

clang-tidy has the following rules:

* modernize-use-nullptr
* bugprone-stringview-nullptr

Code review is required to ensure that '\0' and `nullptr` are used instead of `0`` and `NULL`.

### `sizeof`

Prefer `sizeof(varname)` to `sizeof(type)`.

Use `sizeof(type)` for code unrelated to any particular variable, such as code that manages an
external or internal data format where a variable of an appropriate C++ type is not convenient.

#### Why Use `sizeof(varname)`

* `sizeof(varname)` will update appropriately if someone changes the variable type either now or
later. 

#### How `sizeof()` Is Enforced

clang-tidy has the following rule related to `sizeof()`:

* bugprone-sizeof-container
* bugprone-sizeof-expression

These rules check for very specific questionable uses of `sizeof()`. The only way to check that
`sizeof(varname)` and `sizeof(type)` are used appropriately is through code review.

### Type Deduction (Including `auto`)

Use type deduction only to make the code clearer or safer. Do not use it merely to avoid the
inconvenience of writing the explicit type. When judging whether the code is clearer, keep in mind
that your readers are not necessarily on your team or familiar with the project, so types that you
or your reviewer experience as unnecessary clutter will very often provide useful information to
others. For example, you can assume that `make_unique<Foo>()` is obvious, but the return type of
`MyWidgetFactory()` probably isn't.

These principles apply to all forms of type deduction, but the details vary, as described in the
following sections:

#### Type Deduction Details

##### Function Template Argument Deduction

Function template argument deduction is almost always OK. Type deduction is the expected default
way of interacting with function templates because it allows function templates to act like
infinite sets of ordinary function overloads. Consequently, function templates are almost always
designed so that template argument deduction is clear and safe, or does not compile.

##### Local Variable Type Deduction

For local variables, you can use type deduction to make the code clearer by eliminating type
information that is obvious or irrelevant, so that the reader can focus on meaningful parts of the
code. Here are some examples:

```C++
auto widget = std::make_unique<WidgetWithBellsAndWhistles>(arg1,arg2);
auto it = myMap.find(key);
std::array numbers = { 1, 2, 4, 8, 16, 32 };
```

is clearer than:

```C++
std::unique_ptr<WidgetWithBellsAndWhistles> widget = 
    std::make_unique<WidgetWithBellsAndWhistles>(arg1,arg2);
std::unordered_map<std::string, std::make_unique<WidgetWithBellsAndWhistles>>::const_iterator
    it = myMap.find(key);
std::array numbers =  { 1, 2, 4, 8, 16, 32 };
```

Types sometimes contain a mixture of useful information and boilerplate, such as `it` above. It is
obvious that the type is an `iterator`, and in many contexts the container type and even the key
are not relevant, but the type of values is probably useful. In such situations, it is often
possible to define local variables with explicit types that convey the relevant information:

```C++
if (auto it = myMap.find(key); it != myMap.end()) {
    WidgetWithBellsAndWhistles& widget = *it->second;
    // Do stuff with `widget`
}
```
If the type is a template instance ande the parameters are boilerplate but the template itself is
informative, you can use class template argument deduction to suppress the boilerplate. However,
cases where this actually provides a meaningful benefit are quite rare. Note that class template
argument deduction is also subject to a [[[separate standard - #class-template-argument-deduction]]].

Do not use `decltype(auto)` if a simpler option will work because it is a fairly obscure feature
which therefore has a high cost in code clarity.

##### Return Type Deduction

Return type deduction should not be used as it contradicts the 
[trailing return type syntax](#trailing-return-type-syntax) standard.

##### Parameter Type Deduction

`auto` parameter types for lambdas should be used with caution because the actual type is
determined by the code that calls the lambda, rather than by the definition of the lambda.
Consequently, an explict type will almost always be clearer unless the lambda is explicitly called
very close to where it is defined, or the lambda is passed to an interface so well-known that it is
obvious what arguments it will eventually be called with (e.g. `std::sort`).

##### Lambda init captures

Init captures are covered by a [[[more specific standard = #lambda-expressions]]] which largely
supercedes general rules for type deduction.

##### Structured Bindings

Unlike other forms of type deduction, structured bindings can actually give the reader additional
information by giving meaningful names to the elements of a larger object. This means that a
structured binding declaration may provide a net readability improvement over an explicit type,
even in cases where `auto` would not. Structured bindings are especially beneficial when the object
is a `pair` or `tuple` because they don't have meaningful names to begin with. Note that you 
[should not use pairs or tuples](#structs-vs-pairs-and-tuples) unless a pre-existing API like
`insert` forces you to.

If the object being bound is a `struct`, it may sometimes be helpful to provide names that are
more specific to your usage, but keep in mind that this may also mean the names are less
recognizable to your reader than the field names. We recommend using a comment to indicate the
name of the underlying field if it does not match the name of the binding, using the same
syntax as for function parameter comments:

```C++
auto [/*fieldName1=*/boundName1, /*fieldName2=*/boundName2] = ...
```

As with function parameter comments, this can enable tools to detect if you get the oprder of the
fields wrong.
##### How Type Deduction Rules Are Enforced

clang-tidy has two checks related to type deduction:

* modernize-replace-auto-ptr
* modernize-use-auto

The second check especially handles quite a few cases such as those discussed in the subsections
above.

For cases that are not or cannot be handled by clang-tidy, code review is the only option.

#### Class Template Argument Deduction

Class template argument deduction occurs when a variable is declared with a type that names a
template, and the template argument list is not provided, not even empty angle brackets:

```C++
std::array a = { 1, 2, 3};  // `a` is of type std::array<int, 3>
```
The compiler deduces the arguments from the initializer using the template's "deduction guides".
When you declare a variable that relies on class template argument deduction, the compiler selects
a deduction guide using the rules of constructor overload resolution, and that guide's return
type becomes the type of the variable.

Do not use class template argument deduction with a given template unless the template's
maintainers have opted into supporting its use by providing at least one explicit deduction guide.
All templates in the `std` namespace are presumed to have opted in.

Uses of class template argument deduction must also follow the general rules on
[type deduction](#type-deduction-including-auto).

##### Why Use Class Template Argument Deduction

* Class template argument deduction can sometimes allow you to omit boilerplate from your code.

##### Why Not To Use Class Template Argument Deduction

* The implicit deduction guides that are generated from constructors may have undesirable behavior,
or be outright incorrect. This is particularly problematic for constructors written before class
template argument deduction was introduced in C++17 because the authors of those constructors had
no way of knowing about, much less fixing, any problems that their constructors would cause for
class template argument deduction. Furthermore, addeing explicit deduction guides to fix those
problems might break any existing code that relies on the implicit deduction guides.
* Class template argument deduction also suffers from many of the same drawbacks as `auto` because
they are both mechanisms for deducing all or part of a variable's type from its initializer.
Class template argument deduction does give the reader more information than `auto`, but it also
does not give the reader an obvious cue that information has been omitted.

##### How Class Template Argument Deduction Rules Are Enforced

Code review must be used to ensure that class template rules are followed.

#### Designated Intializers

Designated initializers are a syntax that allows initializing an aggregate (plain old struct) by
naming its fields explicitly:

```C++
struct Point {
    float x = 0.0;
    float y = 0.0;
    float z = 0.0;
};

Point p = {
    .x = 1.0,
    .y = 2.0,
    // z will be 0.0
};
```

The explicitly listed fields will be initialized as specified, and others will be initialized in
the same way as they would be in a traditional aggregrate initializatioi expression like
`Point{ 1.0, 2.0 };`.

Use designated initializers only in the form that is compatible with the C++20 standard: with 
initializers in the same order as the corresponding fields appear in the struct definition.

##### Why Use Designated Initializers

* Designated initializers can make for convenient and highly readable aggregate expressions,
especially for structs with less straightforward ordering of fields than the `Point` example above.

##### Why Not To Use Designated Initializers

* While designated initializers have long been part of the C standard and supported by C++ as an
extension, they were not supported by C++ prior to C++20.
* The rules in the C++ standard are stricter than in C and compiler extensions, requiring that the
designated initializers appear in the same order as the fields appear in the struct definition. So
in the example above, it is legal according to C++20 to initialize `x` and then `z`, but not `y` and
then `x`.

#### Lambda Expressions

Lambda expressions are a concise way of creating anonymous function objects.

Use lambda expression where appropriate, with formatting as described [[[below - #formatting
lambda expressions].

Prefer explicit captures if the lambda may escape the current scope. For example, instead of:

```C++
{
    Foo foo;
    ...
    executor->Schedule([&] { Forbicate(foo); })
    ...
}
// BAD! The fact that the lambda makes use of a reference to `foo` and
// possibly `this` (if `Frobnicate` is a member function) may not be
// apparent on a cursory inspection. If the lambda is invoked after
// the function returns, that would be bad, because both `foo`
// and the enclosing object could have been destroyed.
```

prefer to write:

```C++
{
    Foo foo;
    ...
    executor->Schedule([&foo] { Frobnicate(foo); })
    ...
}
// BETTER - The compile will fail if `Frobnicate` is a member
// function, and it's clearer that `foo` is dangerously captured
// by reference.
```

Use default capture by reference ([&]) only when the lifetime of the lambda is shorter than any
potential captures.

Use default capture by value ([=]) only as a means of boinding a few variables in a short lambda,
where the set of captured variables is obvious at a glance, and which does not result in capturing
`this` implicitly. That means that a lambda that appears in a non-static class member function
and refers to non-static class members in its body must capture `this` explicitly or via [&].
Prefer not to write long or complex lambdas with default capture by value.

Use captures only to actually capture variables from the enclosing scope. Do not use captures with
initializers to introduce new names, or to substantially change the meaning of an existing name.
Instead, declare a new variable in the conventional way and then capture it, or avoid the lambda
shorthand and define a function object explicitly.

See the sections on [type deduction](#type-deduction-including-auto) 
and [trailing return type syntax](#trailing-return-type-syntax)
for guidance on specifying the parameter and return types.

##### Why Use Lambda Expressions

* Lambdas are much more concise than other ways of defining function objects to be passed to STL
algorithms, which can be a readability improvement.
* Appropriate use of default captures can remove redundancy and highlight important exceptions from
the default.
* Lambdas, `std::function`, and `std::bind` can be used in combination as a general purpose
callback mechanism; they make it easy to write functions that thake bound functions as arguments.

##### Why Not To Use Lambda Expressions

* Variable capture in lambdas can be a source of dangling-pointer bugs, particularly if a lambda
escapes the current scope.
* Default captures by value can be misleading because they do not prevent dangling-pointer bugs.
Capturing a pointer by value does not cause a deep copy, so it often has the same lifetime issues
as capture by reference. This is especially confusing when capturing `this` by value since the use
of `this` is often implicit.
* Captures actually declare new variables whether or not the captures have initializers, but they
look like any other variable declaration syntax in C++.. In particular, there is no place for the
variable's type, or even an `auto` placeholder, although init captures can indicate it indirectly,
e.g. with a cast. This can make it difficult to even recognize them as declarations.
* Init captures inherently rely on [type deduction](#type-deduction-including-auto), and suffer
from many of the same drawbacks as `auto`, with the additional problem that the syntax does not
even cure the reader thatn deduction is taking place.
* It is possible for use of lambdas to get out of hand; very long nested anonymous functions can
make code harder to understand.

##### How Lambda Expression Rules Are Enforced

clang-tidy has the following checks related to lambdas:

* bugprone-lambda-function-name
* cppcoreguidelines-avoid-capturing-lambda-coroutines

This is a short list because the various lambda expression rules may vary considerably from one
coding standard or guideline to another. The only way to enforce lambda expression rules is through
the use of code reviews.

#### Template Metaprogramming

Template metaprogramming refers to a family of techniques that exploit the fact that the C++
templaate instantiation mechanism is Turing complete and can be used to perform arbitrary
compile-time computation in the type domain.

Template metaprogramming sometimes allows cleaner and easier-to-use interfaces than would be
possible without it, but it is also often a temptation to be overly clever. It is best used in a
small number of low level components where the extra maintenance burden is spread out over a large
number of uses.

Think twice before using template metaprogramming or other complicated template techniques; think
about whether the average member of your team will be able to understand your code well enough
to maintain it after you switch to another project, or whether a non-C++ programmer or someone
casually browsing the code base will be able to understand the error messages or trace the flow of
a function they want to call. If you are using recursive template instantiations, type lists or
metafunctions or expression templates, or relying of SFINAE or on the `sizeof` trick for detecting
function overload resolution, then there is a good chance you have gone too far.

If you use template metaprogramming, you should expect to put considerable effort into minimizing
and isolating the complexity. You should hide metaprogramming as an implementation detail
whenever possible, so that user-facing headers are readable, and you should make sure that tricky
code is especially well commented. You should carefully document how the code is used, and you
should say something About what the "generate" code looks like. Pay extra attention to the error
messages that the compile emits when users make mistakes. The error messages are port of your
user interface, and your code should be tweaked as necessary so that the error messages are
understandable and actionable from a user point of view.

##### Why Use Template Metaprogramming

* Template metaprogramming allows extremely flexible interfaces that are type safe and high
performance. Facilities like GoogleTest, `std::tuple`, `std::function`, and Boost.Spirit would
be impossible without it.

##### Why Not To Use Template Metaprogramming

* The techniques used in template metaprogramming are often obscure to anyone but language experts.
Code that uses templates in complicated ways is often unreadable, and is often hard to debug or
maintain.
* Template metaprogramming often leads to extremely poor compile time error messages: even if an
interface is simple, the complicated implementation details become visible when the user does
something wrong.
* Template metaprogramming interferes with large scale refactoring by making the job of refactoring
tools harder. First, the template code is expanded in multiple contexts, and it is hard to verify
that the transformation makes sense in all of them. Second, some refactoring tools work with the
AST that only represents the structure of the code after template expansion. It can be difficult
to automatically work back to the original source construct that needs to be rewritten.

##### How Template Metaprogramming Rules Are Enforced

Code review is the only method that call be used to enfore template metaprogramming rules.

#### Concepts and Constraints

The `concept` keyword is a new mechanism for defining requirements such as type traits or
interface specifications lfor a template parameter. The `requires` keyword provides mechanisms
for placingt anonymous constraints on templates and verifying that constraints are satisfied at
compile time. Concepts and constraints are often used together, but can also be used
independently.

Use concepts sparingly. In general, concepts and constraints should only be used in cases where
templates would have been used prior to C++20.

Avoid introducing new concepts in headers, unless the headers are marked as internal to the library.

Do not define concepts that are not enforced by the compiler.

Prefer constraints over template metaprogramming and avoid `template<Concept T> syntax`, instead
use the `requires(Concept<T>)` syntax.

Predefined concepts in te standard library should be preferred to type traits when equivalent ones
exist. Similarly, prefer modern constraint syntax (via `requires(Condition)`). Avoid legacy
template metaprogramming constructs such as `std::enable_if<Condition>` as well as
`tmeplate<Concept T>` syntax.

Do not manually reimplement any existing concepts or traits. For example, use
`requires(std::default_initializable<T>)` instead of `requires({ T v; })` or similar.

New `concept` declarations should be rare, and only defined internally within a library, such that
they are not exposed at API boundaries. More generally, do not use concepts or constraints in cases
where you would not use their legacy template equivalents in C++17.

Do not define concepts that dupliate the function body or impose requirements that would be
insignificant or obvious from reading the body of the code or the resulting error messages.
Instead, prefer to leave code as an ordinary template unless you can demonstrate that concepts
result in significant improvement for that particular case, such as in the resulting error
messages for a deeply nested or non-obvious requirement.

Concepts should be statically verifiable by the compiler. Do not use any concept whose primary
benefits would come from a semantic or otherwise unenforced constraint. Requirements that are
unenforced at compile tie should instead be imposed by other mechanisms such as comments,
assertions, or tests.

##### Why Use Concepts And Constraints

* Concepts allow the compiler to generate much better error messages when templates are involved,
which can reduce confusion and significalntly improve the development experience.
* Concepts can reduce the boilerplate necessary for defining and using compile-time constraints,
often increasing the clarity of the resulting code.
* Constraints provide some capabilities that are difficult to achieve with templates and SFINAE
techniques.

##### Why Not To Use Concepts And Constraints

* As with templates, concepts can make code significantly more complex and difficult to understand.
* Concept syntax can be confusing to readers, as concepts appear similar to class types at their
usage sites.
* Concepts, especially at API boundaries, increase code coupling, rigidity, and ossification.
* Concepts and constraints can replicate logic from a function body, resulting in code duplication
and increased maintenance costs.
* Concepts muddy the source of truth for their underlying contracts, as they are standalone named
entities that can be utilized in multiple locations, all of which evolve separately from each other.
* Concepts and constraints affect overload resolution in novel and non-obvious ways.
* As with SFINAE, constraints make it harder to refactor code at scale.

##### How Concepts And Constraints Rules Are Enforced

There is only one clang-tidy check related to concepts and constraints:

* modernize-use-constraints

This check flags and possibly replaces `std::enable_if` with C++20 `requires` clauses.

Concepts and constraints rules in general must be enforced using code review.

#### Boost

In general, Boost libraries should not be used. Requirements for Boost functionality will be
considered on a case-by-case basis.

#### Aliases

Do not put an alias in your public API just to save typing in the implementation; do so only if you
intend it to be used by the API's clients.

When defining a public alias, document the intent of the new name, including whether it is
guaranteed to always be the asme as the type it is currently aliased to, or whether a more
limited compatibility is intended. This lets the user know whether they can treat the types as
substitutable or whether more specific rules must be followed, and can help the implementation
retain some degree of freedome to change the alias.

Do not put namespace aliases in your public API. See also [Namespaces](#namespaces).

Local convenience aliases are fine in function definitions, `private` sections of classes,
explicitly marked internal namespaces, and in .cpp files.

##### Why Use Aliases

* Aliases can improve readability by simplifying a long or complicated name.
* Aliases can reduce duplication by naming in one place a type used repeatedly in an API, which
might make it easier to change the type later.

##### Why Not To Use Aliases

* When placed in a header where client code can refer to them, aliases increase the number of
entities in that header's API, increasing its complexity.
* Clients can easily rely on unintended details of public aliases, making changes difficult.
* It can be tempting to create a public alias that is only intended for use in the implementation
without considering its impact on the API, or on maintainability.
* Aliases can increase the risk of name collisions.
* Aliases can reduce readability by giving a familiar construct an unfamiliar name.
* Type aliases can create an unfamiliar API contract: it is unclear whether the alias is
guaranteed to be identical to the type it aliases, or only to be usable in specific narrow ways.

##### How Aliasing Rules Are Enforced

Coding reviews are the only way that aliasing rules can be enforced.

#### Switch Statements

If not conditional on an enumerated value, `switch` statements should alwasy have a `default` case.
In the case of an enumerated value, the compiler will warn if not all values have been handled. If
the `default` case should never execute, treat this as an error.

Fall-through from one case label to another must be annotated using the `[[fallthrough]];`
attribute. `[[fallthrough]];` should be placed at a point of execution where a fall-through to
the next label occurs. A common exception is consecutive case labels without intervening code, in
which case no annotation is needed.

##### How Switch Statement Rules Are Enforced

clang-tidy has the following switch statement related rule:

* bugprone-switch-missing-default-case

clang has the `-Wimplicit-fallthrough` command line argument. This test is also included with the
`-Weverything` command line argument.

gcc also has the `-Wimplicit-fallthrough` command line argument. This chceck is also performed with
the `-Wextra` command line argument.

MSVC checks and warns automatically on implict fallthrough in switch statements, so no command line
argument is provided.

### Inclusive Language

In all code, including naming and comments, use inclusive language and avoid terms that other
programmers might find disrespectful or offensive, such as "master", "slave", "blacklist",
"whitelist", and "redline", even if the terms also have an ostensibly neutral meaning.

Similarly, use gender-neutral language unless you are referring to a specific person and using
their pronouns. For example, use "they"/"them"/"their" for people of unspecified gender even when
singular, and "it"/"its" for software, computers, and other things that are not people.

#### Inclusive Language Enforcement

The only way that inclusive language can be enforced is through code review.

## Naming

The most important consistency rules are those that govern naming. The style of a name immediately
informs us what sort of thing the named entity is: a type, a variable, a function, a constant, a
macro, etc. without requiring us to search for the declaration of that entity. The pattern-matching
engine in our brains relies a great deal on these naming rules.

### Why These Naming Rules

Naming rules are fairly arbitrary, but consistency is more important than individual preferences
in this area, so regardless of whether you find them sensible or not, the rules are the rules.

Since this is a project that I started and am leading, I get to make the naming rules.

### General Naming Rules

Optimize for readability using names that would be clear even to people on a different project.

Use names that describe the pupose or intent of the object. Do not worry about saving horizontal
space as it is far more important to make your code immediately understandable by a new reader.
Minimize the use of abbreviations that would likely be unknown to someone outside your project,
especially acronyms and initialisms. Do not abbreviate by deleting letters within a word. As a
rule of thumb, an abbreviation is OK if it is listed in Wikipedia. Generally speaking,
descriptiveness should be proportional to the name's scope of visibility. For example, `n` may
be fine within a 5 line function, but within the scope of a class, it is likely too vague.

Note that certain universally-known abbreviations are OK, such as `i` for an iteration variable and
`T` for a template parameter.

For the purposes of the naming rules below, a "word" is anything that you would write in English
without internal spaces. This included abbreviations such as acronyms and initialisms. For names
written in mixed case, known as "camelCase" or "PascalCase", in which the first letter of each
word is capitalized, prefer to capitalize abbreviations as single words. For example, `StartRpc()`
rather than `StartRPC()`.

Template parameters should follow the naming style for their category: type template parameters
should follow the rules for [type names](#type-names), and non-type template parameters should
follow the rules for [variable names](#variable-names).

#### File Names

The rules in this section apply to C++ specific files; there are rules for other file types
that are based on convention for those types or languages.

Filenames should be all lowercase and include underscores (`_`) as word separators. The general
rule is to name source and header files after the class name that the file contains. For example,
the files containing the declaration and definition for a class named `JSDRMainframe` should be in
the files `jsdr_mainframe.h` and `jsdr_mainframe.cpp`. For files that do not contain class
declarations and definitions, make the filenames very specific. For example, use 
`http_server_logs.h` rather than `logs.h`.

Do not use filenames that already exist in general library locations such as `/usr/include` and
`C:\Windows\`, such as `db.h`.

As indicated in the paragraphs above, C++ source files should end in `.cpp` and header files should
end in `.h`. Files that rely on being textually included at specific points should end in `.inc`.
See also the section on [self-contained headers](#self-contained-headers).

##### How File Names Are Enforced

Filenames are enforced during code review.

#### Type Names

Type names use PascalCase with no underscores.

The names of all types: classes, structs, type aliases, enums, and type template parameters, use
the same naming convention. For example:

```C++
// classes and structs
class UrlTable { ...
class UrlTableTester { ...
struct UrlTableProperties { ...

// typedefs
typedef hash_map<UrlTableProperties *, std::string> PropertiesMap;

// using aliases
using PropertiesMap = hash_map<UrlTableProperties *, std::string>;

// enums
enum class UrlTableError { ...
```

#### Variable Names

The names of variables, including function parameters, and data members are `camelCase`. Private
data members of classes have a single starting underscore.

##### Common Variable Names

For example:

```C++
std::string tableName;
```

##### Class Member Names

Data members of classes, both static and non-static, are named like nonmember variables,, but with
a leading underscore.

```C++
class TableInfo {
...
private:
std::string _tableName;
static Pool<UsrTableProperties>* _pool;
};
```

##### Struct Data Members

Data members of structs, both static and non-static, are named like ordinary nonmember values. They
do not have the prepended underscores that data members in classes do.

```C++
struct UrlTableProperties {
    std::string name;
    int numEntries;
    static Pool<UrlTableProperties>* pool;
};
```

See [Structs vs. Classes](#structs-vs-classes) for a discussoin on when to use a struct versus a
class.

#### Constant Names

Variables declared `constexpr` or `const`, and whose value ios fixed for the duration of the
program, are named with a leading "k" followed by PascalCase. Underscores can be used as separators
in the rare cases where capitalization cannot be be used for separation. For example:

```C
const int kDaysInAWeek = 7;
const int kAndroid8_0_0 = 24;   // Android 8.0.0
```

All such variables with static storage duration (i.e. statics and globals, see
[[[Storage Duration](#storage-duration)]] for details) should be named this way, including those
in templates where different instantiations of the template may hae different values. This
convention is optional for variables of other storage classes (e.g. automatic variables; otherwise
the usual variable naming rules apply. For example:

```C++
void ComputeFoo(std::string_view suffix) {
    // Either of these is acceptable
    const std::string_view kPrefix = "prefix";
    const std::string_view prefix = "prefix:";
    ...
}
```

```C++
void ComputeFoo(std::string_view suffix) {
    // Bad - different invocations of ComputeFoo give kCombined different values.
    const std::string kCombined = std::string_)view(kPrefix, suffix);
}
```

#### Function Names

Regular functions should be named using PascalCase; accessors and mutators may be named like
variables. For example:

```C++
AddTableEntry();
DeleteUrl();
OpenFileOrDie();
```

The same naming rule applies to class- and namespace-scope constants that are exposed as part of an
API and that are intended to look like functions, because the fact that they are objects rather
than functions is an unimportant implementation detail.

Accessors and mutators (get and set functions) may be named like variables. These often correspond
to actual data member variables, but this is not important. For example: `int count()` and
`void setCount(int count)`.

#### Namespace Names

Namespace names are all lower-case, with words separated by underscores. Top-level namespaces are
based on the project name. Avoid collisions between nested namespaces and well-known top-level
namespaces.

The name of a top-level namespace should usually be the name of the project whose code is
contained in that namespace. The code in that namespace should usually be in a directory whose
basename matches the namespace or in subdirectories thereof.

Keep in mind that the [rule against abbreviated names](#general-naming-rules) applies to namespaces
just as much as variable names. Code inside the namespace seldom needs to mention the namespace
name, so there is usually no particular need for abbreviation anyway.

Avoid nested namespaces that match well-known top-level namespaces. Collisions between namespace
names can lead to surprising build breaks because of name lookup rules. In particular, do not
create any nested `std` namespaces. Prefer  unique project identifiers (`websearch::index`,
`websearch::index_util`) over collision-prone names like `websearch::util`. Also avoid overly
deep nesting namespaces.

For `internal` namespaces, be wary of other code being added to the same `internal` namespace
causing a collision (internal helpers within a team tend to be related and may lead to collisions).
In such a situation, using the filename to make a unique internal name is helpful
(`websearch::index::frobber_internal` for use in `frobber.h`).

#### Enumerator Names

Enumerators for both scoped and unscoped `enum`s should be named like [constants](#constant-names),
not like [macros](#macros). That is, use `kEnumName`, not `ENUM_NAME`.

```C++
enum class UrlTableError {
    kOK = 0,
    kOutOfMemory,
    kMalformedInput,
};
```

#### Macro Names

You are not really going to [define a macro](#preprocessor-macros) are you? If you do, they are
like this: `MY_MACRO_THAT_SCARES_SMALL_CHILDREN_AND_ADULTS_ALIKE`.

Please see the [description of macros](#preprocessor-macros); in general, macros should *not* be
used. However, if they are absolutely needed, they should be names with all capitals and
underscores, and with a project-specify prefix.

```C++
#define JSDR_ROUND(x) ...
```

#### Exceptions To Naming Rules

If you are naming something that is analogous to an existing C or C++ entity, then you can follow
the existing naming convention scheme.

    `bigopen()`
        function name, follows form of `open`
    `uint`
        typedef
    `bigpos`
        `struct` or `class`, follows form of `pos`
    'sparse_hash-map`
        STL-like entity; follows STL naming conventions
    `LONGLONG_MAX`
        a constant, as in `INT_MAX`

### Comments

Comments are absolutely vital to keeping code readable. The following rules describe what you should
comment and where. But remember: while comments are very important, the best code is
self-documenting. Giving sensible names to types and variables is much better than using
obscure names that you must explain through comments.

When writing your comments, write for your audience: the next contributor who will need to
understand your code. Be generouse - the next one may be you!

#### Comment Style

Use either `//` or `/*  */` syntax; however, `//` is *much* more common. Be consistent with how
you comment and what style you use where.

When documenting

#### File Comments

Start each file with license boilerplate.

If a source file (such as a `.h` file) declares multiple user-facing abstractions (common functions,
related classes, etc.), include a comment describing the collection of these abstractions.
Include enough detail for future authors to know what does not fit there. However, the detailed
documentation about individual abstractions belongs with those abstractions, not at a file level.

For instance, if you write a file comment for `frobber.h`, you do not need to include a file
comment in `frobber.cpp` or `frobber_test.cpp`. On the other hand, if you write a collection of
classes in `registered_objects.cpp` that has no associated header file, you must include a file
comment in `registered_objects.cpp`.

##### Legal Notice And Author Line

Every file should contain license boilerplate. Chose the appropriaate boilerplate for the license
used by the project (for example: Apache 2.0, BSD, LGPL 2.0, GPL 3.0).

If you make significant changes to a file with an author line, consider deleting the author line.
New files should usually not contain a copyright notice or author line.

#### Struct And Class Comments

Every non-obvious class or struct declaration should have an accompanying comment that describes
what it is for and how it should be used.

```C++
// Iterates over the contents of a GargantuanTable.
// Example:
//    std::unique_ptr<GargantuanTableIterator> iter = table->NewIterator();
//    for (iter->Seek("foo"); !iter->done(); iter->Next()) {
//      process(iter->key(), iter->value());
//    }
class GargantuanTableIterator {
  ...
};
```

##### Class Comments

The class comment should provide the reader with enough information to know how and when to use
the class, as well as additional considerations necessary to correctly use the class. Document
the synchronization assumptions the class makes, if any. If an instance of the class can be
accessed by multiple threads, take extra care to document the rules and invariants surrounding
multithreaded use.

The class comment oftenn is a good place for a small sample code snippet demonstrating a simple and
focused usage of the class.

When sufficiently separated (e.g. `.h` and `.cpp` files), comments describing the use of the class
should go together with its interface definition; comments about the class operation and
implementation should accompany the implementation of the class's methods.

#### Function Comments

Declaration comments describe use of the function when it is not obvious; comments at the definition
of a function describe operation.

##### Function Declarations

Almost every function declaration should have comments immediately preceding it that describe
what the function does and how to use it. Private methods and functions declared in `.cpp` files
are not exempt. Function comments should be written with an implied subject of *This function* and
should start with a verb phrase; for example: "Opens the file", rather than "Open the file". In
general, these comments do not describe how the function performs the task. Instead, that should
be left to comments in the function definition.

Types of things to mention in comments at the function declaration:

* What the inputs and outputs are.
* For class member functions: whether the object remembers reference or pointer arguments beyond
the duration of the method call. This is quite common for pointer/reference arguments to
constructors.
* For each pointer argument, whether it is allowed to be null and what happens if it is.
* For each output or input/output argument, what happens to any state that argument is in.
For example, is the state appended or overwritten?

Here is an example:

```C++
// Returns an iterator for this table, positioned at the first entry
// lexically greater than or equal to `start_word`. If there is no
// such entry, returns a null pointer. The client must not use the
// iterator after the underlying GargantuanTable has been destroyed.
//
// This method is equivalent to:
//    std::unique_ptr<Iterator> iter = table->NewIterator();
//    iter->Seek(start_word);
//    return iter;
std::unique_ptr<Iterator> GetIterator(std::string_view start_word) const;
```

Do not be unnecessarily verbose or state the completely obvious.

When documenting function overrides, focus on the specifics of the override itself, rather than
repeating the comment from the overridden function. In many cases, the override needs no
additional documentation and thus no comment is required. See
[Using Doxygen For API Documentation](#using-doxygen-for-api-documentation) for an exception to this
rule.

When commenting constructors and destructors, remember that the person reading your code knows what
constructors and destructors are for, so comments that just say something like "destroys this
object" are not useful. Document what constructors do with their arguments (for example, if they
take ownership of pointers), and what cleanup the destructor does. If this is trivial, just skip
the comment. It is quite common for destructors not to have a header comment.

##### Function Definitions

If there is anything tricky about how a function does its job, the function definition should have
an explanatory comment. For example, in the definition comment you might describe any coding tricks
you use, give an overview of the steps you go through, or explain why you chose to implement the
function the way you did rather than using a viable alternative. For instance, you might mention
why it must acquire a lock for the first half of the function but why it is not needed for the
second half.

#### Variable Comments

In general, the actual name of the variable should be desctiptive enought to give a good idea of
what the variable is used for. In certain cases, more comments are required.

##### Class Data Members

The purpose fo each class data member must be clear. If there are any variants (special 
values, relationships between members, lifetime requirements) not clearly expressed by the type
and name, they must be commented. However, if the type and name suffice (e.g. `int _numEvents`;),
no comment is needed.

In particular, add comments to describe the existence and meaning of sentinel values, such as
`nullptr` or `-1`, when they are not obvious. For example:

```C++
private:
    // Used to bounds-check table accesses. -1 means
    // that we don't yet know how many entries the table has.
    int _numTotalEntries;
```

##### Global Variables

All global variables should have a comment describing what they are used for, and, if unclear, why
they need to be global. For example:

```C++
// The total number of test cases that we can run in this regression test.
const int kNumTestCases = 6;
```

#### Implementation Comments

In your implementation you should have comments in tricky, non-obvious, interesting, or important
parts of your code.

##### Explanatory Comments

Tricky or complicated code blocks should have comments before them.

#### Function Argument Comments

When the meaning of a function argument is non-obvious, consider one of the following  remedies:

* If the argument is a literal constant, and the same constant is used in multiple function calls
in a way that tacitly assumes they are the same, you should use a named constant to make that
constraint explicit, and to guarantee that it holds.
* Consider changing the function signature to replace a `bool` argument with an `enum` argument.
This will make the argument values self-describing.
* For functions that have several configuration options, consider defining a single class or
struct to hold all the options, and pass an instance of that. This approach has several advantages:
  * Options are referenced by name at the call site, which clarifies their meaning.
  * It reduces the function argument count, which makes the function calls easier to read and write.
  * You don't have to change call sites when you add another option.
* Replace large or complex nested expressions with named variables.
* As a last resort, use comments to clarify argument meanings at the call site.

Consider the following example:

```C++
// What are these arguments?
const DecimalNumber product = CalculateProduct(values, 7, false, nullptr);
```

versus:

```C++
ProductOptions options;
options.setPrecisionDecimals(7);
options.setUseCache(ProductOptons::kDontUseCache);
const DecimalNumber product = 
    CalculateProduct(values, options, /*completionCallback=*/nullptr);
```

#### Don'ts

Do not state the obvious. In particualr, do not literally describe what the code does, unless the
behavior is non-obvious to a reader who understands C++ well. Instead, provide higher level
comments that describe why the code does what it does, or make the code self-describing.

Compare this:

```C++
// Find the element in the vector.  <-- Bad: obvious!
if (std::find(v.begin(), v.end(), element) != v.end()) {
    Process(element);
}
```

to this:

```C++
// Process "element" unless it  was already processed.
f (std::find(v.begin(), v.end(), element) != v.end()) {
    Process(element);
}
```

Self-describing code does not need a comment. The comment from the example above would be obvious:

```C++
if (!IsAlreadyProcessed(element)) {
    Process(element);
}
```

#### Punctuation, Spelling, And Grammar

Comments should be written in American English.

Comments should be readable as narative text, with proper capitalization and punctuation. In many
cases, complete sentences are more readable than sentence fragments. Shorter comments, such as
commnents at the end of a line of code, can sometimes be less formal, but you should be consistent
with your style.

Although it can be frustrating to have a code reviewer point out that you are using a comma when
you should be using a semicolon, it is very important that source code maintain a high level of
clarity and readability. Proper punctuation, spelling, and grammar help with that goal.

##### Why Pay Attention To Punctuation, Spelling, and Grammar

* It is easier to read well-written comments than badly written ones.
* While it can be frustrating and non-intuitive for people who are from cultures that use British 
English or its derivatives to have to write American English, American English is the defacto
standard for writing code and documentation.

#### TODO Comments

Use TODO comments for code that is temporary, a short-term solution, or good-enough but not
perfect.

TODOs should include the string `TODO` in all capitals, follwed by the bug ID, name, or other
identifier of the person or issue with the best context about the problem referenced by the TODO.

```C++
// TODO: bug 12345678 - Remove this after the 2047q4 compatibility window expires.
// TODO: example.com/my-design-doc - Manually fix up this code the next time it's touched.
// TODO(bug 12345678): Update this list after the Foo service is turned down.
// TODO(John): Use a "\*" here for concatenation operator.
```

If your TODO is of the form "At a future date do something", make sure that you either include a
very specific date ("Fix by November 2025") or a very specific event ("Remove this code when all
clients can handle XML responses.").

##### Why Use TODO Comments

* TODO comments provide information about what code needs modification, when, and why.
* TODO comments are recognized by a many IDEs. These IDEs maintain a list of the TODOs, so that
you know what changes still need to be addressed.

### Using Doxygen For API Documentation

Documentation created using Doxygen is intended for both the users and maintainers of a library.
As such, Doxygen comments
should be used only on the external interfaces of the library APIs. Documentation of internal
interfaces of a library are intended only for the library maintainers, and therefore, standard 
comments should be used there.

Doxygen comment blocks use a special format to separate
them from normal comments. There are several ways to specify comment blocks. The preferred way is
listed below.

Doxygen also has a large number of commands that help it to divide the documentation into various
sections. For example, `@class` can be used to document information about a class at the class
level, such as what the class does, and `@param` is used to describe each parameter for a function
or class member. There are too large a number of Doxygen commands to list in this document. An
[alphabetical list of all commands](https://www.doxygen.nl/manual/commands.html) is available in 
the Doxygen documentation.

The following rules apply to using Doxygen to document a project's APIs:

* Doxygen comments should document the external interfaces of an API only, and not any internal
interfaces. Normal comments should be used to document any internal interfaces and implementation
details.
* Doxygen comments are typically placed only in `.h` files of a library because they describe
the external interfaces of the library. While they can also be placed in `.cpp` files, the
external interfaces of a library are not normally described in those files.
* Doxygen comments may be placed before or after the item they are documenting. The preferred
location is before, but occasionally the comments are required to be placed after the item.
* For comments placed before the item being described, begin each Doxygen comment line with `///`.
Alternatively, though not preferred, is to use a comment block as follows:

```
/**
 * Comment line 1
 * Comment line 2
 */
```

* Comments placed after the item being described, either on the same line, or on subsequent lines,
should begin with `///<`. For example:

```C++
int var;    ///< Detailed description of var
///< Additional detailed description
```

* Commands may be specified by prefacing them with either `\` or `@`. For example, `\param` and
`@code`. `@` is the preferred prefix in this project.
* As much information as possible should be included in the Doxygen documentation. For example:
  * Document what each class does.
  * Document what every class method, including constructors and the destructor do:
    *If what the constructors and destructor do is obvious, then just document them as `Constructor`
    or `Destructor`. However, if they do anything that is not immediately obvious, then describe
    what processing each does.
  * For every function or method:
    * Document what the function or method does.
    * Document every parameter:
      * For input parameters, is there a range of acceptable values? What happens if an
      unacceptable or unexpected value is input?
      * For pointer parameters, is ownership transferred from the caller to the callee?
      * For pointer parameters, can the value be `nullptr`? If so, what happens in that case?
      * For output or input/output parameters, what output is expected?
    * Document the return value from the function. What values can the return value contain?
    * What exceptions, if any, does the function throw, and why?
  * Provide a relatively simple sample of how the class and some of is methods is used.
* Doxygen also accepts markdown formatted files. Use them to provide information on topics related
to classes and functions in the API. For example:
  * A mainpage that provides an overview of the classes in the API:
    * General classes
    * Classes related to specific topics.
    * Features. For example:

        ## Features

        * Free and open source (MIT license).
        * A collection of native C++ classes to implement functionality that is not including in the STL.
        * Fully integrated with the STL.
        * Written in efficient, modern C++20 with the RAII programming idiom.
        * Highly portable and available on Windows, MacOS and Linux. They should be portable to
        iOS and Android as well, although they have not been tested on the latter two operating
        systems.

    * Library architecture. For example:
      * Does it contain templated and non-templated classes?
      * Is it functionally divided into multiple parts. If so, a quick description of each part.
    * Getting Started guide. This may be included in the top-level page, or simply contain the link
    to a separate document.
  * How the classes can be used to accomplish a task or set of tasks.
  * Why you might use one or more of the classes in the API instead of classes or functions in the
  STL. Possibly include code to illustrate how to accomplish a task using both libraries.

### Formatting

Coding style and formatting are fairly arbitrary but having a specific style and format makes a
project easier for all developers to follow.

As long as you do not format your code in weird ways, you do not need to worry about how *you*
format your code because, as part of the build process, clang-format is run to reformat the any
new code to match the style and formattng specified in the subsections below.

#### Why Consistent Formatting is Important

* A project is much easier to follow and understand if everyone uses the same style.

#### How Formatting Is Enforced

* clang-format is run on all changed source files to modifies this project's source code according
to the style options provided in the
`.clang-format` configuration file located in this project's top-level directory. The style
options specified in the configuration file enforce the style guidelines specified below.
* Additional tools, as noted in the appropriate subsections, may be used to help you follow the
appropriate style.

#### Line Length

Each line of text in the code should be at most 100 characters long. 80-character line lengths are
a hold-over from the days of coding on punched cards. In the early
days on coding using text-based computer monitors, the monitors could only hold 80 characters. In
the years since the introduction of GUI-based monitors, this limitation no longer is forced upon
you by your hardware. In many projects, line lengths of up to 132 characters have been specified.

100 characters is a compromise between these two line lengths.

##### Why 100-Character Long Lines

* Wider lines can make code more readable. The 80-character line width is a throwback to 1960's
mainframes and early text-based monitors.
* It may not be feasible to split comment lines without harming readability, ease of cut and paste,
or auto-linking.
* 100-character line lengths are a compromise between the 80-character line lengths of old, and
the longer 120- or even 132-character line lengths used in some other projects.


##### Why Use 80-Character Line Lengths

Some developers prefer 80 character line lengths because they prefer to have two or more code
windows open side-by-side and therefor do not have room to widen their windows in that case.
Poeple set up their work enviroment assuming a particular maximum window width, and 80 columns
has been the traditional standard.

##### When Can The 100-Character Line Length Be Exceeded

A line may exceed 100-characters if it is:

* A comment line which is not feasible to split without harming readability, ease of cut and paste,
or auto-linking. For example, if a line contains an example command or a URL longer than
100 characters.
* A string literal that cannot easily be wrapped at 100 columns. This may be because it contains
URIs or other semantically-critical pieces, or because the literal contains embedded language,
or a multiline literal whose newlines are significant like help messsages. In these cases,
breaking up the literal would reduce reabililty, searchability, ability to click links, and so
forth. Except for teest code, such literals should apprear at namespace scope near the top of the
file. If clang-format does not recognize the unsplittable content, 
[disable the tool](https://clang.llvm.org/docs/ClangFormatStyleOptions.html#disabling-formatting-on-a-piece-of-code)
around the content as necessary.t.
* An include statement.
* An [include guard](#22-include-guards).
* A using-declaration.

##### Editor Rulers

Many text editors and IDEs have settable properties that will display a ruler at one or more
columns. For example, the following contents in the `.vscode/settings.json` file sets a visible
ruler at column 100:

```json
{
    "editor.rulers": [
    {
        "column": 100
    }]
}
```

Most text editors and IDEs have equivalent settings.

`editorconfig` is a file format for defining editor properties and a collection of text editor
plugins or extensions for various text editors and IDEs. Unfortunately, it does not currently
support rulers. There is an issue ([#89](https://github.com/editorconfig/editorconfig/issues/89))
requesting that a rulers property be added to editorconfig. That request was orignially made in 2013
and was closed in 2016 when the `max_line_length` property was added. This is not the same as
`rulers`, and comments on issue #89 have been added yearly since 2018 again requesting a rulers
property. `max_line_length` forcibly wraps lines at the specified column, which is not the same and
does not allow the maximum length to be overridden on a case by case basis.

#### Non-ASCII Characters

Non-ASCII characters should be rarely needed. UTF-8 characters should be used in those rare cases.

You should not hard-code user-facing text in source, even in English, so use of non-ASCII
characters should be rare. In such cases, you
should use UTF-8, sice that is an encoding understood by most tools able to handle more than just
ASCII. See also [Inteernationalization](#internationalization).

Hex encoding is also OK, and encouraged when it enhances readability. For example, "xEF\]xBB\xBF",
or, even more simply, "\uFEFF", is the Unicode zero-width no-break space character, which would
be invisible if included in source as straight UTF-8.

When possible, avoid the `u8` prefix. It has significantly different semantics between C++17,
C++20, and C++23. In C++17, arrays are of type `char`, and in C++20, they are of type `char8_t`.
C++23 introduces several changes that result in formerly acceptable sequences as being ill-formed.
See, for example: 
[C++23: Growing unicode support](https://www.sandordargo.com/blog/2023/11/29/cpp23-unicode-support).

You should not use `char16_t` and `char32_t` character types since they are non-UTF-8 text.
Similarly, you should not use `wchar_t` unless you are writing code that interacts with the
Windows API, which uses `wchar_t` extensively.

##### Why Use Non-ASCII Characters

* The main reason for using non-ASCII is to support internationalization (i18n), and a number of
different locales.
* If your code parses data files from foreign sources, it may be appropriate
to hard-code non-ASCII strings used in those data files as delimiters.
* Unittest code, which does not need to be localized, might contain non-ASCII strings.

##### Why Not To Use Non-ASCII Characters

* There are a number of different unicode character sets: UTF-8, UTF-16, and UTF-32.
The encodings are different among all three.
* There are a number of STD functions that will convert from one encoding to another, including
from ASCII, since ASCII is a subset of UTF-8. But using them can somewhat obfuscate your code.
Therefore, if you must convert between various character set encodings, you may wish to perform
these conversions as near to the sources/syncs for the strings as possible.

##### How Use Of Non-ASCII Characters Is Enforced

The only way to enforce the rules related to non-ASCII characters is through the use of code
reviews.

#### Internationalization

Internationalizing user-facing messages and files can be a fairly daunting task. The steps needed
to support internationalizing text involves:

* Identifying the text that is to be internationized.
* Create a translation file between the original text, presumably in English, and the language
that the text is to be translated to. There are two things to note here:
  * Translation files may be based on language only, or on locale. For example, there could be a
  translation file from English to German, or from English to Austrian German, from English to Swiss
  German, and so forth. It is also possible to have a translation file from English to German and
  another translation file for English to Swiss German. The Swiss German translation would be used
  if the user's locale is `de-CH`, but the general German translation file would be used for any
  other locale that begins with `de`.
  * Some languages are written right-to-left, not left-to-right.
  * Translation phrases can affect the order in which variables are inserted into the text.
  * Number, time, and money formatting differ from one language to another.
* Keeping translation files up to date as new character strings are added to a project.

If internationalization must be supported, then the `Boost.Locale` library should be used to
provide message localization. `Boost.Locale` uses the GNU `gettext` localization model. The
advantage of using this model is that a number of open-source C++ libraries that may be used in a
project either use Boost.Locale, or another library based on GNU `gettext`, and so combining
translation files from the libraries and this project is relatively easier that attempting to
maintain separate translation files for each library.

##### Why Internationalize Text

* Not all users are fluent in English, so presenting user interfaces and messages in their languages
helps them to understand the program.

##### Why Not To Internationalize Text

* Text, such as error messages, that are intended for users of a library do not need to be
translated because the standard language for code development is English.
* There are hundreds if not thousands of languages and locales to support. It is not possible to
support them all. For any language/locale that is not supported, or for which no translated text is
provided, the text is displayed in the original language (presumably English).
* In open-source and other non-commercial projects, the developers must rely on outside help to
provide the translation files. That is fine for the original project version, but outside help may
be unavailable or unwilling to update the translation files for new versions of the project.
* Different programming languages and operating systems support different non-ASCII character sets
and encodings (e.g. code pages, UTF-8, UTF-16, UTF-32). Code must support a single encoding for
translations to function correctly.
* Different operating systems and compilers support translations differently. For example, the
default for Windows and MSVC is to place the translation strings in a project's resource files,
while *nix systems tend to use separate translation (`.po`) files.
* Different libraries may use their own translation methods. For example, wxWidgets provides the
`wxLocale` class for message translation, which closely follows the GNU `gettext` package. It may
be difficult to combine the translation files for multiple libraries into a single translation file.

##### How Internationalization Support Is Enforced

Code review can be used to ensure that:
* Internationalization is not used when not needed.
* That `Boost.Locale` is used when internationalization is needed.

#### Spaces vs. Tabs

Use only spaces, and insert 3 spaces at a time.

Why three spaces? Well, 4 is typical, but some people use 2 spaces, so I chose 3 as a comprimize.

##### How Spaces Versus Tabs is Enforced

The .clang-format configuration file for this project contains the following three settings:

```
Indent Width: 3
TabWidth: 3
UseTab: never
```

#### Function Declarations and Definitions

There are quite a few rules related to function declarations and definitions:

* Choose good parameter names.
* A parameter name may be omitted only if the parameter is not used in the function's definition.
* An unused parameter that might not be obvious should comment out the variable name in the
function definition.
* If you cannot fit the `auto` keyword and the function name on a single line, break between them.
* If you break after the return type of a function declaration or definition, do not indent.
* The open parenthesis is always on the same line as the function name.
* There is never a space between the function name and the open parenthesis.
* There is never a space between the parentheses and the parameters.
* The open curly brace is always on the end of the last line of the function declaration, not the
start of the next line.
* The close curly brace is either on the same line as the open curly brace or on the last line
by itself.
* There should be a space between the close parenthesis and the open curly brace.
* All parameters should be placed on the same line as the function name in the declaration or
definition. If not all parameters fit on a single line, then they should be placed with one
parameter per line.
* All parameters should be aligned if possible.
* Default indentation is 3 spaces.
* Wrapped parameters have a 3 space indent.
* The [trailing return type](#trailing-return-type-syntax) should be placed on the same line as the
function declaration or definition's closing parenthesis if possible. If not possible, then it 
should be placed on the immediately following line.
* Consecutive

##### How Function Declaration and Definition Rules Are Enforced

clang-format has a number of configuration rules that will format function declarations and
definitions according to the coding rules listed above:

* AlignAfterOpenBracket: Align
* AllowAllParametersOfDeclarationOnNextLine: false
* AllowShortFunctionsOnASingleLine: All
* AlwaysBreakAfterReturnType: None
* IndentWidth: 3
* IndentWrappedFunctionNames: true

Code review should also be used to ensure that any coding rules, not handled by the clang-format
rules, are followed.

#### Lambda Expressions

Format parameters and bodies as for any other [function](#function-declarations-and-definitions),
and capture lists like other comma-separated lists.

For by-reference captures, do not leave a space between the ampersand (`&`) and the variable name.

Short lambdas may be written inline as function arguments.

##### Enforcing Lambda Expression Rules

The clang-format rules for
[function declarations and functions](#function-declarations-and-definitions) apply automatically
to lambda expressions. Code review should be used as well.

#### Floating-point Literals

Floating point literals should always have a radix point, with digits on both sides. For example:

```C++
float f = 1.0f;
long double ld = -0.5L;
double d = 1248.0e6;
```

##### Why These Rules For Floating-Point Literals

* Readability is improved if all floating-point literals take this familiar form, as this helps
ensure that they are not mistaken for integer literals, and the `E/e` of the exponential notation
is not mistaken for a hexadecimal digit.