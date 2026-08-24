/* Throwaway hello-world compiler: parses `print "..."` and emits LLVM IR text.
 * Carries no design commitment; see the repo workflow docs. */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void die(const char *msg)
{
	fprintf(stderr, "unmangc: %s\n", msg);
	exit(1);
}

int main(int argc, char **argv)
{
	if (argc != 2)
		die("usage: unmangc <file>");

	FILE *f = fopen(argv[1], "r");
	if (!f)
		die("cannot open source file");

	char line[1024];
	if (!fgets(line, sizeof line, f))
		die("empty source file");
	fclose(f);

	line[strcspn(line, "\n")] = '\0';

	const char *keyword = "print \"";
	if (strncmp(line, keyword, strlen(keyword)) != 0)
		die("expected: print \"...\"");

	char *str = line + strlen(keyword);
	char *end = strrchr(str, '"');
	if (!end || end[1] != '\0')
		die("unterminated string literal");
	*end = '\0';

	/* Escape sequences are copied through as raw bytes; good enough here. */
	size_t len = strlen(str);

	printf("@.str = private unnamed_addr constant [%zu x i8] c\"%s\\00\"\n",
	       len + 1, str);
	printf("\n");
	printf("declare i32 @puts(ptr)\n");
	printf("\n");
	printf("define i32 @main() {\n");
	printf("  call i32 @puts(ptr @.str)\n");
	printf("  ret i32 0\n");
	printf("}\n");

	return 0;
}
