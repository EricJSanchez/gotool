package sys

import "gotool/pkg/environment"

func Env() environment.Env {
	return environment.Get()
}

func InitEnv(env environment.Env) error {
	return environment.SetAndLock(env)
}
