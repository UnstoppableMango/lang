; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [3 x i8] c"hi\00"

declare i32 @puts(ptr)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str)
  ret i32 0
}
