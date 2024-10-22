# Build your own Regex Engine

This solves the bonus problem from the grep coding challenge: https://codingchallenges.fyi/challenges/challenge-grep#going-beyond-this

> if you took the easy path, you could go back and try implementing your own regex engine.

## Implementation Notes

- Package structure. helpers to contain all non-public. types to contain structs and constructors. my_regexp.go to contain exported functions.
- Define a compile function. Similar to regexp, this will be repsonsible for taking pattern and outputting a regexp struct, later used to match
- First step, as known from compilers, is to scan and create stream of tokens. Instead of stream, will create a slice. For each, provide a kind (type wasn't valid field) and the char, in case it is needed
- Start with types defined here: https://pkg.go.dev/regexp/syntax
  - x - single character
  - . - any characters
  - - - repetition, 0 or more
  - \w - word character
  - \n - new line
- Consider grouping these, with sub types
- For now, no sub types. Just limit to a subset of regex.
- Optimizations: Do not append to tokens
- Compile creates a set of matcher functions depending on tokens. These matcher functions are passed to regexp instance and will be used during matching.
- Matcher function. Variation of visitor pattern (or my understanding of that pattern). Delegate matching to each matcher. Decides that to do next depending on current match result. Will enable such things as repition.
- a next function is called from each matcher to determine what to call next. This is controlled in outer function.
- For start character, will need to intercept in the for loop to prevent continuing to execute matchers against each char in string
- Refactoring parse to align with steps in compilation process
- Improvements to test according to Claude using a loop
- Refactor repitition matcher to a generic. Use compisition for each manager
- Use a map to determine if is repitiion character and get constructor for repitition matcher
