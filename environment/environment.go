package environment

const (
	// Development is the development environment, which is also the default
	// runtime environment.
	Development Env = "development"

	// Testing is a test environment, usually used for initial quality acceptance.
	Testing Env = "testing"

	// Prerelease is a pre release environment, usually used for grayscale
	// testing or quality acceptance.
	Prerelease Env = "prerelease"

	// Production is the production environment and the final deployment
	// environment of the application.
	Production Env = "production"
)

// The global default runtime environment manager.
var defaultManager = New()

// Env type defines the runtime environment.
type Env string

// String method returns the current runtime environment string.
func (e Env) String() string { return string(e) }

// Is method returns whether the given runtime environment is equal to the
// current runtime environment.
func (e Env) Is(env Env) bool { return e == env }

// In method returns whether the current runtime environment is in the given
// runtime environment list.
func (e Env) In(envs []Env) bool {
	for i, j := 0, len(envs); i < j; i++ {
		if e.Is(envs[i]) {
			return true
		}
	}
	return false
}

// Get returns the current runtime environment.
func Get() Env { return defaultManager.Get() }

// Is returns whether the given runtime environment is equal to the
// current runtime environment.
func Is(env Env) bool { return defaultManager.Is(env) }

// In returns whether the current runtime environment is in the given
// runtime environment list.
func In(envs []Env) bool { return defaultManager.In(envs) }

// Register registers a custom runtime environment.
// If you want to add a custom environment, this method must be called
// before the Set() method.
func Register(env Env) { defaultManager.Register(env) }

// Registered determines whether the given runtime environment is already registered.
func Registered(env Env) bool { return defaultManager.Registered(env) }

// Lock locks the current runtime environment.
// After locking, the current runtime environment cannot be changed.
func Lock() { defaultManager.Lock() }

// Locked returns whether the current runtime environment is locked.
func Locked() bool { return defaultManager.Locked() }

// Set sets the current runtime environment.
// If the given runtime environment is not supported, ErrInvalidEnv error is returned.
// If the current runtime environment is locked, ErrLocked error is returned.
func Set(env Env) error { return defaultManager.Set(env) }

// SetAndLock sets and locks the current runtime environment.
// If the runtime environment settings fail, they are not locked.
func SetAndLock(env Env) error { return defaultManager.SetAndLock(env) }

// Listen adds a given runtime environment listener.
// If the given listener is nil, ignore it.
func Listen(listener Listener) { defaultManager.Listen(listener) }

// UnListen removes and returns to the recently added listener.
// If there is no listener to be removed, nil is returned.
func UnListen() Listener { return defaultManager.UnListen() }

// UnListenAll removes and returns all added listeners.
// If there is no listener to be removed, nil is returned.
func UnListenAll() []Listener { return defaultManager.UnListenAll() }
