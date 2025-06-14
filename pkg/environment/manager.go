package environment

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	// ErrInvalidEnv represents that the given runtime environment is not
	// registered or supported.
	ErrInvalidEnv = errors.New("environment: invalid env")

	// ErrLocked indicates that the current runtime environment is locked
	// and cannot be changed.
	ErrLocked = errors.New("environment: locked")
)

// Manager interface type defines a runtime environment manager.
// An independent system should share an independent manager instance.
type Manager interface {
	// Get method returns the current runtime environment.
	Get() Env

	// Is returns whether the given runtime environment is equal to the
	// current runtime environment.
	Is(Env) bool

	// In returns whether the current runtime environment is in the given
	// runtime environment list.
	In(envs []Env) bool

	// Register method registers a custom runtime environment.
	// If you want to add a custom environment, this method must be called
	// before the Manager.Set() method.
	Register(Env)

	// Registered determines whether the given runtime environment is already registered.
	Registered(Env) bool

	// Lock method locks the current runtime environment.
	// After locking, the current runtime environment cannot be changed.
	Lock()

	// Locked method returns whether the current runtime environment is locked.
	Locked() bool

	// Set method sets the current runtime environment.
	// If the given runtime environment is not supported, ErrInvalidEnv error is returned.
	// If the current runtime environment is locked, ErrLocked error is returned.
	Set(Env) error

	// SetAndLock method sets and locks the current runtime environment.
	// If the runtime environment settings fail, they are not locked.
	SetAndLock(Env) error

	// Listen method adds a given runtime environment listener.
	// If the given listener is nil, ignore it.
	Listen(Listener)

	// UnListen removes and returns to the recently added listener.
	// If there is no listener to be removed, nil is returned.
	UnListen() Listener

	// UnListenAll removes and returns all added listeners.
	// If there is no listener to be removed, nil is returned.
	UnListenAll() []Listener
}

// New creates and returns a new instance of the built-in runtime environment manager.
// The default runtime environment is Development, and all built-in runtime environments
// have been registered.
func New() Manager {
	return &manager{
		current:    Development,
		registered: []Env{Development, Testing, Prerelease, Production},
	}
}

// NewEmpty creates and returns an empty instance of the runtime environment manager.
// The manager returned by this function does not register any runtime environment,
// and the current runtime environment is empty.
func NewEmpty() Manager {
	return new(manager)
}

// Listener defines the runtime environment listener.
// Listeners are used to receive notifications when the runtime environment changes.
type Listener func(after, before Env)

// This is a built-in runtime environment manager.
type manager struct {
	mutex      sync.RWMutex
	current    Env
	locked     int32
	registered []Env
	listeners  []Listener
}

// Get method returns the current runtime environment.
func (m *manager) Get() Env {
	return m.current
}

// Is returns whether the given runtime environment is equal to the
// current runtime environment.
func (m *manager) Is(env Env) bool {
	return m.current.Is(env)
}

// In returns whether the current runtime environment is in the given
// runtime environment list.
func (m *manager) In(envs []Env) bool {
	return m.current.In(envs)
}

// Register method registers a custom runtime environment.
// If you want to add a custom environment, this method must be called
// before the Manager.Set() method.
func (m *manager) Register(env Env) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !env.In(m.registered) {
		m.registered = append(m.registered, env)
	}
}

// Registered determines whether the given runtime environment is already registered.
func (m *manager) Registered(env Env) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return env.In(m.registered)
}

// Lock method locks the current runtime environment.
// After locking, the current runtime environment cannot be changed.
func (m *manager) Lock() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	atomic.StoreInt32(&m.locked, 1)
}

// Locked method returns whether the current runtime environment is locked.
func (m *manager) Locked() bool {
	return atomic.LoadInt32(&m.locked) == 1
}

// Set method sets the current runtime environment.
// If the given runtime environment is not supported, ErrInvalidEnv error is returned.
// If the current runtime environment is locked, ErrLocked error is returned.
func (m *manager) Set(env Env) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.set(env)
}

// SetAndLock method sets and locks the current runtime environment.
// If the runtime environment settings fail, they are not locked.
func (m *manager) SetAndLock(env Env) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.set(env); err != nil {
		return err
	}

	atomic.StoreInt32(&m.locked, 1)
	return nil
}

// Sets the current runtime environment.
func (m *manager) set(env Env) error {
	if m.Locked() {
		return ErrLocked
	}
	if !env.In(m.registered) {
		return ErrInvalidEnv
	}

	if old := m.current; !old.Is(env) {
		m.current = env
		// Trigger all listeners synchronously.
		for i, j := 0, len(m.listeners); i < j; i++ {
			m.listeners[i](env, old)
		}
	}
	return nil
}

// Listen method adds a given runtime environment listener.
// If the given listener is nil, ignore it.
func (m *manager) Listen(listener Listener) {
	if listener != nil {
		m.mutex.Lock()
		m.listeners = append(m.listeners, listener)
		m.mutex.Unlock()
	}
}

// UnListen removes and returns to the recently added listener.
// If there is no listener to be removed, nil is returned.
func (m *manager) UnListen() (r Listener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if n := len(m.listeners) - 1; n >= 0 {
		r = m.listeners[n]
		if n == 0 {
			m.listeners = nil
		} else {
			m.listeners = m.listeners[:n]
		}
	}
	return
}

// UnListenAll removes and returns all added listeners.
// If there is no listener to be removed, nil is returned.
func (m *manager) UnListenAll() (r []Listener) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.listeners) > 0 {
		r = m.listeners
		m.listeners = nil
	}
	return
}
