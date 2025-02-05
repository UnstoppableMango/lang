module UnMango.Lang.Ast

type Node =
    abstract Pos : unit -> int
    abstract End : unit -> int
    
type String = { Val: string }

type File = { Nodes: Node list }
