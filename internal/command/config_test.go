package command

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func countSettingsKeys(t *testing.T, m map[string]interface{}) int {
	t.Helper()

	var keys int
	for _, v := range m {
		if sub, ok := v.(map[string]interface{}); ok {
			keys += len(sub)
		}
	}
	return keys
}

func TestSetNotiDefaults(t *testing.T) {
	v := viper.New()
	setNotiDefaults(v)

	haveKeys := countSettingsKeys(t, v.AllSettings())
	if haveKeys != len(baseDefaults) {
		t.Error("Unexpected base config length")
		t.Errorf("have=%d; want=%d", haveKeys, len(baseDefaults))
	}
}

func getNotiEnv(t *testing.T) map[string]string {
	t.Helper()

	notiEnv := make(map[string]string)
	for _, env := range keyEnvBindings {
		notiEnv[env] = os.Getenv(env)
	}
	return notiEnv
}

func clearNotiEnv(t *testing.T) {
	t.Helper()

	for _, env := range keyEnvBindings {
		if err := os.Unsetenv(env); err != nil {
			t.Fatalf("failed to clear noti env: %s", err)
		}
	}
}

func setNotiEnv(t *testing.T, m map[string]string) {
	t.Helper()

	for env, val := range m {
		if err := os.Setenv(env, val); err != nil {
			t.Fatalf("failed to set noti env: %s", err)
		}
	}
}

func TestBindNotiEnv(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)

	clearNotiEnv(t)

	v := viper.New()
	bindNotiEnv(v)

	haveKeys := countSettingsKeys(t, v.AllSettings())
	if haveKeys != 0 {
		t.Fatal("Environment should be cleared")
	}

	for _, env := range keyEnvBindings {
		if err := os.Setenv(env, "foo"); err != nil {
			t.Errorf("Setenv error: %s", err)
		}
	}

	haveKeys = countSettingsKeys(t, v.AllSettings())
	if haveKeys != len(baseDefaults) {
		t.Error("Unexpected base config length")
		t.Errorf("have=%d; want=%d", haveKeys, len(baseDefaults))
	}
}

func TestSetupConfigFile(t *testing.T) {
	v := viper.New()
	setupConfigFile(v)

	haveKeys := countSettingsKeys(t, v.AllSettings())
	if haveKeys != 0 {
		t.Fatal("Environment should be cleared")
	}

	sample, err := ioutil.ReadFile("testdata/sample_config.yaml")
	if err != nil {
		t.Errorf("Failed to read sample config: %s", err)
	}

	if err := ioutil.WriteFile(".noti.yaml", sample, 0644); err != nil {
		t.Errorf("Failed to write sample config: %s", err)
	}
	defer os.Remove(".noti.yaml")

	if err := v.ReadInConfig(); err != nil {
		t.Errorf("Failed to read config: %s", err)
	}

	haveKeys = countSettingsKeys(t, v.AllSettings())
	if haveKeys != len(baseDefaults) {
		t.Error("Unexpected len keys")
	}
}

func TestConfigureApp(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)
	clearNotiEnv(t)

	v := viper.New()
	flags := pflag.NewFlagSet("testconfigureapp", pflag.ContinueOnError)

	configureApp(v, flags)

	t.Run("default config", func(t *testing.T) {
		have := v.GetString("nsuser.soundName")
		want := baseDefaults["nsuser.soundName"]
		if have != want {
			t.Error("Unexpected config value")
			t.Errorf("have=%s; want=%s", have, want)
		}
	})

	t.Run("env override", func(t *testing.T) {
		want := "foo"
		if err := os.Setenv("NOTI_SOUND", want); err != nil {
			t.Errorf("Failed to set env: %s", err)
		}

		have := v.GetString("nsuser.soundName")
		if have != want {
			t.Error("Unexpected config value")
			t.Errorf("have=%s; want=%s", have, want)
		}
	})
}
