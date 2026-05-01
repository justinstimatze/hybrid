# Contributing

This is research-stage software. The pattern itself is the artifact; the skill is its operational form. Contributions are welcome but the bar for what changes is higher than the bar for what gets discussed.

## Easy first contributions

- **Cross-agent ports.** The skill is portable; the cross-agent stub manifests under `.codex-plugin/`, `.cursor-plugin/`, and `gemini-extension.json` need adjustment per each agent's actual current spec. PRs from users on those platforms are the right path.
- **Catch a documented mistake.** Wrong citation, broken link, factual error in the prior-art tiers, voice slipping out of register. Small, real, welcome.
- **Doc edits.** Copyedits welcome. New entries in `skills/hybrid-loops/references/BLOCK_GRAPHS.md` for shapes you've actually built. Domain-expert on-ramp drafts written for non-engineering audiences. Worked-example writeups for any of the runnable instances (your own or a project you've forked from).
- **Tried this in a real project? Tell me about it.** Whether it helped, where it didn't, what shape it ended up looking like. Negative results especially welcome.

## What changes the bar is high for

- **Renaming the pattern or its roles.** "Hybrid loops" / "lens / substrate / gate / reasoner / action" are the working internal vocabulary. They're not perfect; changing them is disruptive enough to need real argument.
- **Major skill structural changes.** The 5-phase diagnostic (find → scope → shape → quick design → scaffold) has been validated by trigger tests. Restructuring needs an equally rigorous validation.

## Commit style

Semantic commits, one logical change per commit. Co-authored-by trailers welcome when the work was meaningfully assisted by an AI tool.

## Code of conduct

See `CODE_OF_CONDUCT.md`. Short version: be kind, assume good faith, don't drag personal disputes into Issues.

## License

By contributing, you agree your contributions are licensed under the MIT License (see `LICENSE`).
