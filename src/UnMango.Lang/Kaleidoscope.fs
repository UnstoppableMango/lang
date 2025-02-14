module UnMango.Lang.Kalaidoscope

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

let parseNumberExpr: Parser<ExprAST, unit> = pfloat |>> NumberExprAST

let parseVariableExpr: Parser<ExprAST, unit> =
    identifier (IdentifierOptions()) |>> VariableExprAST

let betweenParens x = x |> between (pchar '(') (pchar ')')

let parseBinOp: Parser<char, unit> = anyOf "<+-*"

// let parseExpr = sepBy parsePrimary parseBinOp

// let parseParenExpr = parseExpr |> betweenParens

// let parseCallExpr =
//     (fun name args -> CallExprAST(name, args))
//     |> pipe2 (identifier (IdentifierOptions())) (sepBy parseExpr (pstring ",") |> betweenParens)
//     |>> fun x -> x :> ExprAST

// let parseIdentifierExpr = parseVariableExpr <|> parseCallExpr

// let parsePrimary = parseIdentifierExpr <|> parseNumberExpr <|> parseParenExpr
