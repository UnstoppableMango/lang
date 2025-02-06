module UnMango.Lang.Parser

open FParsec
open UnMango.Lang.Ast

let pstring: Parser<_, _> =
    manySatisfy ((<>) '"') |> between (pchar '"') (pchar '"') |>> String

let pnodelist: Parser<_, _> = many pstring

let pfile: Parser<_, _> = pnodelist |>> (fun s -> { Nodes = s })

let parse (input: string) =
    match run pfile input with
    | Success(result, _, _) -> Result.Ok(result)
    | Failure(msg, err, _) -> Result.Error(msg, err)
