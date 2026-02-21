package command

import (
	"fmt"
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
			// Don't count the object, just its keys.
			keys += len(sub)
		}

		if _, ok := v.(string); ok {
			// v is just a string key.
			keys++
		}

		if _, ok := v.([]string); ok {
			// v is just a string key.
			keys++
		}

		if _, ok := v.(bool); ok {
			// v is just a bool key.
			keys++
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
	t.Run("current env vars set", func(t *testing.T) {
		orig := getNotiEnv(t)
		defer setNotiEnv(t, orig)

		clearNotiEnv(t)

		v := viper.New()
		bindNotiEnv(v)

		haveKeys := countSettingsKeys(t, v.AllSettings())
		if haveKeys != 0 {
			t.Error("Environment should be cleared")
			t.Error(v.AllSettings())
		}

		var numSet int
		for _, env := range keyEnvBindings {
			if err := os.Setenv(env, "foo"); err != nil {
				t.Errorf("Setenv error: %s", err)
				continue
			}
			numSet++
		}

		haveKeys = countSettingsKeys(t, v.AllSettings())
		wantKeys := numSet
		if haveKeys != wantKeys {
			t.Error("Unexpected base config length")
			t.Errorf("have=%d; want=%d", haveKeys, wantKeys)
			t.Error(v.AllSettings())
		}
	})

	t.Run("deprecated env vars set", func(t *testing.T) {
		orig := getNotiEnv(t)
		defer setNotiEnv(t, orig)

		clearNotiEnv(t)

		var numSet int
		for oldEnv := range keyEnvBindingsDeprecated {
			if err := os.Setenv(oldEnv, "foo"); err != nil {
				t.Errorf("Setenv error: %s", err)
				continue
			}
			numSet++
		}

		v := viper.New()
		bindNotiEnv(v)

		haveKeys := countSettingsKeys(t, v.AllSettings())
		wantKeys := numSet
		if haveKeys != wantKeys {
			t.Error("Unexpected base config length")
			t.Errorf("have=%d; want=%d", haveKeys, wantKeys)
			t.Error(v.AllSettings())
		}
	})
}

func TestSetupConfigFile(t *testing.T) {
	v := viper.New()
	if err := setupConfigFile("testdata/noti.yaml", v); err != nil {
		t.Error(err)
	}

	const want = 1
	have := countSettingsKeys(t, v.AllSettings())
	if have != want {
		t.Error("Unexpected number of keys")
		t.Errorf("have=%d; want=%d", have, want)
	}
}

func TestConfigureApp(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)

	cases := []struct {
		name       string
		configFile string
		env        string
		want       string
	}{
		{
			// Config file should take precedence.
			name:       "defaults and file",
			configFile: "testdata/noti.yaml",
			want:       "testSoundName",
		},
		{
			// Env should take precedence.
			name:       "defaults, file, and env",
			configFile: "testdata/noti.yaml",
			env:        "NOTI_NSUSER_SOUNDNAME",
			want:       "testSoundName",
		},
		{
			// Defaults should take precedence.
			name: "defaults",
			want: baseDefaults["nsuser.soundName"].(string),
		},
	}

	for _, c := range cases {
		// Pin case scope.
		c := c

		t.Run(c.name, func(t *testing.T) {
			clearNotiEnv(t)

			v := viper.New()
			flags := pflag.NewFlagSet("testconfigureapp", pflag.ContinueOnError)
			InitFlags(flags)

			if c.configFile != "" {
				flags.Set("file", c.configFile)
			}
			if c.env != "" {
				if err := os.Setenv(c.env, c.want); err != nil {
					t.Errorf("Failed to set env: %s", err)
				}
			}

			if err := configureApp(v, flags); err != nil {
				t.Error(err)
			}

			have := v.GetString("nsuser.soundName")
			if have != c.want {
				t.Error("Unexpected config value")
				t.Errorf("have=%s; want=%s", have, c.want)
				t.Error("nsuser:", v.Sub("nsuser").AllSettings())
			}
		})
	}
}

