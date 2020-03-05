package utils

import (
	"image-clone-controller/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var _ predicate.Predicate = &BlacklistNamespacePredicate{}

// BlacklistNamespacePredicate is a predicate to not process events from the blacklisted namespaces
type BlacklistNamespacePredicate struct {
	list map[string]bool
}

// NewBlacklistNamespacePredicate returns an instance of BlacklistNamespacePredicate
func NewBlacklistNamespacePredicate(blacklist map[string]bool) *BlacklistNamespacePredicate {
	return &BlacklistNamespacePredicate{
		list: blacklist,
	}
}

// NewBlacklistNamespacePredicateFromConfig returns an instance of BlacklistNamespacePredicate set from the program configuration
func NewBlacklistNamespacePredicateFromConfig() *BlacklistNamespacePredicate {
	return &BlacklistNamespacePredicate{
		list: config.GlobalConfig.NamespaceBlacklist(),
	}
}

// Create returns true if the create event should be processed
func (p *BlacklistNamespacePredicate) Create(e event.CreateEvent) bool {
	_, exists := p.list[e.Meta.GetNamespace()]
	return !exists
}

// Update returns true if the update event should be processed
func (p *BlacklistNamespacePredicate) Update(e event.UpdateEvent) bool {
	_, exists := p.list[e.MetaNew.GetNamespace()]
	return !exists
}

// Delete returns true if the delete event should be processed
func (p *BlacklistNamespacePredicate) Delete(e event.DeleteEvent) bool {
	_, exists := p.list[e.Meta.GetNamespace()]
	return !exists
}

// Generic  returns true if the generic event should be processed
func (p *BlacklistNamespacePredicate) Generic(e event.GenericEvent) bool {
	_, exists := p.list[e.Meta.GetNamespace()]
	return !exists
}
