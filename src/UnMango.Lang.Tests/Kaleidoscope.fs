module Kaleidoscope

open FParsec
open Swensen.Unquote
open UnMango.Lang
open Xunit
open FsCheck
open FsCheck.Xunit
open System

[<Theory>]
[<InlineData("+", '+')>]
[<InlineData("-", '-')>]
[<InlineData("*", '*')>]
[<InlineData("<", '<')>]
let ``Should parse a binop`` (i: string, o: char) =
    test <@ run Kaleidoscope.binOp i |> Util.parseSuccess = o @>

[<Property>]
let ``Should parse a string between parentheses`` (UnicodeString s) =
    let i = "(" + s + ")"
    test <@ run (Kaleidoscope.betweenParens (anyString s.Length)) i |> Util.parseSuccess = s @>

[<Property>]
let ``Should parse a float`` (NormalFloat f) =
    let i = Convert.ToString f

    test <@ run Kaleidoscope.numberExpr i |> Util.parseSuccess = Kaleidoscope.NumberExprAST f @>

[<Fact(Skip = "Still working this out")>]
let ``Should parse a variable identifier`` () =
    test <@ run Kaleidoscope.identifierExpr "test" |> Util.parseSuccess = Kaleidoscope.VariableExprAST "test" @>

[<Fact(Skip = "Still working this out")>]
let ``Should parse a call expression identifier`` () =
    test <@ run Kaleidoscope.identifierExpr "test()" |> Util.parseSuccess = Kaleidoscope.CallExprAST("test", []) @>