func TestEnabledServices(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)
	clearNotiEnv(t)

	t.Run("flag override", func(t *testing.T) {
		v := viper.New()
		// For tests, we prepend the testdata dir so that we check for a config
		// file there first.
		v.AddConfigPath("testdata")

		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		InitFlags(flags)

		configureApp(v, flags)

		want := true
		flags.Set("slack", fmt.Sprint(want))
		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["slack"]
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})

	t.Run("non-service flags", func(t *testing.T) {
		v := viper.New()
		// For tests, we prepend the testdata dir so that we check for a config
		// file there first.
		v.AddConfigPath("testdata")

		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		InitFlags(flags)

		configureApp(v, flags)

		flags.Set("verbose", "true")
		services := enabledServices(v, flags)

		// We should end up taking the defaults.

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
			t.Error("services:", services)
		}

		want := true
		_, have := services["banner"]
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})

	t.Run("env override", func(t *testing.T) {
		v := viper.New()
		// For tests, we prepend the testdata dir so that we check for a config
		// file there first.
		v.AddConfigPath("testdata")

		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		InitFlags(flags)

		configureApp(v, flags)

		if err := os.Setenv("NOTI_DEFAULT", "slack"); err != nil {
			t.Fatal(err)
		}
		defer os.Unsetenv("NOTI_DEFAULT")

		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["slack"]
		want := true
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		v := viper.New()
		// For tests, we prepend the testdata dir so that we check for a config
		// file there first.
		v.AddConfigPath("testdata")

		flags := pflag.NewFlagSet("testenabledservices", pflag.ContinueOnError)
		InitFlags(flags)

		configureApp(v, flags)

		services := enabledServices(v, flags)

		if len(services) != 1 {
			t.Error("Unexpected number of enabled services")
			t.Errorf("have=%d; want=%d", len(services), 1)
		}

		_, have := services["banner"]
		want := true
		if have != want {
			t.Error("Unexpected enabled state")
			t.Errorf("have=%t; want=%t", have, want)
		}
	})
}

func TestBannerIconConfig(t *testing.T) {
	orig := getNotiEnv(t)
	defer setNotiEnv(t, orig)

	t.Run("default is empty", func(t *testing.T) {
		clearNotiEnv(t)

		v := viper.New()
		flags := pflag.NewFlagSet("testbannericon", pflag.ContinueOnError)
		InitFlags(flags)

		if err := configureApp(v, flags); err != nil {
			t.Error(err)
		}

		have := v.GetString("banner.icon")
		if have != "" {
			t.Errorf("have=%q; want=%q", have, "")
		}
	})

	t.Run("env var", func(t *testing.T) {
		clearNotiEnv(t)

		const want = "/tmp/icon.png"
		os.Setenv("NOTI_BANNER_ICON", want)
		defer os.Unsetenv("NOTI_BANNER_ICON")

		v := viper.New()
		flags := pflag.NewFlagSet("testbannericon", pflag.ContinueOnError)
		InitFlags(flags)

		if err := configureApp(v, flags); err != nil {
			t.Error(err)
		}

		have := v.GetString("banner.icon")
		if have != want {
			t.Errorf("have=%q; want=%q", have, want)
		}
	})

	t.Run("flag", func(t *testing.T) {
		clearNotiEnv(t)

		const want = "/tmp/flag-icon.png"
		v := viper.New()
		flags := pflag.NewFlagSet("testbannericon", pflag.ContinueOnError)
		InitFlags(flags)
		flags.Set("icon", want)

		if err := configureApp(v, flags); err != nil {
			t.Error(err)
		}

		have := v.GetString("banner.icon")
		if have != want {
			t.Errorf("have=%q; want=%q", have, want)
		}
	})

	t.Run("flag overrides env", func(t *testing.T) {
		clearNotiEnv(t)

		const want = "/tmp/flag-wins.png"
		os.Setenv("NOTI_BANNER_ICON", "/tmp/env-loses.png")
		defer os.Unsetenv("NOTI_BANNER_ICON")

		v := viper.New()
		flags := pflag.NewFlagSet("testbannericon", pflag.ContinueOnError)
		InitFlags(flags)
		flags.Set("icon", want)

		if err := configureApp(v, flags); err != nil {
			t.Error(err)
		}

		have := v.GetString("banner.icon")
		if have != want {
			t.Errorf("have=%q; want=%q", have, want)
		}
	})
}

func TestGetNotifications(t *testing.T) {
	services := []string{
		"banner",
		"bearychat",
		"keybase",
		"pushbullet",
		"pushover",
		"pushsafer",
		"simplepush",
		"slack",
		"speech",
		"zulip",
	}

	for _, name := range services {
		// Pin name scope.
		name := name
		t.Run(fmt.Sprintf("get %s notification", name), func(t *testing.T) {
			v := viper.New()
			s := map[string]struct{}{name: {}}

			notis := getNotifications(v, s)
			if len(notis) != 1 {
				t.Error("Unexpected number of notifications")
				t.Errorf("have=%d; want=%d", len(notis), 1)
			}
		})
	}
}
