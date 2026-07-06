# mangolang Language Specification

> Status: Sketch. Examples show intended feel, not finalized syntax.

## Source Files

UTF-8. Extension `.um`. No semicolons. Tabs for indentation.

## Naming

`camelCase` for values and functions. `PascalCase` for types.

## Examples

### Functions

```um
fn add(Int x, Int y) -> Int =
	x + y

fn greet(String name) -> String =
	"Hello, " + name

fn compute(Int x) -> Int = {
	let doubled = x * 2
	doubled + 1
}
```

### Types

```um
type Point = { Float64 x, Float64 y }

type Shape =
	| Circle { Float64 radius }
	| Rectangle { Float64 width, Float64 height }

type Option<T> =
	| Some T
	| None

type Result<T, E> =
	| Ok T
	| Err E
```

### Bindings

```um
let x = 42
let mut counter = 0
let String name = "mango"
```

### Pattern Matching

```um
match shape {
	Circle { radius } -> radius * radius * 3.14159
	Rectangle { width, height } -> width * height
}
```

### Pipeline

```um
result = data
	|> filter(isValid)
	|> map(transform)
	|> reduce(0, add)
```

### Error Handling

```um
fn parseNum(String s) -> Result<Int, String> =
	// ...

fn main() -> Result<Unit, String> = {
	let n = match parseNum("42") {
		Ok value -> value
		Err msg -> return Err msg
	}
	Ok ()
}
```

### Modules

```um
module math

import strings
import io

fn pi() -> Float64 = 3.141592653589793
```

All `.um` files in a directory share a module. Module name matches directory name.

---

## Open Questions

- Typeclass / constraint syntax for parametric polymorphism
- Region/arena syntax; any explicit lifetime annotation?
- Concurrency primitives: goroutine spawn syntax, channel types
- String interpolation
- Operator overloading policy
- FFI design (C interop)
- Exact grammar for newline-as-statement-terminator edge cases
