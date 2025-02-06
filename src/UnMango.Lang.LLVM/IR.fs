module UnMango.Lang.IR

open LLVMSharp
open UnMango.Lang.Ast

let emitString c (String s) : Value =
    ConstantDataArray.GetString(c, s)
