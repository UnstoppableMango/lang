module Parser

open FsCheck
open FsCheck.FSharp
open FsCheck.Xunit
open Swensen.Unquote
open UnMango.Lang

type String =
    static member Double() =
        ArbMap.defaults
        |> ArbMap.arbitrary<UnicodeString>
        |> Arb.filter (fun (UnicodeString s) -> s |> String.forall (fun c -> c <> '\\' && c <> '"'))

let parseSuccess =
    function
    | Ok(x) -> x
    | Error(msg, _) -> failwith msg

[<Property(Arbitrary = [| typeof<String> |])>]
let ``Should parse a quoted unicode string`` (UnicodeString s) =
    test <@ Parser.parse ("\"" + s + "\"") |> parseSuccess = { Nodes = [ (Ast.String s) ] } @>
