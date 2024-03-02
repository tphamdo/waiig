### Writing an Interpreter in Go

Based on "Writing an Interpreter in Go" by Thorsten Ball :)

##### REPL:
```shell-session
$ go run main.go
Hello tphamdo!. This is the Monkey programming language!
Feel free to type in commands
>> 5 + 5 * 10
55
>> (5 > 5 == true) != false
false
>> 500 / 2 == 250
true
```

##### Notes:
<ol>
  <li>Interpeters are simple</li>
  <li>Lexer -> Parse -> Eval</li>
  <li>Tokenize input string</li>
  <li>Create AST by parsing tokens. We use Pratt-parsing here</li>
  <li>Traverse AST and eval recursively</li>
</ol>
