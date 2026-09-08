; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [13 x i8] c"line one\0Atwo\00"

declare i32 @puts(ptr)

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str)
  ret i32 0
}
