# How to Build jsdr

The jsdr project is mainly written in C++ and uses CMake to build the project.
This means that the project can be built either from the command line, or from
just about any IDE that supports C++.

Since I use Visual Studio Code, I have
provided a `settings.json` file in the `.vscode` directory to match my settings.
If you use a different IDE, you should create a settings file with the
equivalent settings for your IDE. I will gladly accept pull requests for
settings files for other IDEs.

## MacOS

You will need a number of tools to build the jsdr project. The following
terminal commands show how to install them. It assumes that you have already
installed [Homebrew](https://brew.sh/) on your system.

```zsh
xcode-select --install
brew update
brew upgrade
brew install git cmake ninja
brew install wxWidgets
brew install clang-format
brew install llvm
# if Apple silicon enter the following line:
ln -s "$(brew --prefix llvm)/bin/clang-tidy" "/opt/homeberew/bin/clang-tidy"
# else Intel enter the following:
ln -s "$(brew --prefix llvm)/bin/clang-tidy" "/usr/local/bin/clang-tidy"

```

Enter the following from the command line to download the jsdr project files:

```zsh
cd <your-root-project-directory>
git clone https://github.com/jimorc/jsdr.git
```

For example, my root project directory is `~/Projects`, not `~/Projects/jsdr`.

### Building jsdr from the Command Line

To build jsdr from the command line, enter the following:

```zsh
cd jsdr
mkdir build
cd build
cmake .. 
ninja
```

The `cmake` command will generate a ninja build file. The `ninja` command
will perform the build. By default, `ninja` runs build commands in parallel.
You can specify the number of build commands to run in parallel using the
`-j` option. For example: `ninja -j4` will run 4 build commands in parallel.
Because `ninja` runs build commands in parallel by default, you don't gain
anything by using the `-j` command line option unless you want to force running
a single threaded build.

### Building jsdr using an IDE

If your IDE automatically runs `cmake`, you can simply open the jSDR project
directory in your IDE. For example, Visual Studio Code does this.
If you must generate IDE-specific build files, then at
the command line, enter:

```zsh
cd jsdr
mkdir build
cd build
cmake .. -G <build-tool-generator>
```

where `<build-tool-generator>` specifies the generator type. The documentation
for your IDE should specify what generator you need to use. Alternatively,
you can view the list of supported generators at
[cmake-generators](https://cmake.org/cmake/help/latest/manual/cmake-generators.7.html#manual:cmake-generators(7)).
