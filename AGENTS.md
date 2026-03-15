# Agent Instructions for StreamPeek

This file provides context and instructions for AI agents (Antigravity, etc.) working on this repository.

## Tech Stack
- **Language:** Go (Golang)
- **Domain:** Apache Kafka REST Proxy
- **Focus:** High performance, low latency, JVM-free.

## Shipping Process
All changes must follow the Feature Branch Workflow. Use the local workflow defined in `/.agent/workflows/ship.md`.
- **Primary Tool:** GitHub CLI (`gh`)
- **Main Branch:** `main` (Protected)
- **Naming Convention:** `feat/feature-name` or `fix/issue-name`

## Coding Standards
- **Performance:** Avoid unnecessary allocations in the hot path (message production/consumption).
- **Concurrency:** Use idiomatic Go channels and context for timeout handling.
- **Commits:** Use Conventional Commits (e.g., `feat:`, `fix:`, `perf:`, `chore:`).

## Available Workflows
- **Ship:** `/ship` - Automates branch, commit, push, and PR.
- **Test:** `/test` - Runs `go test ./...` and benchmarks.

---
*Note: If a task involves Enterprise features, ensure you are working in the `StreamPeek-enterprise` repository, not this one.*