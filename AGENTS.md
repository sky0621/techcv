# Repository Guidelines

## Project Structure & Module Organization
The repository currently hosts curated Markdown docs describing the engineer's technical skillset. Core content lives in `README.md`, while licensing info sits in `LICENSE`. Add new source files under `docs/` when expanding the CV; reserve `assets/` for images or PDFs shared across pages.

## Build, Test, and Development Commands
This project is Markdown-first and has no mandatory build step. Use `npx markdownlint README.md AGENTS.md` to catch style issues before pushing. If you introduce compiled assets, document the command set in this file and provide scripts via a `package.json`.

## Coding Style & Naming Conventions
Write plain Markdown with semantic headings that match CV sections (e.g., `## Framework Expertise`). Use sentence case for headings, bold for role titles, and ordered lists when chronology matters. Keep line length under 100 characters to ease diff reviews. Image or asset filenames should be kebab-case (`frontend-stack.png`) and stored in `assets/`.

## Testing Guidelines
Manual review remains the primary validation—cross-check links, spellings, and skill lists each PR. When adding scripts or data files, provide reproducible checks (e.g., a validation script or link to dataset source) and mention them in the PR body. Flag any external dependencies in the README.

## Commit & Pull Request Guidelines
Follow concise, imperative commit messages (`Add backend proficiency section`). Scope each commit to one logical change. PRs should summarize the narrative of the CV update, link any related tracking issue, and include before/after screenshots if you modify visual assets.

## Security & Configuration Tips
Do not store proprietary data or credentials in the repo. For private work samples, link to redacted external references instead of embedding files. If you need environment variables for future automation, provide sample values in `.env.example`.
