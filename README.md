# pi2go
Compile pi calculus expressions to golang

We use a syncronous pi calculus with parametric recursion. 
The grammar is taken from [SLMC](http://ctp.di.fct.unl.pt/SLMC/)

```
lower = ['a'-'z']
upper = ['A'-'Z']
letter = lower | upper
digit  = ['0'-'9']
name = lower (letter | digit | '_')*
namelist := 	epsilon | name (',' name)*
prefix 	:=	name!(namelist)
	|	name?(namelist)
	|	[name = name]
	|	[name != name]
	| 	tau
process := 	0
	|	process | process
	| 	'new' namelist 'in' process
	|	prefix.process
	| 	select{ prefix.process (';' prefix.process)* }
	|	CapsID(namelist)
	| 	( process )

```

# Roadmap
- parse multiple lines (sequence of processes)
- support process definitions
- compiler


