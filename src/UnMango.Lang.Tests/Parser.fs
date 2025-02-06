module Parser

open System
open FsCheck
open FsCheck.FSharp
open FsCheck.Xunit
open Swensen.Unquote
open UnMango.Lang

type String =
    static member Double() =
        ArbMap.defaults
        |> ArbMap.arbitrary<NonEmptyString>
        |> Arb.filter (fun (NonEmptyString s) -> s |> String.forall Char.IsLetterOrDigit)

let parseSuccess =
    function
    | Ok(x) -> x
    | Error(msg, _) -> failwith msg

[<Property(Arbitrary = [| typeof<String> |])>]
let ``Should parse a string`` (NonEmptyString s) =
    test <@ Parser.parse ("\"" + s + "\"") |> parseSuccess = { Nodes = [ (Ast.String s) ] } @>
