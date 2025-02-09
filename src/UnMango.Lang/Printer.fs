module UnMango.Lang.Printer

open UnMango.Lang.Ast

let printNode =
    function
    | String s -> $"\"{s}\""

let sprint f =
    f.Nodes |> List.map printNode |> String.concat "\n"
