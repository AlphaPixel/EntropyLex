export SHELL := /bin/bash

PYTHON ?= python

# ____ all __________________________________________________________________

.PHONY: all

all: test lint

# ____ lint _________________________________________________________________

.PHONY: fmt check typecheck vuln lint

fmt:
	@printf "\nRunning Ruff formatter...\n"
	$(PYTHON) -m ruff format .

check:
	@printf "\nRunning Ruff checks...\n"
	$(PYTHON) -m ruff check --fix .

typecheck:
	@printf "\nRunning mypy...\n"
	$(PYTHON) -m mypy .

vuln:
	@printf "\nRunning pip-audit...\n"
	$(PYTHON) -m pip_audit --local

lint: fmt check typecheck vuln
	@printf "\nChecking for changed files...\n"
	git status --porcelain
	@printf "\n"
	@if [[ "$${CI}" == "true" ]]; then $$(exit $$(git status --porcelain | wc -l)); fi

# ____ test _________________________________________________________________

.PHONY: unit_test test

unit_test:
	@printf "\nRunning Python unit tests...\n"
	$(PYTHON) -m coverage run -m unittest discover -s . -p "test_*.py" -v
	$(PYTHON) -m coverage report

test: unit_test

# ____ clean artifacts ______________________________________________________

.PHONY: clean

clean:
	@printf "\nClean...\n"
	-rm -rf .coverage .mypy_cache .ruff_cache build dist *.egg-info
	-find . -type d -name __pycache__ -prune -exec rm -rf {} +
