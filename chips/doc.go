// Chips is simple wrapper around float originally as an replacement
// for decimal.Decimal used in project.
//
// Reason for not using float32 directy is fact we can later
// add assembly instractions for basic math to improve performance.
//
// Same additional methods such as Abs, Pow, etc. can be useful are available
// to make it easier to work with chips.
package chips
