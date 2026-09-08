//! Source-anchored diagnostics.
//!
//! A diagnostic is a message plus a byte offset into the source text, and
//! nothing about it is specific to parsing: later phases report the same way.

/// Byte offset of `tail` within `source`.
///
/// `tail` must be a subslice of `source`, which every parser input is, since
/// the parsers only ever narrow the one string read from disk.
pub fn offset_of(source: &str, tail: &str) -> usize {
    tail.as_ptr() as usize - source.as_ptr() as usize
}

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Diagnostic {
    pub message: String,
    pub offset: usize,
}

impl Diagnostic {
    pub fn new(message: impl Into<String>, offset: usize) -> Self {
        Diagnostic {
            message: message.into(),
            offset,
        }
    }

    /// Anchor a diagnostic at the start of `tail`, a subslice of `source`.
    pub fn at(message: impl Into<String>, source: &str, tail: &str) -> Self {
        Self::new(message, offset_of(source, tail))
    }

    /// Render in the rustc style. The path is a parameter because a
    /// diagnostic never knows where its source text came from.
    pub fn render(&self, path: &str, source: &str) -> String {
        let (line_no, col, line) = locate(source, self.offset);
        let gutter = " ".repeat(line_no.to_string().len());
        let pad: String = line
            .chars()
            .take(col - 1)
            .map(|c| if c == '\t' { '\t' } else { ' ' })
            .collect();

        format!(
            "error: {message}\n{gutter}--> {path}:{line_no}:{col}\n{gutter} |\n{line_no} | {line}\n{gutter} | {pad}^\n",
            message = self.message,
        )
    }
}

/// Resolve a byte offset to a one-based line number, a one-based column
/// counted in characters, and the text of that line without its terminator.
fn locate(source: &str, offset: usize) -> (usize, usize, &str) {
    let mut offset = offset.min(source.len());

    // At end of input the offset points past the final newline, which is a
    // line that does not exist. Back up to the end of the last real line.
    if offset == source.len() {
        while offset > 0 && matches!(source.as_bytes()[offset - 1], b'\n' | b'\r') {
            offset -= 1;
        }
    }

    let start = source[..offset].rfind('\n').map_or(0, |i| i + 1);
    let end = source[start..]
        .find('\n')
        .map_or(source.len(), |i| start + i);

    let line_no = source[..start].matches('\n').count() + 1;
    let col = source[start..offset].chars().count() + 1;
    let line = source[start..end].trim_end_matches('\r');

    (line_no, col, line)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn renders_on_the_first_line() {
        let source = "print 42\n";
        let d = Diagnostic::new("expected a string literal", 6);
        assert_eq!(
            d.render("hello.lang", source),
            "error: expected a string literal\n \
             --> hello.lang:1:7\n  \
             |\n1 | print 42\n  |       ^\n"
        );
    }

    #[test]
    fn renders_mid_file_with_a_wide_gutter() {
        let mut source = "print \"a\"\n".repeat(9);
        source.push_str("print 42\n");
        let offset = source.find("42").unwrap();
        let d = Diagnostic::new("expected a string literal", offset);
        assert_eq!(
            d.render("hello.lang", &source),
            "error: expected a string literal\n  \
             --> hello.lang:10:7\n   \
             |\n10 | print 42\n   |       ^\n"
        );
    }

    #[test]
    fn renders_at_end_of_input_on_the_last_real_line() {
        let source = "print \"a\"\n";
        let d = Diagnostic::new("unexpected end of input", source.len());
        assert_eq!(
            d.render("hello.lang", source),
            "error: unexpected end of input\n \
             --> hello.lang:1:10\n  \
             |\n1 | print \"a\"\n  |          ^\n"
        );
    }

    #[test]
    fn pads_the_caret_with_the_tabs_it_passes_over() {
        let source = "\tprint 42\n";
        let d = Diagnostic::new("expected a string literal", 7);
        assert_eq!(
            d.render("hello.lang", source),
            "error: expected a string literal\n \
             --> hello.lang:1:8\n  \
             |\n1 | \tprint 42\n  | \t      ^\n"
        );
    }

    #[test]
    fn offset_of_finds_a_subslice() {
        let source = "print 42";
        assert_eq!(offset_of(source, &source[6..]), 6);
    }
}
