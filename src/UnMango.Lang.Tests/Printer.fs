module Printer

open FsCheck
open FsCheck.Xunit
open Swensen.Unquote
open UnMango.Lang

[<Property>]
let ``Should print a unicode string`` (UnicodeString s) =
    test <@ { Ast.File.Nodes = [Ast.String s] } |> Printer.sprint = $"\"{s}\"" @>
