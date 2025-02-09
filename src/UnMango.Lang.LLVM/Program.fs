open System
open LLVMSharp
open LLVMSharp.Interop
open UnMango.Lang

let ctx = LLVMContext()
let m = LLVMModuleRef.CreateWithName("test")

let v = IR.emitString ctx (Ast.String("test"))

Console.WriteLine(v.ToString())
