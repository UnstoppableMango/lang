module Kaleidoscope

open FParsec
open Swensen.Unquote
open UnMango.Lang
open Xunit
open FsCheck
open FsCheck.Xunit

[<Theory>]
[<InlineData("+", '+')>]
[<InlineData("-", '-')>]
[<InlineData("*", '*')>]
[<InlineData("<", '<')>]
let ``Should parse a binop`` (i: string, o: char) =
  test <@ run Kalaidoscope.parseBinOp i |> Util.parseSuccess = o @>

[<Property>]
let ``Should parse a string between parentheses`` (UnicodeString s) =
  let i = "(" + s + ")"
  test <@ run (Kalaidoscope.betweenParens (anyString s.Length)) i |> Util.parseSuccess = s @>
