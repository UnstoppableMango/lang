; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [7 x i8] c"count:\00"
@.str.1 = private unnamed_addr constant [6 x i8] c"%lld\0A\00"

declare i32 @puts(ptr)

declare i32 @printf(ptr, ...)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str)
  %call1 = call i32 (ptr, ...) @printf(ptr @.str.1, i64 39)
  ret i32 0
}
