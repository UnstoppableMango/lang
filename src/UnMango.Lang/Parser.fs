module UnMango.Lang.Parser

open FParsec
open UnMango.Lang.Ast

let validStringLiteralChar c = c <> '\\' && c <> '"'

let stringLiteral: Parser<Node, _> =
    manySatisfy validStringLiteralChar |> between (pchar '"') (pchar '"') |>> String

let nodeList: Parser<Node list, _> = many stringLiteral

let file: Parser<File, _> = nodeList |>> (fun s -> { Nodes = s })

let parse (input: string) =
    match run file input with
    | Success(result, _, _) -> Result.Ok(result)
    | Failure(msg, _, _) -> Result.Error(msg)
