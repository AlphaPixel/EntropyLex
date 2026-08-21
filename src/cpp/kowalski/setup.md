# Setup

## Build & Test

```bash
cmake -S . -B ./build -DCMAKE_EXPORT_COMPILE_COMMANDS=ON
cmake --build ./build
ctest --test-dir ./build --output-on-failure
```

## Lint & Format

Requires `compile_commands.json` (generated above).

```bash
# lint
find . -name '*.cc' -not -path './build/*' | xargs clang-tidy -p build/

# format check
find . \( -name '*.cc' -o -name '*.h' \) -not -path './build/*' | xargs clang-format --dry-run --Werror

# format fix
find . \( -name '*.cc' -o -name '*.h' \) -not -path './build/*' | xargs clang-format -i
```

## Pre-commit Hook

```bash
# requires python installed
pip install pre-commit
pre-commit install
```

Runs clang-tidy/clang-format on staged files at commit time. Build `./build` at least once first so `compile_commands.json` exists.

Run against all files manually:
```bash
pre-commit run --all-files
```

