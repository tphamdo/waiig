package eval

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, e *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node, e)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, e)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, e)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.BlockStatement:
		return evalBlockStatement(node, e)

	case *ast.LetStatement:
		val := Eval(node.Value, e)
		if isError(val) {
			return val
		}
		e.Set(node.Name.Value, val)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: e}

	case *ast.PrefixExpression:
		right := Eval(node.Right, e)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, e)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, e)
		if isError(right) {
			return right
		}

		return evalInfixExpression(left, node.Operator, right)

	case *ast.IfExpression:
		return evalIfExpression(node, e)

	case *ast.Identifier:
		return evalIdentifier(node, e)

	case *ast.CallExpression:
		return evalCallExpression(node, e)

	}

	return nil
}

func evalProgram(program *ast.Program, e *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, e)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(bs *ast.BlockStatement, e *object.Environment) object.Object {
	var result object.Object

	for _, statement := range bs.Statements {
		result = Eval(statement, e)

		if ret, ok := result.(*object.ReturnValue); ok {
			return ret
		}

		if err, ok := result.(*object.Error); ok {
			return err
		}
	}

	return result
}

func nativeBoolToBooleanObject(val bool) object.Object {
	if val {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalNegOperatorExpression(right)
	default:
		return NULL
	}
}

func evalInfixExpression(left object.Object, operator string,
	right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(left, operator, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(left, operator, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(left object.Object, operator string,
	right object.Object) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(left object.Object, operator string,
	right object.Object) object.Object {

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	res := right.(*object.Integer)
	res.Value = -res.Value
	return res
}

func evalIfExpression(ie *ast.IfExpression, e *object.Environment) object.Object {
	if cond := Eval(ie.Condition, e); isTruthy(cond) {
		return evalBlockStatement(ie.Consequence, e)
	} else if ie.Alternative != nil {
		return evalBlockStatement(ie.Alternative, e)
	}

	return NULL
}

func evalIdentifier(ident *ast.Identifier, e *object.Environment) object.Object {
	val, ok := e.Get(ident.Value)

	if !ok {
		return newError("identifier not found: %s", ident.Value)
	}

	return val
}

func evalCallExpression(node *ast.CallExpression, e *object.Environment) object.Object {
	f := Eval(node.Function, e)

	if isError(f) {
		return f
	}

	fn, ok := f.(*object.Function)
	if !ok {
		return newError("not a function: %s", f.Type())
	}

	if len(node.Arguments) != len(fn.Parameters) {
		return newError("Expected %d arguments. Got=%d", len(fn.Parameters), len(node.Arguments))
	}

	// extend function environment
	ne := object.NewEnclosedEnvironment(fn.Env)

	for i := range node.Arguments {
		arg := Eval(node.Arguments[i], e)
		if isError(arg) {
			return arg
		}
		ne.Set(fn.Parameters[i].String(), arg)
	}

	evaluated := Eval(fn.Body, ne)
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		// unwrap return ojbect
		return returnValue.Value
	}
	return evaluated

}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) object.Object {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
