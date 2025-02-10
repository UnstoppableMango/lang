module rec UnMango.Lang.Kalaidoscope

open FParsec

[<AbstractClass>]
type ExprAST() = class end

type NumberExprAST(Val: float) =
    inherit ExprAST()
    member _.Val = Val

type VariableExprAST(Name: string) =
    inherit ExprAST()
    member _.Name = Name

type BinaryExprAST(Op: char, LHS: ExprAST, RHS: ExprAST) =
    inherit ExprAST()
    member _.Op = Op
    member _.LHS = LHS
    member _.RHS = RHS

type CallExprAST(Callee: string, Args: ExprAST list) =
    inherit ExprAST()
    member _.Callee = Callee
    member _.Args = Args

type PrototypeAST(Name: string, Args: string list) =
    inherit ExprAST()
    member _.Name = Name
    member _.Args = Args

type FunctionAST(Proto: PrototypeAST, Body: ExprAST) =
    inherit ExprAST()
    member _.Proto = Proto
    member _.Body = Body

let binOpPrec =
    function
    | '<' -> 10
    | '+' -> 20
    | '-' -> 30
    | '*' -> 40
    | _ -> -1

let parseNumberExpr = pfloat |>> NumberExprAST |>> fun x -> x :> ExprAST

let parseVariableExpr =
    identifier (IdentifierOptions()) |>> VariableExprAST |>> (fun x -> x :> ExprAST)

let betweenParens x = between (pchar '(') (pchar ')') x

let parseBinOp = anyOf "<+-*"

let parseExpr = sepBy parsePrimary parseBinOp

let parseParenExpr = parseExpr |> betweenParens

let parseCallExpr =
    (fun name args -> CallExprAST(name, args))
    |> pipe2 (identifier (IdentifierOptions())) (sepBy parseExpr (pstring ",") |> betweenParens)
    |>> fun x -> x :> ExprAST

let parseIdentifierExpr = parseVariableExpr <|> parseCallExpr

let parsePrimary = parseIdentifierExpr <|> parseNumberExpr <|> parseParenExpr
