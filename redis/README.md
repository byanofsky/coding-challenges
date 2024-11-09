# Build Your Own Redis Server

Coding Challenge: https://codingchallenges.fyi/challenges/challenge-redis

## Implementation Notes

- RESP protocol: https://redis.io/docs/latest/develop/reference/protocol-spec/
- Handling generics for a function that accepts a union of types while retaiing type safety. See Serialize function and Serializable type
- No null type. Created a null struct.
- Parsing for array, returning remaining. Use a regex to ensure no remaining values at end.
