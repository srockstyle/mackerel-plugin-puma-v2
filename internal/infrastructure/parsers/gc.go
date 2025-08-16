package parsers

import (
	"encoding/json"
	"fmt"

	"github.com/srockstyle/mackerel-plugin-puma-v2/internal/domain"
)

// GCStats represents the /gc-stats response (based on lib/gc.go)
type GCStats struct {
	// Ruby 2.0+
	Count                json.Number `json:"count"`
	HeapFinalNum         json.Number `json:"heap_final_num"`
	HeapFreeNum          json.Number `json:"heap_free_num"`
	HeapIncrement        json.Number `json:"heap_increment"`
	HeapLength           json.Number `json:"heap_length"`
	HeapLiveNum          json.Number `json:"heap_live_num"`
	HeapUsed             json.Number `json:"heap_used"`
	TotalAllocatedObject json.Number `json:"total_allocated_object"`
	TotalFreedObject     json.Number `json:"total_freed_object"`
	
	// Ruby 2.1+
	HeapLiveSlot               json.Number `json:"heap_live_slot"`
	HeapFreeSlot               json.Number `json:"heap_free_slot"`
	HeapFinalSlot              json.Number `json:"heap_final_slot"`
	HeapSweptSlot              json.Number `json:"heap_swept_slot"`
	HeapEdenPageLength         json.Number `json:"heap_eden_page_length"`
	HeapTombPageLength         json.Number `json:"heap_tomb_page_length"`
	MallocIncrease             json.Number `json:"malloc_increase"`
	MallocLimit                json.Number `json:"malloc_limit"`
	MinorGcCount               json.Number `json:"minor_gc_count"`
	MajorGcCount               json.Number `json:"major_gc_count"`
	RememberedShadyObject      json.Number `json:"remembered_shady_object"`
	RememberedShadyObjectLimit json.Number `json:"remembered_shady_object_limit"`
	OldObject                  json.Number `json:"old_object"`
	OldObjectLimit             json.Number `json:"old_object_limit"`
	OldmallocIncrease          json.Number `json:"oldmalloc_increase"`
	OldmallocLimit             json.Number `json:"oldmalloc_limit"`
	
	// Ruby 2.2+
	HeapAllocatedPages                  json.Number `json:"heap_allocated_pages"`
	HeapSortedLength                    json.Number `json:"heap_sorted_length"`
	HeapAllocatablePages                json.Number `json:"heap_allocatable_pages"`
	HeapAvailableSlots                  json.Number `json:"heap_available_slots"`
	HeapLiveSlots                       json.Number `json:"heap_live_slots"`
	HeapFreeSlots                       json.Number `json:"heap_free_slots"`
	HeapFinalSlots                      json.Number `json:"heap_final_slots"`
	HeapMarkedSlots                     json.Number `json:"heap_marked_slots"`
	HeapSweptSlots                      json.Number `json:"heap_swept_slots"`
	HeapEdenPages                       json.Number `json:"heap_eden_pages"`
	HeapTombPages                       json.Number `json:"heap_tomb_pages"`
	TotalAllocatedPages                 json.Number `json:"total_allocated_pages"`
	TotalFreedPages                     json.Number `json:"total_freed_pages"`
	TotalAllocatedObjects               json.Number `json:"total_allocated_objects"`
	TotalFreedObjects                   json.Number `json:"total_freed_objects"`
	MallocIncreaseBytes                 json.Number `json:"malloc_increase_bytes"`
	MallocIncreaseBytesLimit            json.Number `json:"malloc_increase_bytes_limit"`
	RememberedWbUnprotectedObjects      json.Number `json:"remembered_wb_unprotected_objects"`
	RememberedWbUnprotectedObjectsLimit json.Number `json:"remembered_wb_unprotected_objects_limit"`
	OldObjects                          json.Number `json:"old_objects"`
	OldObjectsLimit                     json.Number `json:"old_objects_limit"`
	OldmallocIncreaseBytes              json.Number `json:"oldmalloc_increase_bytes"`
	OldmallocIncreaseBytesLimit         json.Number `json:"oldmalloc_increase_bytes_limit"`
}

// GCParser parses GC statistics
type GCParser struct{}

// ParseGCStats parses the /gc-stats response
func (p *GCParser) ParseGCStats(data []byte) (*domain.MetricCollection, error) {
	var stats GCStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse GC stats: %w", err)
	}

	collection := domain.NewMetricCollection()

	// Basic GC count (always available)
	if count, err := stats.Count.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.count",
			Value: count,
			Type:  domain.MetricTypeCounter,
		})
	}

	// Minor/Major GC counts (Ruby 2.1+)
	if minor, err := stats.MinorGcCount.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.minor_count",
			Value: minor,
			Type:  domain.MetricTypeCounter,
		})
	}
	if major, err := stats.MajorGcCount.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.major_count",
			Value: major,
			Type:  domain.MetricTypeCounter,
		})
	}

	// Heap slots - try different field names for compatibility
	// Ruby 2.2+ uses heap_available_slots
	if slots, err := stats.HeapAvailableSlots.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.heap_available_slots",
			Value: slots,
			Type:  domain.MetricTypeGauge,
		})
	}

	// Live slots - try all possible field names
	for _, field := range []json.Number{stats.HeapLiveSlots, stats.HeapLiveSlot, stats.HeapLiveNum} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.heap_live_slots",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Free slots
	for _, field := range []json.Number{stats.HeapFreeSlots, stats.HeapFreeSlot, stats.HeapFreeNum} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.heap_free_slots",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Final slots
	for _, field := range []json.Number{stats.HeapFinalSlots, stats.HeapFinalSlot, stats.HeapFinalNum} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.heap_final_slots",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Marked slots (Ruby 2.2+)
	if marked, err := stats.HeapMarkedSlots.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.heap_marked_slots",
			Value: marked,
			Type:  domain.MetricTypeGauge,
		})
	}

	// Old objects
	for _, field := range []json.Number{stats.OldObjects, stats.OldObject} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.old_objects",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Old objects limit
	for _, field := range []json.Number{stats.OldObjectsLimit, stats.OldObjectLimit} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.old_objects_limit",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Old malloc bytes
	for _, field := range []json.Number{stats.OldmallocIncreaseBytes, stats.OldmallocIncrease} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.oldmalloc_bytes",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// Old malloc limit
	for _, field := range []json.Number{stats.OldmallocIncreaseBytesLimit, stats.OldmallocLimit} {
		if val, err := field.Float64(); err == nil {
			_ = collection.Add(domain.Metric{
				Name:  "ruby.gc.oldmalloc_limit",
				Value: val,
				Type:  domain.MetricTypeGauge,
			})
			break
		}
	}

	// For backward compatibility with simple parsers
	if used, err := stats.HeapUsed.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.heap_used",
			Value: used,
			Type:  domain.MetricTypeGauge,
		})
	}
	if length, err := stats.HeapLength.Float64(); err == nil {
		_ = collection.Add(domain.Metric{
			Name:  "ruby.gc.heap_length",
			Value: length,
			Type:  domain.MetricTypeGauge,
		})
	}

	return collection, nil
}