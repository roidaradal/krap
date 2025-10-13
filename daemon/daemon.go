package daemon

import (
	"fmt"
	"time"

	"github.com/roidaradal/fn/clock"
	"github.com/roidaradal/fn/dict"
	"github.com/roidaradal/fn/io"
)

/*
Note: DaemonConfig is expected to have this structure:
Use int for Interval and Margin values so that we can disable
a daemon by setting the value to -1 or any value < 0

type Config sturct {
	<Domain> struct {
		<FeatureInterval> int
		...
	}
	...
}
*/

// Load Daemon Config which follows the expected structure,
// Validates if any of the Interval values are 0 (invalid)
func LoadConfig[T any](path string) (*T, error) {
	cfg, err := io.ReadJSON[T](path)
	if err != nil {
		return nil, err
	}
	cfgMap, err := dict.FromStruct[T, map[string]int](cfg)
	if err != nil {
		return nil, err
	}
	for key := range cfgMap {
		for cfgKey, value := range cfgMap[key] {
			if value == 0 {
				return nil, fmt.Errorf("invalid daemon %s.%s: %d", key, cfgKey, value)
			}
		}
	}
	return cfg, nil
}

// Runs a task every given interval,
// TimeScale = time.Hour, time.Minute, time.Second
func Run(name string, task func(), interval int, timeScale time.Duration) {
	if interval < 0 {
		fmt.Printf("Daemon:%s is disabled\n", name)
		return
	}
	timeInterval := time.Duration(interval) * timeScale
	go func() {
		for {
			start := clock.TimeNow()
			task()
			clock.Sleep(timeInterval, start)
		}
	}()
}
