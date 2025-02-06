module UnMango.Lang.Parser

open FParsec
open UnMango.Lang.Ast

let validStringLiteralChar c = c <> '\\' && c <> '"'

let escapedChar =
    pstring "\\"
    >>. (anyOf "\\nrt\""
         |>> function
             | 'n' -> "\n"
             | 'r' -> "\r"
             | 't' -> "\t"
             | c -> string c)

let stringLiteral: Parser<_, _> =
    (manySatisfy validStringLiteralChar <|> escapedChar)
    |> manyStrings
    |> between (pchar '"') (pchar '"')
    |>> String

let nodeList: Parser<_, _> = many stringLiteral

let file: Parser<_, _> = nodeList |>> (fun s -> { Nodes = s })

let parse (input: string) =
    match run file input with
    | Success(result, _, _) -> Result.Ok(result)
    | Failure(msg, err, _) -> Result.Error(msg, err)
