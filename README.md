### Writing an Interpreter in Go

Based on "Writing an Interpreter in Go" by Thorsten Ball :)

```shell-session
$ go run main.go
Hello tphamdo!. This is the Monkey programming language!
Feel free to type in commands
>> let x =2;
let x = 2;
>> let
            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
Woops! We ran into some monkey business here!
 parser errors:
        Expected next token to be IDENT. Got EOF instead
```
