# Vendor JavaScript Libraries

This directory contains third-party JavaScript libraries used in the frontend.

## Libraries

### Alpine.js
- **Version**: 3.x (minified)
- **License**: MIT
- **Source**: https://alpinejs.dev/
- **Size**: 44 KB

### Tom Select
- **Version**: 2.3.1 (canonical `tom-select.complete.min.js` build)
- **License**: Apache-2.0
- **Source**: https://tom-select.js.org/
- **Size**: 50 KB
- **Source maps**: `tom-select.complete.min.js.map` (JS) and `../css/tom-select.min.css.map`
  (CSS) are vendored alongside the minified files. They are the unmodified upstream maps
  for v2.3.1 (verified byte-identical), with only the `file` field and the
  `sourceMappingURL` comment adjusted to match the local file names.

### Flatpickr
- **Version**: 4.6.13 (minified)
- **License**: MIT
- **Source**: https://flatpickr.js.org/
- **Size**: 49 KB

## License Compliance

All libraries above are third-party dependencies used under their respective licenses:
Alpine.js and Flatpickr are MIT-licensed; Tom Select is Apache-2.0-licensed. Their use
in this project (which is licensed under the ISC License) is fully compliant with all of
these license terms. This README serves as the required license attribution.

### License Permissions Summary
Both MIT and Apache-2.0 are permissive licenses:
- ✅ Can be used commercially
- ✅ Can be modified
- ✅ Can be distributed
- ✅ Can be used privately
- ⚠️ **Requirement**: Include license notice with distributions

## Notes

These are production minified versions. Source maps for Tom Select are vendored
alongside the minified files (see the Tom Select entry above); source maps for the
other libraries are available from their upstream project repositories.
