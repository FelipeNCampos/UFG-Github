package cmd

import (
	"errors"
	"testing"
)

func TestParseExecutionStrategyDefault(t *testing.T) {
	strategy, err := ParseExecutionStrategy("")
	if err != nil {
		t.Fatalf("erro inesperado para string vazia: %v", err)
	}
	if strategy != StrategyAuto {
		t.Fatalf("esperava StrategyAuto, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyAuto(t *testing.T) {
	strategy, err := ParseExecutionStrategy("auto")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyAuto {
		t.Fatalf("esperava StrategyAuto, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyHTTP(t *testing.T) {
	strategy, err := ParseExecutionStrategy("http")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyHTTP {
		t.Fatalf("esperava StrategyHTTP, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyDireto(t *testing.T) {
	strategy, err := ParseExecutionStrategy("direto")
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if strategy != StrategyDirect {
		t.Fatalf("esperava StrategyDirect, obteve %v", strategy)
	}
}

func TestParseExecutionStrategyInvalid(t *testing.T) {
	_, err := ParseExecutionStrategy("invalido")
	if err == nil {
		t.Fatalf("esperava erro de validacao para estrategia invalida")
	}

	var exitErr *ExitError
	if !errors.As(err, &exitErr) || exitErr.Code != 2 {
		t.Fatalf("esperava erro de validacao com codigo 2")
	}
}

func TestParseExecutionStrategyCaseInsensitive(t *testing.T) {
	testCases := []struct {
		input    string
		expected ExecutionStrategy
	}{
		{"AUTO", StrategyAuto},
		{"HTTP", StrategyHTTP},
		{"DIRETO", StrategyDirect},
		{"Auto", StrategyAuto},
	}

	for _, tc := range testCases {
		strategy, err := ParseExecutionStrategy(tc.input)
		if err != nil {
			t.Fatalf("erro inesperado para %q: %v", tc.input, err)
		}
		if strategy != tc.expected {
			t.Fatalf("para entrada %q: esperava %v, obteve %v", tc.input, tc.expected, strategy)
		}
	}
}
