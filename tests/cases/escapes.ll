; ModuleID = 'unmang'
source_filename = "unmang"

@.str = private unnamed_addr constant [10 x i8] c"tab:\09here\00"
@.str.1 = private unnamed_addr constant [12 x i8] c"quoted: \22x\22\00"
@.str.2 = private unnamed_addr constant [11 x i8] c"back\\slash\00"

declare i32 @puts(ptr)

define i32 @main() {
entry:
  %call = call i32 @puts(ptr @.str)
  %call1 = call i32 @puts(ptr @.str.1)
  %call2 = call i32 @puts(ptr @.str.2)
  ret i32 0
}
