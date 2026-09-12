; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [14 x i8] c"never printed\00"
@.str.1 = private unnamed_addr constant [10 x i8] c"only this\00"

declare i32 @puts(ptr)

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str.1)
  ret i32 0
}
