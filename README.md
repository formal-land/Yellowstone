# Yellowstone
> A free-monad library compatible with any programming languages.

![logo](static/yellowstone.png)

<a href="https://www.flaticon.com/fr/icones-gratuites/yellowstone" title="yellowstone icons">Yellowstone icon created by Chanut-is-Industries - Flaticon</a>

## Aim?

As more and more code is written using AI, and vulnerabilities discovery increases with AI, it becomes increasingly important to be able to trust and understand the code we use.

Yellowstone aims to increasing the understand of the code by following functional programming principles:

- Pure functions
- Everything else in a free-monad to represent side effects

This uniform representation of effects allows to:

- Generate code diagrams from execution traces
- Write unit-tests with better reproducibility
- Connect the code to formal verification tools ([Rocq](https://rocq-prover.org/)) to prove general security properties

This approach is related to the [Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)) or existing libraries like [Effect TS](https://github.com/effect-ts/effect).

## Implementation

To be compatible with existing programming languages and codebases, our approach is as follows:

- Make an extension of your code/programming language with additional information to represent the explicitly effects.
- Have a tool to erase these additional information to get back the original code.

This is verbose, but:

- Verbosity should be less of an issue with AI when it is checked (here the check is to get back the original code).
- This gives full control over how to plug into an existing codebase, and how to serialize the effects.
- This is compatible with more programming languages; for example, the Effect library of TypeScript relies on the existence of generators in the language.
