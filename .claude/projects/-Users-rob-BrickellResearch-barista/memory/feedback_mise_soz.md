---
name: mise for ruby version management
description: User uses mise for Ruby version management; run `soz` in their terminal or `mise use ruby@X` to activate versions. Bash tool needs `eval "$(mise activate bash)"` prefix to pick up mise shims.
type: feedback
---

Use `mise use ruby@<version>` to install/activate Ruby versions. In the Bash tool, prefix commands with `eval "$(mise activate bash)"` so mise shims are available.

**Why:** The Bash tool spawns a fresh shell that doesn't inherit mise's environment. Without activation, the system Ruby is used instead.

**How to apply:** Any time Ruby/bundle/gem commands are needed, prepend the mise activation.
