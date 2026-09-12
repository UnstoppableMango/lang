; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [6 x i8] c"first\00"
@.str.1 = private unnamed_addr constant [7 x i8] c"second\00"
@.str.2 = private unnamed_addr constant [6 x i8] c"third\00"
@.str.3 = private unnamed_addr constant [7 x i8] c"fourth\00"

declare i32 @puts(ptr)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str)
  %call1 = call i32 @puts(ptr @.str.1)
  %call2 = call i32 @puts(ptr @.str.2)
  %call3 = call i32 @puts(ptr @.str.3)
  ret i32 0
}
