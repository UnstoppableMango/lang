open Angstrom

let parse (s : string) = parse_string ~consume:All (many (char 'a')) s
