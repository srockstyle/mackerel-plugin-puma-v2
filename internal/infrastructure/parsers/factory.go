package parsers

import (
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/infrastructure"
)

// Parser interface for different Puma version parsers
type Parser interface {
	Parse(stats *infrastructure.PumaStats) (*domain.MetricCollection, error)
}

// ParserFactory creates appropriate parser based on version
type ParserFactory struct{}

// NewParserFactory creates a new parser factory
func NewParserFactory() *ParserFactory {
	return &ParserFactory{}
}

// GetParser returns appropriate parser for the given version
func (f *ParserFactory) GetParser(version string) Parser {
	switch {
	case isPuma6OrLater(version):
		return NewV6Parser()
	case isPuma5(version):
		return NewV5Parser()
	default:
		return NewV4Parser()
	}
}

// isPuma6OrLater checks if version is Puma 6.x or later
func isPuma6OrLater(version string) bool {
	return version == "6.x" || (len(version) > 0 && version[0] >= '6')
}

// isPuma5 checks if version is Puma 5.x
func isPuma5(version string) bool {
	return version == "5.x" || (len(version) > 0 && version[0] == '5')
}