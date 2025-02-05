module UnMango.Lang.Parser

open FParsec
open UnMango.Lang.Ast

let pstring: Parser<_, _> = pchar '"' >>. manySatisfy ((<>) '"') >>. pchar '"'

let parse (input: string) =
    let t = run pstring input
    { Nodes = [] }
