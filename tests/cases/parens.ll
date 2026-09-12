; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [6 x i8] c"%lld\0A\00"

declare i32 @puts(ptr)

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %call = call i32 (ptr, ...) @printf(ptr @.str, i64 9)
  ret i32 0
}
