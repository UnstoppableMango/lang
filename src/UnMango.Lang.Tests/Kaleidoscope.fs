module Kaleidoscope

open FParsec
open Swensen.Unquote
open UnMango.Lang.Kaleidoscope
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
    test <@ run binOp i |> Util.parseSuccess = o @>

[<Property>]
let ``Should parse a string between parentheses`` (UnicodeString s) =
    let i = "(" + s + ")"
    test <@ run (betweenParens (anyString s.Length)) i |> Util.parseSuccess = s @>

[<Property>]
let ``Should parse a float`` (NormalFloat f) =
    let i = Convert.ToString f
    test <@ run numberExpr i |> Util.parseSuccess = NumberExprAST f @>

[<Fact(Skip = "Still working this out")>]
let ``Should parse a variable identifier`` () =
    test <@ run identifierExpr "test" |> Util.parseSuccess = VariableExprAST "test" @>

[<Fact(Skip = "Still working this out")>]
let ``Should parse a call expression identifier`` () =
    test <@ run identifierExpr "test()" |> Util.parseSuccess = CallExprAST("test", []) @>

[<Fact>]
let ``Should work the way I think it works`` () =
    test <@ run tmp "test()" |> Util.parseSuccess = CallExprAST("test", []) @>
