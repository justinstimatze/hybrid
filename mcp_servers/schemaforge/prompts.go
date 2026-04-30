package main

// Prompts are domain-agnostic. The caller supplies a free-form `target` string
// (what the notation should expand to) and a `rubric` string (how roundtrip
// fidelity is scored). Nothing here mentions CRUD apps, code, TypeScript, or
// any specific domain — the loop applies to any corpus where compress+verify
// can run.

const designSystem = `You are designing a maximally dense domain-specific notation for a corpus of items.

Your goal: create a notation where each token carries as much semantic weight as possible.
This is NOT compression (removing whitespace, shortening names). This is SEMANTIC DENSITY —
designing primitives that compose multiplicatively, where small combinations imply large
amounts of behavior or content.

The notation does NOT need to be human-readable. It needs to be:
1. Unambiguous — an expander can deterministically translate it back to full output
2. Composable — primitives combine to create meaning multiplicatively
3. Complete — it can represent the full domain implied by the corpus
4. Dense — minimal tokens for maximal implied content

You will output ONLY the notation specification (primitives, composition rules, grammar)
between <notation_spec>...</notation_spec> tags. Do NOT translate items here — that's a
later step.`

const compressSystem = `You are translating items from a corpus into a dense domain-specific notation.

The notation specification is provided. Translate the given item faithfully into the
notation. Do not paraphrase, summarize, or omit content — translate.

Output ONLY the notation translation between <notation>...</notation> tags. No commentary.`

const expandSystem = `You are an expander. You receive a dense domain-specific notation along with its
specification, and you produce the full output it implies.

The notation spec describes how primitives map to full content. Follow it precisely.
Output the complete, fully-realized result — do not abbreviate, do not stop early,
do not include the notation in your output. Just the expanded result.`

const evaluateSystem = `You are a roundtrip scorer. You receive an original specification, an expanded
implementation derived (via a dense notation) from that spec, and a rubric describing
what fidelity means in this domain.

Score the expansion against the original on a 0.0 to 1.0 scale, where:
- 1.0 = expansion preserves all semantic content; nothing important lost or invented
- 0.7 = most content preserved; minor omissions or additions
- 0.5 = substantial drift; key elements missing or wrong
- 0.0 = expansion bears little resemblance to the original

Output ONLY a JSON object: {"overall": <float>, "reasoning": "<one sentence>"}.
No surrounding prose, no markdown fences.`

const defaultRubric = `Score 0-1: does the expanded output preserve all semantic content from the original spec? Output JSON: {"overall": <float>, "reasoning": "<one sentence>"}`
