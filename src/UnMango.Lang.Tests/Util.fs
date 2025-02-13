module Util

open FParsec

let parseOk =
    function
    | Result.Ok x -> x
    | Result.Error e -> failwith e

let parseSuccess =
  function
  | Success(x, _, _) -> x
  | Failure(e, _, _) -> failwith e
