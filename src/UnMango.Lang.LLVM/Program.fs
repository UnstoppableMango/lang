open System
open LLVMSharp
open UnMango.Lang

let ctx = LLVMContext()
let m = Module()

let v = IR.emitString ctx (Ast.String("test"))

Console.WriteLine (v.ToString())
