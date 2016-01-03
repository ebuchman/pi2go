# pi2go
Compile pi calculus expressions to golang

We use a syncronous pi calculus without replication. The grammar is as follows:

```
Process := process [ "|" Process ]
process := "0" | "(" Process ")" | sum
sum := prefix "." { prefix "." } Process [ "+" sum ]
prefix := fire | pull
fire := ident "!" "(" ident ")"
pull := ident "?" "(" ident ")"
```

Here, `[ ]` encloses an optional term, while `{ }` encloses a term repeated 0 or more times.
Terminals are enclosed in `" "`.

# Roadmap
- parse multiple lines (sequence of processes)
- support process definitions
- compiler


