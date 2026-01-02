package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kawai-network/veridium/pkg/fantasy/providers/pooled"
)

// ProviderMetricsService tracks and monitors pooled provider metrics
type ProviderMetricsService struct {
	mu        sync.RWMutex
	providers map[string]*pooled.PooledProvider // provider name -> provider
	stopCh    chan struct{}
}

// NewProviderMetricsService creates a new metrics service
func NewProviderMetricsService() *ProviderMetricsService {
	return &ProviderMetricsService{
		providers: make(map[string]*pooled.PooledProvider),
		stopCh:    make(chan struct{}),
	}
}

// RegisterProvider registers a pooled provider for monitoring
func (s *ProviderMetricsService) RegisterProvider(name string, provider *pooled.PooledProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.providers[name] = provider
	log.Printf("📊 [Metrics] Registered provider: %s", name)
}

// GetProviderMetrics returns metrics for a specific provider
func (s *ProviderMetricsService) GetProviderMetrics(name string) (*pooled.PoolMetricsSnapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	provider, exists := s.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider not found: %s", name)
	}
	
	return provider.GetMetrics(), nil
}

// GetAllMetrics returns metrics for all registered providers
func (s *ProviderMetricsService) GetAllMetrics() map[string]*pooled.PoolMetricsSnapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make(map[string]*pooled.PoolMetricsSnapshot)
	for name, provider := range s.providers {
		result[name] = provider.GetMetrics()
	}
	
	return result
}

// StartMonitoring starts background monitoring and logging
func (s *ProviderMetricsService) StartMonitoring(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	
	log.Printf("📊 [Metrics] Monitoring started (interval: %v)", interval)
	
	for {
		select {
		case <-ctx.Done():
			log.Printf("📊 [Metrics] Monitoring stopped (context cancelled)")
			return
		case <-s.stopCh:
			log.Printf("📊 [Metrics] Monitoring stopped (manual stop)")
			return
		case <-ticker.C:
			s.logMetrics()
		}
	}
}

// StopMonitoring stops the monitoring service
func (s *ProviderMetricsService) StopMonitoring() {
	close(s.stopCh)
}

// logMetrics logs current metrics for all providers
func (s *ProviderMetricsService) logMetrics() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if len(s.providers) == 0 {
		return
	}
	
	log.Printf("📊 ========== Provider Metrics Report ==========")
	
	for name, provider := range s.providers {
		metrics := provider.GetMetrics()
		if metrics == nil {
			continue
		}
		
		log.Printf("📊 Provider: %s", name)
		log.Printf("   Total Keys: %d | Healthy: %d", metrics.TotalKeys, metrics.HealthyKeys)
		log.Printf("   Total Requests: %d | Success Rate: %.2f%%", 
			metrics.TotalRequests, metrics.SuccessRate*100)
		log.Printf("   Rotations: %d | Current Key: #%d", 
			metrics.RotationCount, metrics.CurrentKeyIndex)
		
		// Log per-key metrics
		for i, km := range metrics.KeyMetrics {
			status := "✅"
			if !km.IsHealthy {
				status = "❌"
			}
			log.Printf("   Key #%d %s: Requests=%d | Success=%d | Failed=%d | RateLimited=%d", 
				i, status, km.TotalRequests, km.SuccessCount, km.FailureCount, km.RateLimited)
			
			if km.LastError != "" {
				log.Printf("      Last Error: %s (at %v)", km.LastError, km.LastErrorTime.Format("15:04:05"))
			}
		}
	}
	
	log.Printf("📊 =============================================")
}

// GetHealthStatus returns overall health status
func (s *ProviderMetricsService) GetHealthStatus() HealthStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	status := HealthStatus{
		Timestamp:      time.Now(),
		TotalProviders: len(s.providers),
		HealthyProviders: 0,
		Providers:      make(map[string]ProviderHealth),
	}
	
	for name, provider := range s.providers {
		metrics := provider.GetMetrics()
		if metrics == nil {
			continue
		}
		
		providerHealth := ProviderHealth{
			Name:         name,
			TotalKeys:    metrics.TotalKeys,
			HealthyKeys:  metrics.HealthyKeys,
			SuccessRate:  metrics.SuccessRate,
			IsHealthy:    metrics.HealthyKeys > 0 && metrics.SuccessRate > 0.5,
		}
		
		status.Providers[name] = providerHealth
		
		if providerHealth.IsHealthy {
			status.HealthyProviders++
		}
	}
	
	return status
}

// HealthStatus represents overall system health
type HealthStatus struct {
	Timestamp        time.Time                  `json:"timestamp"`
	TotalProviders   int                        `json:"total_providers"`
	HealthyProviders int                        `json:"healthy_providers"`
	Providers        map[string]ProviderHealth  `json:"providers"`
}

// ProviderHealth represents health of a single provider
type ProviderHealth struct {
	Name        string  `json:"name"`
	TotalKeys   int     `json:"total_keys"`
	HealthyKeys int     `json:"healthy_keys"`
	SuccessRate float64 `json:"success_rate"`
	IsHealthy   bool    `json:"is_healthy"`
}

// GetMetricsSummary returns a human-readable summary
func (s *ProviderMetricsService) GetMetricsSummary() string {
	health := s.GetHealthStatus()
	
	summary := fmt.Sprintf("Provider Metrics Summary (%s)\n", health.Timestamp.Format("2006-01-02 15:04:05"))
	summary += fmt.Sprintf("Total Providers: %d | Healthy: %d\n", health.TotalProviders, health.HealthyProviders)
	summary += "\nPer-Provider Status:\n"
	
	for name, ph := range health.Providers {
		status := "✅ Healthy"
		if !ph.IsHealthy {
			status = "❌ Unhealthy"
		}
		summary += fmt.Sprintf("  %s: %s (Keys: %d/%d, Success: %.1f%%)\n", 
			name, status, ph.HealthyKeys, ph.TotalKeys, ph.SuccessRate*100)
	}
	
	return summary
}

