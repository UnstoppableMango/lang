module UnMango.Lang.Ast

type Node = String of string

type File = { Nodes: Node list }
