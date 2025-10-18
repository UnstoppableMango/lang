open Ast

let print = function
  | `String s -> print_endline s
  | `Number n -> print_endline n
