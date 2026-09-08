//! Golden tests: every `tests/cases/*.lang` file is compiled and its output
//! compared against a sibling golden file.
//!
//! A case that compiles is compared against `<name>.ll`, and a case that fails
//! against `<name>.err`, holding the rendered diagnostic verbatim. Which one
//! applies is decided by the result, so there is no naming convention to
//! remember and no way for a case to test nothing by sitting in the wrong
//! place.
//!
//! Rerun with `UPDATE_GOLDENS=1` (or `make bless`) to rewrite the goldens,
//! then read the diff. Blessing without reading defeats the point.

use std::env;
use std::fs;
use std::path::{Path, PathBuf};

use unmangc::compile;

fn cases_dir() -> PathBuf {
    Path::new(env!("CARGO_MANIFEST_DIR")).join("tests/cases")
}

#[test]
fn goldens() {
    let update = env::var_os("UPDATE_GOLDENS").is_some();

    let mut cases: Vec<PathBuf> = fs::read_dir(cases_dir())
        .expect("tests/cases is missing")
        .map(|entry| entry.expect("cannot read tests/cases").path())
        .filter(|path| path.extension().is_some_and(|ext| ext == "lang"))
        .collect();
    cases.sort();

    // A nix source filter that drops data files would leave this empty and the
    // suite would pass while testing nothing.
    assert!(
        !cases.is_empty(),
        "no .lang cases found in {}",
        cases_dir().display()
    );

    let mut failures = Vec::new();

    for case in cases {
        let source = fs::read_to_string(&case).expect("cannot read case");
        let name = case.file_name().unwrap().to_string_lossy().into_owned();

        let (golden, actual) = match compile(&source) {
            Ok(ir) => (case.with_extension("ll"), ir),
            Err(d) => (case.with_extension("err"), d.render(&name, &source)),
        };

        if update {
            fs::write(&golden, &actual).expect("cannot write golden");
            continue;
        }

        match fs::read_to_string(&golden) {
            Ok(expected) if expected == actual => {}
            Ok(expected) => failures.push(report(&golden, &expected, &actual)),
            Err(_) => failures.push(format!(
                "{}: golden is missing\n---- actual ----\n{actual}----------------",
                golden.display()
            )),
        }
    }

    assert!(
        failures.is_empty(),
        "\n{}\nrerun with UPDATE_GOLDENS=1 to accept",
        failures.join("\n")
    );
}

fn report(golden: &Path, expected: &str, actual: &str) -> String {
    let at = expected
        .lines()
        .zip(actual.lines())
        .position(|(e, a)| e != a)
        .unwrap_or_else(|| expected.lines().count().min(actual.lines().count()));

    format!(
        "{}: golden mismatch at line {}\n  expected: {}\n    actual: {}\n---- full actual output ----\n{actual}----------------------------",
        golden.display(),
        at + 1,
        expected.lines().nth(at).unwrap_or("<end of file>"),
        actual.lines().nth(at).unwrap_or("<end of file>"),
    )
}
