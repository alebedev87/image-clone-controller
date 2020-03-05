package utils

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestBlacklistNamespacePredicate(t *testing.T) {
	testCases := []struct {
		name      string
		blacklist map[string]bool
		ns        string
		expected  bool
	}{
		{
			name:      "Go",
			blacklist: map[string]bool{"bad": true},
			ns:        "good",
			expected:  true,
		},
		{
			name:      "No go",
			blacklist: map[string]bool{"bad": true},
			ns:        "bad",
			expected:  false,
		},
		{
			name:      "No go 2",
			blacklist: map[string]bool{"bad": true, "too bad": true},
			ns:        "bad",
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pred := NewBlacklistNamespacePredicate(tc.blacklist)
			output := pred.Create(newTestCreateEvent(tc.ns))
			if output != tc.expected {
				t.Errorf("Create event not handled correctly. Expected %t, got %t", tc.expected, output)
			}
			output = pred.Update(newTestUpdateEvent(tc.ns))
			if output != tc.expected {
				t.Errorf("Update event not handled correctly. Expected %t, got %t", tc.expected, output)
			}
			output = pred.Delete(newTestDeleteEvent(tc.ns))
			if output != tc.expected {
				t.Errorf("Delete event not handled correctly. Expected %t, got %t", tc.expected, output)
			}
			output = pred.Generic(newTestGenericEvent(tc.ns))
			if output != tc.expected {
				t.Errorf("Generic event not handled correctly. Expected %t, got %t", tc.expected, output)
			}
		})
	}
}

func newTestCreateEvent(ns string) event.CreateEvent {
	return event.CreateEvent{
		Meta: &metav1.ObjectMeta{
			Namespace: ns,
		},
	}
}

func newTestUpdateEvent(ns string) event.UpdateEvent {
	return event.UpdateEvent{
		MetaOld: &metav1.ObjectMeta{
			Namespace: "ERROR! this namespace must not be used",
		},
		MetaNew: &metav1.ObjectMeta{
			Namespace: ns,
		},
	}
}

func newTestDeleteEvent(ns string) event.DeleteEvent {
	return event.DeleteEvent{
		Meta: &metav1.ObjectMeta{
			Namespace: ns,
		},
	}
}

func newTestGenericEvent(ns string) event.GenericEvent {
	return event.GenericEvent{
		Meta: &metav1.ObjectMeta{
			Namespace: ns,
		},
	}
}
