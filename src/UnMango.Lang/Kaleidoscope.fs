module rec UnMango.Lang.Kaleidoscope

open FParsec

type PrototypeAST = { Name: string; Args: string list }

type ExprAST =
    | NumberExprAST of float
    | VariableExprAST of string
    | BinaryExprAST of op: char * lhs: ExprAST * rhs: ExprAST
    | CallExprAST of callee: string * args: ExprAST list
    | PrototypeAST of name: string * args: string list
    | FunctionAST of proto: PrototypeAST * body: ExprAST

let binOpPrec =
    function
    | '<' -> 10
    | '+' -> 20
    | '-' -> 30
    | '*' -> 40
    | _ -> -1

let (&&>) a b = fun x -> a x && b x

let (||>) a b = fun x -> a x || b x

let numberExpr: Parser<ExprAST, unit> = pfloat |>> NumberExprAST

let identifier = CharParsers.identifier (IdentifierOptions())

let variableExpr: Parser<ExprAST, unit> = identifier |>> VariableExprAST

let betweenParens x = x |> between (pchar '(') (pchar ')')

let binOp: Parser<char, unit> = anyOf "<+-*"

let parseExpr () = primary // TODO

let parenExpr = parseExpr () |> betweenParens

let argList = sepBy (parseExpr ()) (pstring ",") |> betweenParens

let callExpr =
    (fun name args -> CallExprAST(name, args)) |> pipe2 identifier argList

let identifierExpr = callExpr <|> variableExpr

let primary = identifierExpr <|> numberExpr <|> parenExpr

let tmp = identifier >>. followedBy (satisfy ((=) '('))
